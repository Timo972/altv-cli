package main

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/timo972/altv-cli/pkg/platform"
	"github.com/timo972/altv-cli/pkg/util"
)

var branch string
var arch string
var path string
var modules []string
var timeout int
var debug bool
var silent bool
var manifests bool

func setFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&branch, "branch", "b", "release", "server version branch")
	cmd.Flags().StringVarP(&arch, "arch", "a", platform.Platform().String(), "server binary architecture")
	cmd.Flags().StringVarP(&path, "path", "p", ".", "server installation path")
	cmd.Flags().StringArrayVarP(&modules, "modules", "m", []string{"server"}, "server components to install")
	cmd.Flags().IntVarP(&timeout, "timeout", "t", -1, "server download timeout (in seconds)")
	cmd.Flags().BoolVarP(&debug, "debug", "d", false, "enable debug logging")
	cmd.Flags().BoolVarP(&silent, "silent", "s", false, "disable logging (except errors)")
	cmd.Flags().BoolVarP(&manifests, "manifests", "M", false, "download manifests for all modules, useful to verify server files later on")
}

func timeoutContext(ctx context.Context) (context.Context, context.CancelFunc) {
	return util.ContextWithOptionalTimeout(ctx, timeout)
}
