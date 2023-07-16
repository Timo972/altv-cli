package vcs

import (
	"context"

	"github.com/timo972/altv-cli/pkg/cdn"
	"github.com/timo972/altv-cli/pkg/logging"
	"github.com/timo972/altv-cli/pkg/platform"
	"github.com/timo972/altv-cli/pkg/version"
)

type Updater interface {
	AddCDN(cdn.CDN)
	Update(context.Context) error
}

type updater struct {
	check Checker
	dl    Downloader
}

func NewUpdater(path string, arch platform.Arch, branch version.Branch, modules []string, reg CDNRegistry) Updater {
	u := &updater{
		check: NewChecker(path, arch, branch, modules, reg),
		dl:    NewDownloader(path, arch, branch, modules, reg),
	}

	return u
}

func (u *updater) AddCDN(cdn cdn.CDN) {
	u.check.AddCDN(cdn)
	u.dl.AddCDN(cdn)
}

func (u *updater) Update(ctx context.Context) error {
	status, err := u.check.Verify(ctx)
	if err != nil {
		return err
	}

	logging.DebugLogger.Printf("status: %v", status)

	return nil
}
