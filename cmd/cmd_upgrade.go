package cmd

import (
	"fmt"

	"github.com/gomponents/gontainer/pkg/consts"
	"github.com/spf13/cobra"
)

func NewGetUpgradeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get-upgrade",
		Short: "Prints command to update dependencies",
		Run: func(cmd *cobra.Command, _ []string) {
			cmd.Print(fmt.Sprintf("go get -u %s\n", consts.GontainerHelperPath))
		},
	}
}
