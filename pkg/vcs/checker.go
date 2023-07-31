package vcs

import (
	"context"

	"github.com/timo972/altv-cli/pkg/cdn"
	"github.com/timo972/altv-cli/pkg/platform"
	"github.com/timo972/altv-cli/pkg/version"
)

type ModuleStatus struct{}

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

func (c *checker) Verify(ctx context.Context, path string) ([]*ModuleStatus, error) {
	/*status := make(chan *ModuleStatus, len(c.modules))

	for _, mod := range c.modules {

	}*/

	return nil, nil
}

// func checkFile()
