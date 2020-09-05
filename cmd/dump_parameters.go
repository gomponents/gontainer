package cmd

import (
	"github.com/gomponents/gontainer/pkg/imports"
	"github.com/spf13/cobra"
)

const (
	minImportLen     = 10
	defaultImportLen = 15
)

type fakeImports struct {
	maxLen int
}

func (f fakeImports) GetAlias(i string) string {
	const preffix = "(...)"
	r := []rune(i)
	if len(r) > f.maxLen {
		r = r[len(r)-(f.maxLen-len([]rune(preffix))):]
		i = preffix + string(r)
	}
	return "\"" + i + "\""
}

func (f fakeImports) GetImports() []imports.Import {
	return nil
}

func (f fakeImports) RegisterPrefix(shortcut string, path string) error {
	return nil
}

func NewDumpParamsCmd() *cobra.Command {
	var (
		inputFiles []string
		l          uint
		cmd        *cobra.Command
	)

	cmd = &cobra.Command{
		Use:   "dump-params",
		Short: "Dump parameters",
		Long:  "",
		RunE: func(cmd *cobra.Command, args []string) error {
			if l < minImportLen {
				l = minImportLen
			}
			imps := &fakeImports{maxLen: int(l)}
			runner := newStepRunner(
				newReadConfig(cmd.OutOrStdout(), inputFiles),
				newCompile(newDefaultCompiler(imps)),
				newDumpParams(cmd.OutOrStdout()),
			)
			return runner.run()
		},
	}

	cmd.Flags().StringArrayVarP(&inputFiles, "input", "i", nil, "input name (required)")
	cmd.Flags().UintVarP(&l, "import-maxLen", "l", defaultImportLen, "maximum length of import path")
	_ = cmd.MarkFlagRequired("input")

	return cmd
}
