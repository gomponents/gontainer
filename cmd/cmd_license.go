package cmd

import (
	"github.com/spf13/cobra"
)

//go:generate go run ../embed-file/main.go ../LICENSE cmd license license_var.go

func NewLicenseCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "license",
		Short: "License",
		Run: func(cmd *cobra.Command, _ []string) {
			cmd.Print(license)
		},
	}
}
