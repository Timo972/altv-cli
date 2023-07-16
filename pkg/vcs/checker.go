package vcs

import (
	"context"

	"github.com/timo972/altv-cli/pkg/platform"
	"github.com/timo972/altv-cli/pkg/version"
)

type ModuleStatus struct{}

type Checker interface {
	CDNRegistry
	Verify(context.Context) ([]*ModuleStatus, error)
}

type checker struct {
	CDNRegistry
}

func NewChecker(path string, arch platform.Arch, branch version.Branch, modules []string, registry CDNRegistry) Checker {
	return &checker{
		CDNRegistry: registry,
	}
}

func (c *checker) Verify(ctx context.Context) ([]*ModuleStatus, error) {
	return nil, nil
}
