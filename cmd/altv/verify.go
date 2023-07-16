package main

import (
	"context"
	"time"

	"github.com/spf13/cobra"
	"github.com/timo972/altv-cli/pkg/ghcdn"
	"github.com/timo972/altv-cli/pkg/logging"
	"github.com/timo972/altv-cli/pkg/platform"
	"github.com/timo972/altv-cli/pkg/vcs"
	"github.com/timo972/altv-cli/pkg/version"
)

var verifyCmd = &cobra.Command{
	Use:     "verify",
	Short:   "Verify alt:V server",
	Long:    `Verify the alt:V server in a directory.`,
	Aliases: []string{"v"},
	Run: func(cmd *cobra.Command, args []string) {
		logging.SetDebug(debug)
		if silent {
			logging.Disable()
		}

		logging.InfoLogger.Println("alt:V server verifier")

		checker := vcs.NewChecker(path, platform.Arch(arch), version.Branch(branch), modules, vcs.DefaultRegistry)
		checker.AddCDN(ghcdn.New(ghcdn.ModuleMap{
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

		if _, err := checker.Verify(ctx); err != nil {
			logging.ErrLogger.Fatalln(err)
		}

		logging.InfoLogger.Println("successfully installed")
	},
}

func init() {
	setFlags(verifyCmd)
	rootCmd.AddCommand(verifyCmd)
}
