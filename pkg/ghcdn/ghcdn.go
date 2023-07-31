package ghcdn

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v53/github"
	"github.com/timo972/altv-cli/pkg/cdn"
	"github.com/timo972/altv-cli/pkg/logging"
	"github.com/timo972/altv-cli/pkg/platform"
	"github.com/timo972/altv-cli/pkg/version"
	"golang.org/x/exp/slices"
)

type CDN struct {
	modules ModuleMap
	client  *github.Client
}

type Repository struct {
	Name  string
	Owner string
}

type ExtendedManifest struct {
	*cdn.Manifest
	Assets map[string]*github.ReleaseAsset
}

type ModuleMap map[string]*Repository

var Default cdn.CDN = New(nil)

func New(modules ModuleMap) *CDN {
	return &CDN{
		modules: modules,
		client:  github.NewClient(nil),
	}
}

func (c *CDN) Has(module string) bool {
	_, ok := c.modules[module]
	return ok
}

func (c *CDN) Repo(module string) *Repository {
	return c.modules[module]
}

func (c *CDN) matchingRelease(branch version.Branch, arch platform.Arch, module string) (*github.RepositoryRelease, error) {
	repo := c.Repo(module)
	releases, _, err := c.client.Repositories.ListReleases(context.Background(), repo.Owner, repo.Name, nil)
	if err != nil {
		return nil, err
	}
	logging.DebugLogger.Printf("found %d releases", len(releases))

	// sort releases by creation date, latest -> oldest
	slices.SortFunc[*github.RepositoryRelease](releases, func(a, b *github.RepositoryRelease) bool {
		// checks wether release a was created before release b
		return a.GetCreatedAt().Before(b.GetCreatedAt().Time)
	})

	var target *github.RepositoryRelease
	// TODO: search for release matching branch
	for _, release := range releases {
		logging.DebugLogger.Printf(release.GetName())
		if strings.Contains(release.GetName(), branch.String()) {
			target = release
		}
	}

	if target == nil {
		return nil, fmt.Errorf("no release for branch %s found", branch)
	}

	return target, nil
}

func (c *CDN) Manifest(branch version.Branch, arch platform.Arch, module string) (*cdn.Manifest, error) {
	target, err := c.matchingRelease(branch, arch, module)
	if err != nil {
		return nil, err
	}

	manifest := &cdn.Manifest{
		BuildNumber: -1,
		Version:     target.GetTagName(),
		HashList:    map[string]string{},
		SizeList:    map[string]int{},
	}

	for _, asset := range target.Assets {
		logging.DebugLogger.Printf(asset.GetName())
		if !strings.HasSuffix(asset.GetName(), arch.SharedLibExt()) {
			continue
		}

		fp := fmt.Sprintf("modules/go-module/%s", asset.GetName())
		manifest.HashList[fp] = "n/a"
		manifest.SizeList[fp] = asset.GetSize()
	}

	return manifest, nil
}

func (c *CDN) ExtendedManifest(branch version.Branch, arch platform.Arch, module string) (*ExtendedManifest, error) {
	target, err := c.matchingRelease(branch, arch, module)
	if err != nil {
		return nil, err
	}

	manifest := &ExtendedManifest{
		Manifest: &cdn.Manifest{
			BuildNumber: -1,
			Version:     target.GetTagName(),
			HashList:    map[string]string{},
			SizeList:    map[string]int{},
		},
		Assets: map[string]*github.ReleaseAsset{},
	}

	for _, asset := range target.Assets {
		logging.DebugLogger.Printf(asset.GetName())
		if !strings.HasSuffix(asset.GetName(), arch.SharedLibExt()) {
			continue
		}

		fp := fmt.Sprintf("modules/go-module/%s", asset.GetName())
		manifest.HashList[fp] = ""
		manifest.SizeList[fp] = asset.GetSize()
		manifest.Assets[fp] = asset
	}

	return manifest, nil
}

func (c *CDN) Files(branch version.Branch, arch platform.Arch, module string, manifest bool) ([]*cdn.File, error) {
	man, err := c.ExtendedManifest(branch, arch, module)
	if err != nil {
		return nil, fmt.Errorf("unable to gather module files: %w", err)
	}

	logging.DebugLogger.Println("got manifest")

	files := make([]*cdn.File, len(man.HashList))
	i := 0
	for name, hash := range man.HashList {
		logging.DebugLogger.Printf("adding file %s", name)
		files[i] = &cdn.File{
			Name: name,
			Hash: hash,
			Url:  man.Assets[name].GetBrowserDownloadURL(),
		}
		i++
	}

	logging.DebugLogger.Println("extracted files")

	return files, nil
}
