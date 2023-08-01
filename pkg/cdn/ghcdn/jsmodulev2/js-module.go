package jsmodulev2

import (
	"fmt"
	"strings"

	"github.com/google/go-github/v53/github"
	"github.com/timo972/altv-cli/pkg/cdn"
	"github.com/timo972/altv-cli/pkg/cdn/ghcdn"
	"github.com/timo972/altv-cli/pkg/logging"
	"github.com/timo972/altv-cli/pkg/platform"
	"github.com/timo972/altv-cli/pkg/version"
	"golang.org/x/exp/slices"
)

func New() *ghcdn.Repository {
	return ghcdn.NewRepo(
		"altmp",
		"altv-js-module-v2",
		releaseFilter,
		assetFilter,
		manifestBuilder,
	)
}

// releaseFilter gets latest release for branch by filtering for [branch]/x
func releaseFilter(branch version.Branch, arch platform.Arch, releases []*github.RepositoryRelease) (*github.RepositoryRelease, error) {
	// sort releases by creation date, latest -> oldest
	slices.SortFunc[*github.RepositoryRelease](releases, func(a, b *github.RepositoryRelease) bool {
		// checks wether release a was created before release b
		return a.GetCreatedAt().Before(b.GetCreatedAt().Time)
	})

	var target *github.RepositoryRelease
	for _, release := range releases {
		logging.DebugLogger.Printf(release.GetName())
		if strings.Split(release.GetName(), "/")[0] == branch.String() {
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
		if !strings.Contains(asset.GetName(), arch.String()) {
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

	for _, asset := range assets {
		fp := fmt.Sprintf("modules/js-module-v2/%s", asset.GetName())
		manifest.HashList[fp] = ""
		manifest.SizeList[fp] = asset.GetSize()
	}

	return manifest, nil, nil
}
