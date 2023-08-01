package ghcdn

import (
	"context"
	"fmt"

	"github.com/google/go-github/v53/github"
	"github.com/timo972/altv-cli/pkg/cdn"
	"github.com/timo972/altv-cli/pkg/logging"
	"github.com/timo972/altv-cli/pkg/platform"
	"github.com/timo972/altv-cli/pkg/version"
)

type CDN struct {
	modules ModuleMap
	client  *github.Client
}

type ReleaseFilter func(version.Branch, platform.Arch, []*github.RepositoryRelease) (*github.RepositoryRelease, error)
type AssetFilter func(version.Branch, platform.Arch, []*github.ReleaseAsset) ([]*github.ReleaseAsset, error)
type DownloadURLMap map[string]string
type ManifestBuilder func(version.Branch, platform.Arch, *github.RepositoryRelease, []*github.ReleaseAsset) (*cdn.Manifest, DownloadURLMap, error)

type Repository struct {
	Name            string
	Owner           string
	ReleaseFilter   ReleaseFilter
	AssetFilter     AssetFilter
	ManifestBuilder ManifestBuilder
}

type ExtendedManifest struct {
	*cdn.Manifest
	Assets map[string]*github.ReleaseAsset
}

type ModuleMap map[string]*Repository

func New(modules ModuleMap) *CDN {
	return &CDN{
		modules: modules,
		client:  github.NewClient(nil),
	}
}

func NewRepo(owner string, name string, releaseFilter ReleaseFilter, assetFilter AssetFilter, manBuilder ManifestBuilder) *Repository {
	return &Repository{
		Name:            name,
		Owner:           owner,
		ReleaseFilter:   releaseFilter,
		AssetFilter:     assetFilter,
		ManifestBuilder: manBuilder,
	}
}

func (c *CDN) Has(module string) bool {
	_, ok := c.modules[module]
	return ok
}

func (c *CDN) Repo(module string) *Repository {
	return c.modules[module]
}

func (c *CDN) matchingRelease(branch version.Branch, arch platform.Arch, module string) (*Repository, *github.RepositoryRelease, error) {
	repo := c.Repo(module)
	releases, _, err := c.client.Repositories.ListReleases(context.Background(), repo.Owner, repo.Name, nil)
	if err != nil {
		return nil, nil, err
	}
	logging.DebugLogger.Printf("found %d releases", len(releases))

	release, err := repo.ReleaseFilter(branch, arch, releases)
	return repo, release, err
}

func (c *CDN) buildManifest(branch version.Branch, arch platform.Arch, module string) (*cdn.Manifest, DownloadURLMap, error) {
	repo, target, err := c.matchingRelease(branch, arch, module)
	if err != nil {
		return nil, nil, err
	}

	assets, err := repo.AssetFilter(branch, arch, target.Assets)
	if err != nil {
		return nil, nil, err
	}

	return repo.ManifestBuilder(branch, arch, target, assets)
}

func (c *CDN) Manifest(branch version.Branch, arch platform.Arch, module string) (*cdn.Manifest, error) {
	manifest, _, err := c.buildManifest(branch, arch, module)
	return manifest, err
}

func (c *CDN) Files(branch version.Branch, arch platform.Arch, module string, manifest bool) ([]*cdn.File, error) {
	man, urls, err := c.buildManifest(branch, arch, module)
	if err != nil {
		return nil, fmt.Errorf("unable to build module manifest (github): %w", err)
	}

	logging.DebugLogger.Println("got manifest")

	files := make([]*cdn.File, len(man.HashList))
	i := 0
	for name, hash := range man.HashList {
		logging.DebugLogger.Printf("adding file %s", name)
		files[i] = &cdn.File{
			Type: cdn.ModuleFile,
			Name: name,
			Hash: hash,
			Url:  urls[name],
		}
		i++
	}

	logging.DebugLogger.Println("extracted files")

	return files, nil
}
