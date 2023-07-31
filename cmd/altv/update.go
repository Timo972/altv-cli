package main

import (
	"github.com/spf13/cobra"
	"github.com/timo972/altv-cli/pkg/logging"
	"github.com/timo972/altv-cli/pkg/platform"
	"github.com/timo972/altv-cli/pkg/vcs"
	"github.com/timo972/altv-cli/pkg/version"
)

var updateCmd = &cobra.Command{
	Use:     "update",
	Short:   "Update alt:V server",
	Long:    `Update the alt:V server a directory.`,
	Aliases: []string{"u"},
	Run: func(cmd *cobra.Command, args []string) {
		logging.SetDebug(debug)
		if silent {
			logging.Disable()
		}

		logging.InfoLogger.Println("alt:V server updater")

		upd := vcs.NewUpdater(platform.Arch(arch), version.Branch(branch), modules, vcs.DefaultRegistry)

		ctx, cancel := timeoutContext(cmd.Context())
		defer cancel()

		if err := upd.Update(ctx, path); err != nil {
			logging.ErrLogger.Fatalln(err)
		}

		logging.InfoLogger.Println("successfully installed")
	},
}

func init() {
	setFlags(updateCmd)
	rootCmd.AddCommand(updateCmd)
}
