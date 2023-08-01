package main

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/timo972/altv-cli/pkg/cdn/ghcdn"
	"github.com/timo972/altv-cli/pkg/cdn/ghcdn/gomodule"
	"github.com/timo972/altv-cli/pkg/cdn/ghcdn/jsmodulev2"
	"github.com/timo972/altv-cli/pkg/platform"
	"github.com/timo972/altv-cli/pkg/util"
	"github.com/timo972/altv-cli/pkg/vcs"
)

var branch string
var arch string
var path string
var modules []string
var timeout int
var debug bool
var silent bool
var manifests bool
var github bool

func setFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&branch, "branch", "b", "release", "server version branch")
	cmd.Flags().StringVarP(&arch, "arch", "a", platform.Platform().String(), "server binary architecture")
	cmd.Flags().StringVarP(&path, "path", "p", ".", "server installation path")
	cmd.Flags().StringArrayVarP(&modules, "modules", "m", []string{"server"}, "server components to install")
	cmd.Flags().IntVarP(&timeout, "timeout", "t", -1, "server download timeout (in seconds)")
	cmd.Flags().BoolVarP(&manifests, "manifests", "M", false, "download manifests for all modules, useful to verify server files later on")
	cmd.Flags().BoolVarP(&github, "github", "g", false, "add experimental github cdn (required for js-module-v2 and go-module)")
	setLogFlags(cmd)
}

func setLogFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&debug, "debug", "d", false, "enable debug logging")
	cmd.Flags().BoolVarP(&silent, "silent", "s", false, "disable logging (except errors)")
}

func timeoutContext(ctx context.Context) (context.Context, context.CancelFunc) {
	return util.ContextWithOptionalTimeout(ctx, timeout)
}

func experimentalGithubCDN() {
	if github {
		vcs.DefaultRegistry.AddCDN(ghcdn.New(ghcdn.ModuleMap{
			"go-module":    gomodule.New(),
			"js-module-v2": jsmodulev2.New(),
		}))
	}
}
