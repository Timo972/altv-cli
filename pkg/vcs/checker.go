package vcs

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/timo972/altv-cli/pkg/cdn"
	"github.com/timo972/altv-cli/pkg/logging"
	"github.com/timo972/altv-cli/pkg/platform"
	"github.com/timo972/altv-cli/pkg/version"
)

type ModuleStatus uint8

const (
	StatusValid ModuleStatus = 1 << iota
	StatusInvalid
	StatusUpToDate
	StatusUpgradable
)

func (status ModuleStatus) Add(status2 ModuleStatus) ModuleStatus {
	status |= status2
	return status
}

func (status ModuleStatus) Has(s ModuleStatus) bool {
	return status&s != 0
}

type ModuleStatusResult map[string]ModuleStatus

type extManifest struct {
	*cdn.Manifest
	mod string
}

type Checker interface {
	Verify(ctx context.Context, path string, remote bool) (ModuleStatusResult, error)
	AddCDN(cdn.CDN)
}

type checker struct {
	CDNRegistry
	arch    platform.Arch
	branch  version.Branch
	modules []string
}

func NewChecker(arch platform.Arch, branch version.Branch, modules []string, registry CDNRegistry) Checker {
	return &checker{
		CDNRegistry: registry,
		arch:        arch,
		branch:      branch,
		modules:     modules,
	}
}

// TODO: utilize goroutines to aggregate manifests simultaneously
func (c *checker) aggregateRemoteManifests() ([]*extManifest, error) {
	allMans := make([]*extManifest, 0)
	errs := make([]error, 0)
	for _, mod := range c.modules {
		cdn, ok := c.moduleCDN(mod)
		if !ok {
			err := newErrNoCDN(mod)
			logging.WarnLogger.Printf(err.Error())
			errs = append(errs, err)
			continue
		}
		logging.DebugLogger.Printf("cdn %v for module %s", cdn, mod)

		man, err := cdn.Manifest(c.branch, c.arch, mod)
		if err != nil {
			err = newErrNoManifest(mod, err)
			logging.WarnLogger.Printf(err.Error())
			errs = append(errs, err)
			continue
		}

		logging.DebugLogger.Printf("got manifest for module %s", mod)
		allMans = append(allMans, &extManifest{
			Manifest: man,
			mod:      mod,
		})
	}
	return allMans, errors.Join(errs...)
}

func (c *checker) aggregateLocalManifests(path string) ([]*extManifest, []string, error) {
	mans := make([]*extManifest, 0)
	mods := make([]string, 0)

	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		logging.DebugLogger.Printf("walking path %s", path)

		if d.IsDir() {
			return nil
		}

		if !strings.HasSuffix(d.Name(), ".update.json") {
			return nil
		}

		f, err := os.OpenFile(path, os.O_RDONLY, 0644)
		if err != nil {
			return err
		}
		defer f.Close()

		var man cdn.Manifest
		if err = json.NewDecoder(f).Decode(&man); err != nil {
			return err
		}

		mod := strings.TrimSuffix(d.Name(), ".update.json")
		mans = append(mans, &extManifest{
			Manifest: &man,
			mod:      mod,
		})
		mods = append(mods, mod)
		return nil
	})

	return mans, mods, err
}

func (c *checker) verifyWithManifest(path string, man *extManifest) (ModuleStatus, error) {
	// logging.DebugLogger.Printf("verify using manifest: %+v", man)
	status := StatusValid
	for fname, fhash := range man.HashList {
		if err := VerifyFileChecksum(path, fname, fhash, man.SizeList[fname]); err != nil {
			status = StatusInvalid
		}
	}
	return status, nil
}

type moduleStatusResp struct {
	err  error
	stat ModuleStatus
	mod  string
}

func (c *checker) verifyWithManifests(ctx context.Context, path string, mans []*extManifest) (ModuleStatusResult, error) {
	msrch := make(chan *moduleStatusResp, len(mans))
	for i, man := range mans {
		go func(man *extManifest, i int) {
			logging.DebugLogger.Printf("start module verify: %s", c.modules[i])
			stat, err := c.verifyWithManifest(path, man)
			logging.DebugLogger.Printf("got module status: %s %+v", c.modules[i], stat)
			msrch <- &moduleStatusResp{
				stat: stat,
				err:  err,
				mod:  man.mod,
			}
		}(man, i)
	}

	var err error
	i := 0
	msrs := ModuleStatusResult{}
	for {
		if i >= len(mans) {
			logging.DebugLogger.Printf("%d manifests done!", len(mans))
			break
		}

		select {
		case msr := <-msrch:
			logging.DebugLogger.Printf("received module status: %d %+v", i, msr)
			if msr.err != nil && err != nil {
				err = errors.Join(err, msr.err)
			} else if msr.err != nil {
				err = msr.err
			}

			msrs[msr.mod] = msr.stat
			i++
		case <-ctx.Done():
			return msrs, fmt.Errorf("verify canceled by context: %w", ctx.Err())
		}
	}

	return msrs, err
}

