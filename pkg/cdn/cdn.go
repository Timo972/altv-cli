package cdn

import (
	"io"

	"github.com/timo972/altv-cli/pkg/platform"
	"github.com/timo972/altv-cli/pkg/version"
)

type FileType uint8

const (
	ModuleFile FileType = iota
	ModuleManifestFile
)

type CDN interface {
	// Has checks wether the CDN hosts the given module files.
	Has(module string) bool
	// Manifest returns the manifest for the given branch, arch and module.
	Manifest(branch version.Branch, arch platform.Arch, module string) (*Manifest, error)
	// Files returns a list of required files for the given branch, arch and module.
	Files(branch version.Branch, arch platform.Arch, module string, manifest bool) ([]*File, error)
}

type Manifest struct {
	BuildNumber int               `json:"latestBuildNumber"`
	Version     string            `json:"version"`
	HashList    map[string]string `json:"hashList"`
	SizeList    map[string]int    `json:"sizeList"`
}

type File struct {
	Type FileType
	Name string
	Url  string
	Hash string
}

type BuiltFile struct {
	Name    string
	Content io.Reader
}
