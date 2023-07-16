package server

import (
	"github.com/timo972/altv-cli/pkg/platform"
	"github.com/timo972/altv-cli/pkg/version"
)

type Updater interface {
	// Update updates the server to the latest version.
	Update() error
}

type updater struct {
	arch    platform.Arch
	branch  version.Branch
	include []string
	checker Checker
}

func NewUpdater(path string, arch platform.Arch, branch version.Branch, include []string, checker Checker) Updater {
	return &updater{
		arch:    arch,
		branch:  branch,
		include: include,
		checker: checker,
	}
}

func (u *updater) Update() error {
	return nil
}
