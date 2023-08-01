package vcs

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/timo972/altv-cli/pkg/cdn"
	"github.com/timo972/altv-cli/pkg/cdn/altcdn"
	"github.com/timo972/altv-cli/pkg/logging"
	"github.com/timo972/altv-cli/pkg/platform"
	"github.com/timo972/altv-cli/pkg/version"
)

type Downloader interface {
	Download(ctx context.Context, path string, manifests bool) error
}

type downloader struct {
	CDNRegistry
	cdns    []cdn.CDN
	arch    platform.Arch
	branch  version.Branch
	modules []string
}

func NewDownloader(arch platform.Arch, branch version.Branch, modules []string, registry CDNRegistry) Downloader {
	return &downloader{
		CDNRegistry: registry,
		arch:        arch,
		branch:      branch,
		modules:     modules,
		cdns:        []cdn.CDN{altcdn.Default},
	}
}

// TODO: utilize goroutines to aggregate files simultaneously
func (d *downloader) AggregateFiles(manifests bool) []*cdn.File {
	allFiles := make([]*cdn.File, 0)
	for _, module := range d.modules {
		cdn, ok := d.moduleCDN(module)
		if !ok {
			logging.WarnLogger.Printf(newErrNoCDN(module).Error())
			continue
		}
		logging.DebugLogger.Printf("cdn %v for module %s", cdn, module)

		files, err := cdn.Files(d.branch, d.arch, module, manifests)
		if err != nil {
			logging.WarnLogger.Printf("no files for module %s found, skipping: %v", module, err)
			continue
		}

		logging.DebugLogger.Printf("%d files for module %s", len(files), module)
		allFiles = append(allFiles, files...)
	}

	return allFiles
}

// DownloadFiles downloads all files from the given slice of files to the given path concurrently
func (d *downloader) DownloadFiles(ctx context.Context, path string, files []*cdn.File) error {
	// spin up a goroutine for each file download process
	errs := make(chan error, len(files))
	for _, file := range files {
		go downloadFile(errs, path, file)
	}

	for range files {
		select {
		case err := <-errs:
			logging.DebugLogger.Printf("got resp: %v", err)
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return nil
}

// Download aggregates all files from the given modules and downloads them to the given path
func (d *downloader) Download(ctx context.Context, path string, manifests bool) error {
	files := d.AggregateFiles(manifests)
	logging.InfoLogger.Printf("downloading %d files", len(files))
	return d.DownloadFiles(ctx, path, files)
}

// downloadFile is a utility to download the given file to the given path and verify its checksum
func downloadFile(c chan error, p string, file *cdn.File) {
	resp, err := http.DefaultClient.Get(file.Url)
	if err != nil {
		c <- err
		return
	}
	defer resp.Body.Close()
	logging.DebugLogger.Printf("got response: %s", resp.Status)

	if resp.StatusCode != http.StatusOK {
		c <- fmt.Errorf("unexpected statusCode at download of%s: %s", file.Name, resp.Status)
		return
	}

	logging.DebugLogger.Printf("opening file %s", file.Name)

	fol := fmt.Sprintf("%s/%s", p, path.Dir(file.Name))
	logging.DebugLogger.Printf("file folder: %s", p)

	if _, err := os.Stat(fol); os.IsNotExist(err) {
		if err = os.MkdirAll(fol, 0700); err != nil {
			c <- fmt.Errorf("can not create directory %s: %w", file.Name, err)
			return
		}
	}

	f, err := os.OpenFile(fmt.Sprintf("%s/%s", p, file.Name), os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		c <- fmt.Errorf("can not open file %s: %w", file.Name, err)
		return
	}
	defer f.Close()

	logging.DebugLogger.Printf("writing file %s", file.Name)

	h := sha1.New()
	bodyReader := io.TeeReader(resp.Body, h)

	if _, err = io.Copy(f, bodyReader); err != nil {
		c <- fmt.Errorf("can not write file %s: %w", file.Name, err)
		return
	}
	logging.DebugLogger.Printf("wrote file %s, checking checksum", file.Name)

	if file.Hash == "" {
		if file.Type != cdn.ModuleManifestFile {
			logging.WarnLogger.Printf("no checksum for %s, be careful!", file.Name)
		}
		c <- nil
		return
	}

	checksum := hex.EncodeToString(h.Sum(nil))
	if checksum != file.Hash {
		c <- fmt.Errorf("checksum mismatch for %s: expected %s, got %s; be careful! file might be corrupted", file.Name, file.Hash, checksum)
		return
	}

	logging.DebugLogger.Printf("checksum for %s is ok", file.Name)

	c <- nil
}
