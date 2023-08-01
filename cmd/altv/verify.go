package main

import (
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/timo972/altv-cli/pkg/cdn/ghcdn"
	"github.com/timo972/altv-cli/pkg/cdn/ghcdn/gomodule"
	"github.com/timo972/altv-cli/pkg/cdn/ghcdn/jsmodulev2"
	"github.com/timo972/altv-cli/pkg/logging"
	"github.com/timo972/altv-cli/pkg/platform"
	"github.com/timo972/altv-cli/pkg/vcs"
	"github.com/timo972/altv-cli/pkg/version"
)

var noUpdate bool

var verifyCmd = &cobra.Command{
	Use:     "verify",
	Short:   "Verify alt:V server",
	Long:    "Verify the alt:V server in a directory.",
	Aliases: []string{"v"},
	Run: func(cmd *cobra.Command, args []string) {
		logging.SetDebug(debug)
		if silent {
			logging.Disable()
		}

		logging.InfoLogger.Println("alt:V server verifier")

		checker := vcs.NewChecker(platform.Arch(arch), version.Branch(branch), modules, vcs.DefaultRegistry)
		checker.AddCDN(ghcdn.New(ghcdn.ModuleMap{
			"go-module":    gomodule.New(),
			"js-module-v2": jsmodulev2.New(),
		}))

		ctx, cancel := timeoutContext(cmd.Context())
		defer cancel()

		status, err := checker.Verify(ctx, path, !noUpdate)
		if err != nil {
			logging.ErrLogger.Fatalln(err)
		}

		printSummary(logging.InfoLogger, status)
	},
}

func init() {
	setFlags(verifyCmd)
	verifyCmd.Flags().BoolVarP(&noUpdate, "no-update", "n", false, "do not check for updates, just verify files")
	rootCmd.AddCommand(verifyCmd)
}

func printSummary(logger *log.Logger, status vcs.ModuleStatusResult) {
	head := "| %-18s | %-9s %1s | %-9s %1s |"
	row := "| %-18s | %-19s | %-19s |"
	logger.Printf("status summary")
	logger.Printf(head, "Module", "Integrity", "[✅|💥|⭕]", "Version", "[✅|🔼|⭕]")
	logger.Printf("|%s|", strings.Repeat("-", 66))
	for mod, stat := range status {
		emojis := statusToEmoji(stat)
		logger.Printf(row, mod, emojis[0], emojis[1])
	}

}

func statusToEmoji(status vcs.ModuleStatus) [2]string {
	switch status {
	case vcs.StatusInvalid:
		return [2]string{"💥", "⭕"}
	case vcs.StatusInvalidUpgradable:
		return [2]string{"💥", "🔼"}
	case vcs.StatusInvalidUpToDate:
		return [2]string{"💥", "✅"}
	case vcs.StatusValid:
		return [2]string{"✅", "⭕"}
	case vcs.StatusValidUpgradable:
		return [2]string{"✅", "🔼"}
	case vcs.StatusValidUpToDate:
		return [2]string{"✅", "✅"}
	case vcs.StatusUpToDate:
		return [2]string{"⭕", "✅"}
	case vcs.StatusUpgradable:
		return [2]string{"⭕", "🔼"}
	default:
		return [2]string{"⭕", "⭕"}
	}
}
