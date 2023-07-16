package checker

import (
	"context"

	"github.com/timo972/altv-cli/pkg/platform"
	"github.com/timo972/altv-cli/pkg/version"
)

type ModuleStatus struct{}

type Checker interface {
	Validate(context.Context) ([]*ModuleStatus, error)
}

type checker struct{}

func New(path string, arch platform.Arch, branch version.Branch, modules []string) Checker {
	return &checker{}
}

func (c *checker) Validate(ctx context.Context) ([]*ModuleStatus, error) {
	return nil, nil
}
