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
	Update(ctx context.Context, path string) error
}

type updater struct {
	check Checker
	dl    Downloader
	reg   CDNRegistry
}

func NewUpdater(arch platform.Arch, branch version.Branch, modules []string, reg CDNRegistry) Updater {
	u := &updater{
		check: NewChecker(arch, branch, modules, reg),
		dl:    NewDownloader(arch, branch, modules, reg),
		reg:   reg,
	}

	return u
}

func (u *updater) AddCDN(cdn cdn.CDN) {
	u.reg.AddCDN(cdn)
}

func (u *updater) Update(ctx context.Context, path string) error {
	status, err := u.check.Verify(ctx, path)
	if err != nil {
		return err
	}

	logging.DebugLogger.Printf("status: %v", status)

	return nil
}
