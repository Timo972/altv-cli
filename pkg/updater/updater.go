package updater

import (
	"context"

	"github.com/timo972/altv-cli/pkg/checker"
	"github.com/timo972/altv-cli/pkg/platform"
	"github.com/timo972/altv-cli/pkg/version"
)

type Updater interface {
	// Update updates the server to the latest version.
	Update(context.Context) error
}

type updater struct {
	arch    platform.Arch
	branch  version.Branch
	modules []string
	checker checker.Checker
}

func New(path string, arch platform.Arch, branch version.Branch, modules []string) Updater {
	return &updater{
		arch:    arch,
		branch:  branch,
		modules: modules,
		checker: checker.New(path, arch, branch, modules),
	}
}

func (u *updater) Update(ctx context.Context) error {
	return nil
}
