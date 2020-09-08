package cmd

import (
	"fmt"

	"github.com/gomponents/gontainer/pkg/consts"
	"github.com/spf13/cobra"
)

func NewGetUpdateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get-update",
		Short: "Prints command to update dependencies",
		Long:  "Gontainer requires additional dependencies to run.\nTo update them, run `gontainer get-update | bash`.",
		Run: func(cmd *cobra.Command, _ []string) {
			cmd.Print(fmt.Sprintf("go get -u %s\n", consts.GontainerHelperPath))
		},
	}
}
