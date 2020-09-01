package cmd

import (
	"github.com/spf13/cobra"
)

func NewVersionCmd(version, commit, date string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Version",
		Long:  ``,
		Run: func(c *cobra.Command, _ []string) {
			c.Printf("gontainer has version %s built from %s on %s\n", version, commit, date)
		},
	}
}
