package vcs

import (
	"fmt"
)

type errNoCDN struct {
	mod string
}

func (e *errNoCDN) Error() string {
	return fmt.Sprintf("no cdn for module %s found, skipping", e.mod)
}

func newErrNoCDN(mod string) error {
	return &errNoCDN{mod: mod}
}

type errNoManifest struct {
	mod string
	e   error
}

func (e *errNoManifest) Error() string {
	return fmt.Sprintf("no manifest for module %s found: %v", e.mod, e.e)
}

func newErrNoManifest(mod string, e error) error {
	return &errNoManifest{mod: mod, e: e}
}
