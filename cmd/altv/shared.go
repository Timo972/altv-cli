package main

import (
	"github.com/spf13/cobra"
	"github.com/timo972/altv-cli/pkg/platform"
)

var branch string
var arch string
var path string
var modules []string
var timeout int
var debug bool
var silent bool

func setFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&branch, "branch", "b", "release", "server version branch")
	cmd.Flags().StringVarP(&arch, "arch", "a", platform.Platform().String(), "server binary architecture")
	cmd.Flags().StringVarP(&path, "path", "p", ".", "server installation path")
	cmd.Flags().StringArrayVarP(&modules, "modules", "m", []string{"server"}, "server components to install")
	cmd.Flags().IntVarP(&timeout, "timeout", "t", -1, "server download timeout (in seconds)")
	cmd.Flags().BoolVarP(&debug, "debug", "d", false, "enable debug logging")
	cmd.Flags().BoolVarP(&silent, "silent", "s", false, "disable logging (except errors)")
}
