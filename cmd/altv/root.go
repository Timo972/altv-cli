package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/timo972/altv-cli/pkg/logging"
)

var rootCmd = &cobra.Command{
	Use:   "altv",
	Short: "alt:V command line tool",
	Long:  `A blazingly fast alt:V server manager cli written in go.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logging.ErrLogger.Fatalln(err)
		os.Exit(1)
	}
}
