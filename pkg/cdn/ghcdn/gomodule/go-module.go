package gomodule

import (
	"fmt"
	"slices"
	"strings"

	"github.com/google/go-github/v53/github"
	"github.com/timo972/altv-cli/pkg/cdn"
	"github.com/timo972/altv-cli/pkg/cdn/ghcdn"
	"github.com/timo972/altv-cli/pkg/logging"
	"github.com/timo972/altv-cli/pkg/platform"
	"github.com/timo972/altv-cli/pkg/version"
)

func New() *ghcdn.Repository {
	return ghcdn.NewRepo("timo972",
		"altv-go",
		releaseFilter,
		assetFilter,
		manifestBuilder)
}

func releaseFilter(branch version.Branch, arch platform.Arch, releases []*github.RepositoryRelease) (*github.RepositoryRelease, error) {
	// sort releases by creation date, latest -> oldest
	slices.SortFunc[[]*github.RepositoryRelease](releases, func(a, b *github.RepositoryRelease) int {
		// checks wether release a was created before release b
		return int(a.GetCreatedAt().Sub(b.GetCreatedAt().Time).Seconds())
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

func assetFilter(branch version.Branch, arch platform.Arch, assets []*github.ReleaseAsset) ([]*github.ReleaseAsset, error) {
	passets := make([]*github.ReleaseAsset, 0)
	for _, asset := range assets {
		if !strings.HasSuffix(asset.GetName(), arch.SharedLibExt()) {
			continue
		}

		passets = append(passets, asset)
	}

	return passets, nil
}

func manifestBuilder(branch version.Branch, arch platform.Arch, release *github.RepositoryRelease, assets []*github.ReleaseAsset) (*cdn.Manifest, ghcdn.DownloadURLMap, error) {
	manifest := &cdn.Manifest{
		BuildNumber: -1,
		Version:     release.GetTagName(),
		HashList:    map[string]string{},
		SizeList:    map[string]int{},
	}
	urls := map[string]string{}

	for _, asset := range assets {
		fp := fmt.Sprintf("modules/go-module/%s", asset.GetName())
		manifest.HashList[fp] = ""
		manifest.SizeList[fp] = asset.GetSize()
		urls[fp] = asset.GetBrowserDownloadURL()
	}

	return manifest, urls, nil
}
