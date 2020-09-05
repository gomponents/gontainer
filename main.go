package main

import (
	"os"

	"github.com/gomponents/gontainer/cmd"
	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	rootCmd := cobra.Command{
		Use:   "gontainer",
		Short: "Gontainer is a missing DI container for GO",
		Long: `Gontainer allows you to build DI container based on provided configuration.
Re-use dependencies whenever you need and forget about dependency hell in main.go.`,
		Version:      version,
		SilenceUsage: true,
	}

	rootCmd.AddCommand(
		cmd.NewBuildCmd(),
		cmd.NewDumpParamsCmd(),
		cmd.NewLicenseCmd(),
		cmd.NewVersionCmd(version, commit, date),
	)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
