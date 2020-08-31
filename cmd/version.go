package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	goVersion "go.hein.dev/go-version"
)

func NewVersionCmd(version, commit, date string) *cobra.Command {
	var (
		format     = "json"
		versionCmd = &cobra.Command{
			Use:   "version",
			Short: "Version will format the current build information",
			Long:  ``,
			Run: func(c *cobra.Command, _ []string) {
				c.OutOrStdout()
				_, _ = fmt.Fprint(
					c.OutOrStdout(),
					goVersion.FuncWithOutput(false, version, commit, date, format),
				)
			},
		}
	)

	versionCmd.Flags().StringVarP(&format, "format", "f", "yaml", "Output format. One of 'yaml' or 'json'.")

	return versionCmd
}
