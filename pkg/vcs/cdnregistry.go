package vcs

import (
	"github.com/timo972/altv-cli/pkg/altcdn"
	"github.com/timo972/altv-cli/pkg/cdn"
	"github.com/timo972/altv-cli/pkg/ghcdn"
)

var DefaultRegistry CDNRegistry = NewRegistry(altcdn.Default, ghcdn.Default)

type CDNRegistry interface {
	AddCDN(cdn.CDN)
	moduleCDN(module string) (cdn.CDN, bool)
}

type cdnRegistry struct {
	cdns []cdn.CDN
}

func NewRegistry(cdns ...cdn.CDN) CDNRegistry {
	return &cdnRegistry{
		cdns: cdns,
	}
}

// getCDN returns the CDN that hosts the given module.
func (r *cdnRegistry) moduleCDN(module string) (cdn.CDN, bool) {
	for _, cdn := range r.cdns {
		if cdn.Has(module) {
			return cdn, true
		}
	}
	return nil, false
}

func (r *cdnRegistry) AddCDN(cdn cdn.CDN) {
	r.cdns = append(r.cdns, cdn)
}