// func (c *checker) processLocalManifests(ctx context.Context, path string, mans []*extManifest) (ModuleStatusResult, error) {
// 	status, err := c.verifyWithManifests(ctx, path, mans)
// 	if err != nil {
// 		logging.WarnLogger.Printf("encountered errors while checking using local manifests: %v", err)
// 	}
// 	logging.DebugLogger.Printf("all module status: %+v", status)
// 	return status, err
// }

// func (c *checker) processRemoteManifests(ctx context.Context, path string, rmans []*extManifest) (ModuleStatusResult, error) {
// 	status, err := c.verifyWithManifests(ctx, path, rmans)
// 	return nil, nil
// }

func (c *checker) Verify(ctx context.Context, path string, remote bool) (ModuleStatusResult, error) {
	lmans, mods, err := c.aggregateLocalManifests(path)
	if err != nil && len(lmans) < 1 {
		return nil, err
	} else if err != nil {
		logging.WarnLogger.Printf("encountered errors while looking for local module manifests: %v", err)
	}

	lmansFound := len(lmans) > 0
	if lmansFound {
		c.modules = mods
	}

	logging.DebugLogger.Printf("got %d local manifests", len(lmans))

	// logic:
	// if no local manifests are found and remote = false: throw error
	// if local manifests are found and remote = false: only check with local
	// if local manifests are found and remote = true: first check with local, then for updates using remote
	// if no local manifests are found and remote = true: check with remote

	var rmans []*extManifest
	if remote {
		rmans, err = c.aggregateRemoteManifests()
		if err != nil && len(rmans) < 1 {
			return nil, err
		} else if err != nil {
			logging.WarnLogger.Printf("encountered errors while looking for remote module manifests: %v", err)
		}
	}

	switch true {
	case !lmansFound && !remote:
		return nil, fmt.Errorf("unable to verify files: no local manifests found and not allowed to fetch remote manifests")
	case lmansFound && !remote:
		return c.verifyWithManifests(ctx, path, lmans)
	case lmansFound && remote:
		logging.DebugLogger.Printf("checking local and remote manifests")
		lstatus, lerr := c.verifyWithManifests(ctx, path, lmans)
		if lerr != nil {
			err = errors.Join(err, lerr)
		}
		logging.DebugLogger.Printf("local manifests done!")

		rstatus, rerr := c.verifyWithManifests(ctx, path, rmans)
		if rerr != nil {
			err = errors.Join(err, rerr)
		}
		logging.DebugLogger.Printf("remote manifests done!")

		result := ModuleStatusResult{}
		for mod, lstat := range lstatus {
			if rstat, ok := rstatus[mod]; ok {
				switch rstat {
				case StatusInvalid:
					result[mod] = lstat.Add(StatusUpgradable)
				case StatusValid:
					result[mod] = lstat.Add(StatusUpToDate)
				default:
					result[mod] = lstat.Add(rstat)
				}
			} else {
				logging.DebugLogger.Printf("no remote status for %s", mod)
				result[mod] = lstat
			}
		}
		logging.DebugLogger.Printf("merged status!")
		return result, err
	case !lmansFound && remote:
		return c.verifyWithManifests(ctx, path, rmans)
	default:
		return nil, fmt.Errorf("unexpected switch case")
	}
}

func VerifyFileChecksum(path, fname, fhash string, fsize int) error {
	fpath := fmt.Sprintf("%s/%s", path, fname)
	logging.DebugLogger.Printf("verifying file: %s", fpath)
	file, err := os.OpenFile(fpath, os.O_RDONLY, 0644)
	if err != nil {
		logging.DebugLogger.Printf("error while verifying file: could not open %s", fpath)
		return err
	}
	defer file.Close()

	h := sha1.New()
	if _, err := io.Copy(h, file); err != nil {
		logging.DebugLogger.Printf("error while verifying file: could hash contents of %s", fpath)
		return err
	}

	checksum := hex.EncodeToString((h.Sum(nil)))
	if checksum != fhash {
		logging.DebugLogger.Printf("checksum missmatch for %s: expected %s, got %s", fpath, fhash, checksum)
		return fmt.Errorf("checksum missmatch for %s: expected %s, got %s", fpath, fhash, checksum)
	}

	if fsize < 0 {
		return nil
	}

	stat, err := file.Stat()
	if err != nil {
		logging.DebugLogger.Printf("file size read errpr for %s", fpath)
		return err
	}
	if stat.Size() != int64(fsize) {
		logging.DebugLogger.Printf("file size missmatch for %s: expected %d, got %d", fpath, fsize, stat.Size())
		return fmt.Errorf("file size missmatch for %s: expected %d, got %d", fpath, fsize, stat.Size())
	}

	return nil
}
