package main

import (
	"github.com/spf13/cobra"
	"github.com/timo972/altv-cli/pkg/ghcdn"
	"github.com/timo972/altv-cli/pkg/logging"
	"github.com/timo972/altv-cli/pkg/platform"
	"github.com/timo972/altv-cli/pkg/vcs"
	"github.com/timo972/altv-cli/pkg/version"
)

var noUpdate bool

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

		checker := vcs.NewChecker(platform.Arch(arch), version.Branch(branch), modules, vcs.DefaultRegistry)
		checker.AddCDN(ghcdn.New(ghcdn.ModuleMap{
			"go-module": &ghcdn.Repository{
				Owner: "timo972",
				Name:  "altv-go",
			},
		}))

		ctx, cancel := timeoutContext(cmd.Context())
		defer cancel()

		if _, err := checker.Verify(ctx, path); err != nil {
			logging.ErrLogger.Fatalln(err)
		}

		logging.InfoLogger.Println("server files valid")
	},
}

func init() {
	setFlags(verifyCmd)
	verifyCmd.Flags().BoolVarP(&noUpdate, "no-update", "n", false, "do not check for updates, just verify files")
	rootCmd.AddCommand(verifyCmd)
}
