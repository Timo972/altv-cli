package cmd

import (
	"context"
	"time"

	"github.com/spf13/cobra"
	"github.com/timo972/altv-cli/pkg/downloader"
	"github.com/timo972/altv-cli/pkg/ghcdn"
	"github.com/timo972/altv-cli/pkg/logging"
	"github.com/timo972/altv-cli/pkg/platform"
	"github.com/timo972/altv-cli/pkg/version"
)

var branch string
var arch string
var path string
var modules []string
var timeout int
var debug bool

var installCmd = &cobra.Command{
	Use:     "install",
	Short:   "Install alt:V server",
	Long:    `Install the alt:V server into a directory.`,
	Aliases: []string{"i"},
	Run: func(cmd *cobra.Command, args []string) {
		logging.InfoLogger.Println("alt:V server installer")
		logging.SetDebug(debug)

		inst := downloader.New(path, platform.Arch(arch), version.Branch(branch), modules)
		inst.AddCDN(ghcdn.New(ghcdn.ModuleMap{
			"go-module": &ghcdn.Repository{
				Owner: "timo972",
				Name:  "altv-go",
			},
		}))

		var ctx context.Context
		var cancel context.CancelFunc
		if timeout > 0 {
			ctx, cancel = context.WithTimeout(cmd.Context(), time.Duration(timeout)*time.Second)
		} else {
			ctx, cancel = context.WithCancel(cmd.Context())
		}
		defer cancel()

		if err := inst.Download(ctx, path); err != nil {
			logging.ErrLogger.Fatalln(err)
		}

		logging.InfoLogger.Println("successfully installed")
	},
}

func init() {
	installCmd.Flags().StringVarP(&branch, "branch", "b", "release", "server version branch")
	installCmd.Flags().StringVarP(&arch, "arch", "a", platform.Platform().String(), "server binary architecture")
	installCmd.Flags().StringVarP(&path, "path", "p", ".", "server installation path")
	installCmd.Flags().StringArrayVarP(&modules, "modules", "m", []string{"server"}, "server components to install")
	installCmd.Flags().IntVarP(&timeout, "timeout", "t", -1, "server download timeout (in seconds)")
	installCmd.Flags().BoolVarP(&debug, "debug", "d", false, "enable debug logging")

	rootCmd.AddCommand(installCmd)
}
