package vcs

import (
	"context"
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/timo972/altv-cli/pkg/cdn"
	"github.com/timo972/altv-cli/pkg/logging"
	"github.com/timo972/altv-cli/pkg/platform"
	"github.com/timo972/altv-cli/pkg/version"
)

type ModuleStatus struct{}
type ModuleManifest struct {
	*cdn.Manifest
	Mod string
}

type Checker interface {
	Verify(ctx context.Context, path string) ([]*ModuleStatus, error)
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
func (c *checker) AggregateRemoteManifests() ([]*ModuleManifest, error) {
	allMans := make([]*ModuleManifest, 0)
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
		allMans = append(allMans, &ModuleManifest{
			Manifest: man,
			Mod:      mod,
		})
	}
	return allMans, errors.Join(errs...)
}

func (c *checker) AggregateLocalManifests(path string) ([]*ModuleManifest, error) {
	mans := make([]*ModuleManifest, 0)
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

		mans = append(mans, &ModuleManifest{
			Mod: strings.TrimSuffix(d.Name(), ".update.json"),
		})

		return nil
	})
	return mans, err
}

func (c *checker) CheckManifest(man *ModuleManifest) (*ModuleStatus, error) {
	return nil, nil
}

func (c *checker) Verify(ctx context.Context, path string) ([]*ModuleStatus, error) {
	lmans, err := c.AggregateLocalManifests(path)
	if err != nil && len(lmans) < 1 {
		return nil, err
	} else if err != nil {
		logging.WarnLogger.Printf("encountered errors while looking for local module manifests: %v", err)
	}

	logging.DebugLogger.Printf("got %d local manifests", len(lmans))
	if len(lmans) > 0 {
		// TODO: check local files against local manifests
	}

	// get remote manifests
	mans, err := c.AggregateRemoteManifests()
	if err != nil && len(mans) < 1 {
		return nil, err
	} else if err != nil {
		logging.WarnLogger.Printf("encountered errors while looking for remote module manifests: %v", err)
	}

	if len(mans) > 0 {
		// TODO: check local files against remote manifests
	}

	return nil, nil
}

// func checkFile()
