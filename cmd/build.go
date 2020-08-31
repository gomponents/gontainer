package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/fatih/color"
	"github.com/gomponents/gontainer/pkg"
	"github.com/gomponents/gontainer/pkg/imports"
	"github.com/gomponents/gontainer/pkg/template"
	"github.com/spf13/cobra"
)

func NewBuildCmd() *cobra.Command {
	var (
		inputFiles []string
		outputFile string
		cmd        *cobra.Command
	)

	printErr := func(s string) {
		_, _ = color.New(color.BgRed, color.FgBlack).Fprint(cmd.OutOrStderr(), s)
	}

	write := func(s string) {
		_, err := cmd.OutOrStdout().Write([]byte(s))
		if err != nil {
			panic(err)
		}
	}

	handleErr := func(err error, msg string) {
		if err != nil {
			printErr(fmt.Sprintf(" %s: %s ", msg, err.Error()))
			fmt.Println()
			os.Exit(1)
		}
	}

	printHeader := func(s string) {
		_, _ = color.New(color.FgBlack, color.FgGreen).Fprint(cmd.OutOrStdout(), s)
	}

	callback := func(cmd *cobra.Command, args []string) {
		reader := pkg.NewDefaultConfigReader(func(s string) {
			write(fmt.Sprintf("    %s\n", s))
		})
		printHeader("Reading files...\n")
		input, err := reader.Read(inputFiles)
		handleErr(err, "Configuration error")

		imps := imports.NewSimpleImports()
		compiler := pkg.NewDefaultCompiler(imps)

		compiledInput, ciErr := compiler.Compile(input)
		handleErr(ciErr, "Cannot build container")

		tpl, tplErr := template.NewSimpleBuilder(imps).Build(compiledInput)
		handleErr(tplErr, "Unexpected error has occurred during building container")
		fileErr := ioutil.WriteFile(outputFile, []byte(tpl), 0644)
		handleErr(fileErr, "Error has occurred during saving file")
		write(fmt.Sprintf("Container has been built\n    %s\n", outputFile))
	}

	cmd = &cobra.Command{
		Use:   "build",
		Short: "build-container",
		Long:  "build-long",
		Run:   callback,
	}

	cmd.Flags().StringArrayVarP(&inputFiles, "input", "i", nil, "input name (required)")
	_ = cmd.MarkFlagRequired("input")

	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "output name (required)")
	_ = cmd.MarkFlagRequired("output")

	return cmd
}
