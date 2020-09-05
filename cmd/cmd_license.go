package cmd

import (
	"github.com/spf13/cobra"
)

//go:generate go run ../embed-file/main.go ../LICENSE cmd const license license_const.go

func NewLicenseCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "license",
		Short: "License",
		Run: func(cmd *cobra.Command, _ []string) {
			cmd.Print(license)
		},
	}
}
