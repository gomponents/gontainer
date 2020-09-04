package cmd

import (
	"io/ioutil"
	"os"
	"path/filepath"

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

	handleErr := func(h string, err error) {
		if err == nil {
			return
		}
		cmd.PrintErrf("%s: %s\n", h, err.Error())
		os.Exit(1)
	}

	callback := func(cmd *cobra.Command, args []string) {
		reader := pkg.NewDefaultConfigReader(func(s string) {
			cmd.Printf("    %s\n", s)
		})
		cmd.Printf("Reading files...\n")
		input, rErr := reader.Read(inputFiles)
		handleErr("Configuration error", rErr)

		imps := imports.NewSimpleImports()
		c := newDefaultCompiler(imps)

		compiledInput, ciErr := c.Compile(input)
		handleErr("Cannot build container", ciErr)

		tpl, tplErr := template.NewSimpleBuilder(imps, imps).Build(compiledInput)
		handleErr("Unexpected error has occurred during building container", tplErr)
		of := filepath.Clean(outputFile)
		fileErr := ioutil.WriteFile(of, []byte(tpl), 0644)
		handleErr("Error has occurred during saving file", fileErr)
		cmd.Printf("Container has been built: %s\n", of)
	}

	cmd = &cobra.Command{
		Use:   "build",
		Short: "Build container",
		Long:  "",
		Run:   callback,
	}

	cmd.Flags().StringArrayVarP(&inputFiles, "input", "i", nil, "input name (required)")
	_ = cmd.MarkFlagRequired("input")

	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "output name (required)")
	_ = cmd.MarkFlagRequired("output")

	return cmd
}
