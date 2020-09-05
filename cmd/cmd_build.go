package cmd

import (
	"github.com/gomponents/gontainer/pkg/imports"
	"github.com/spf13/cobra"
)

func NewBuildCmd() *cobra.Command {
	var (
		inputFiles []string
		outputFile string
		cmd        *cobra.Command
	)

	cmd = &cobra.Command{
		Use:   "build",
		Short: "Build container",
		Long:  "",
		RunE: func(cmd *cobra.Command, args []string) error {
			imps := imports.NewSimpleImports()
			runner := newStepRunner(
				newReadConfig(cmd.OutOrStdout(), inputFiles),
				newCompile(newDefaultCompiler(imps)),
				newTemplatePrinter(imps, imps, outputFile),
			)
			return runner.run()
		},
	}

	cmd.Flags().StringArrayVarP(&inputFiles, "input", "i", nil, "input name (required)")
	_ = cmd.MarkFlagRequired("input")

	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "output name (required)")
	_ = cmd.MarkFlagRequired("output")

	return cmd
}
