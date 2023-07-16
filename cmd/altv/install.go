package main

import (
	"context"
	"time"

	"github.com/spf13/cobra"
	"github.com/timo972/altv-cli/pkg/logging"
	"github.com/timo972/altv-cli/pkg/platform"
	"github.com/timo972/altv-cli/pkg/vcs"
	"github.com/timo972/altv-cli/pkg/version"
)

var installCmd = &cobra.Command{
	Use:     "install",
	Short:   "Install alt:V server",
	Long:    `Install the alt:V server into a directory.`,
	Aliases: []string{"i"},
	Run: func(cmd *cobra.Command, args []string) {
		logging.SetDebug(debug)
		if silent {
			logging.Disable()
		}

		logging.InfoLogger.Println("alt:V server installer")

		inst := vcs.NewDownloader(path, platform.Arch(arch), version.Branch(branch), modules, vcs.DefaultRegistry)

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
	setFlags(installCmd)
	rootCmd.AddCommand(installCmd)
}
