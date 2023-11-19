package main

import (
	"log"
	"strings"

	"github.com/spf13/cobra"
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

		experimentalGithubCDN()
		checker := vcs.NewChecker(platform.Arch(arch), version.Branch(branch), modules, vcs.DefaultRegistry)

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
	logger.Printf("integrity / version summary")
	logger.Printf(head, "Module", "Integrity", "[✅|💥|⭕]", "Version", "[✅|🔼|⭕]")
	logger.Printf("|%s|", strings.Repeat("-", 66))
	for mod, stat := range status {
		emojis := statusToEmoji(stat)
		logger.Printf(row, mod, emojis[0], emojis[1])
	}

}

func statusToEmoji(status vcs.ModuleStatus) [2]string {
	str := [2]string{}

	if status.Has(vcs.StatusValid) {
		str[0] = "✅"
	} else if status.Has(vcs.StatusInvalid) {
		str[0] = "💥"
	} else {
		str[0] = "⭕"
	}

	if status.Has(vcs.StatusUpToDate) {
		str[1] = "✅"
	} else if status.Has(vcs.StatusUpgradable) {
		str[1] = "🔼"
	} else {
		str[1] = "⭕"
	}

	return str
}
