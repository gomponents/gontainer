package cmd

import (
	"os"

	"github.com/gomponents/gontainer/pkg"
	"github.com/gomponents/gontainer/pkg/imports"
	"github.com/landoop/tableprinter"
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

type paramRow struct {
	Name    string `header:"Name"`
	Pattern string `header:"Param"`
}

func NewDumpParamsCmd() *cobra.Command {
	var (
		inputFiles []string
		l          uint
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

		if l < minImportLen {
			l = minImportLen
		}
		imps := &fakeImports{maxLen: int(l)}
		c := newDefaultCompiler(imps)

		compiledInput, ciErr := c.Compile(input)
		handleErr("Cannot build container", ciErr)

		var rows []paramRow
		for _, p := range compiledInput.Params {
			rows = append(rows, paramRow{
				Name:    p.Name,
				Pattern: p.Code,
			})
		}

		if len(rows) == 0 {
			cmd.Println("Could not find any parameters")
			return
		}

		p := tableprinter.New(cmd.OutOrStdout())
		p.ColumnSeparator = "│"
		p.RowSeparator = "─"
		cmd.Println()
		p.Print(rows)
	}

	cmd = &cobra.Command{
		Use:   "dump-params",
		Short: "Dump parameters",
		Long:  "",
		Run:   callback,
	}

	cmd.Flags().StringArrayVarP(&inputFiles, "input", "i", nil, "input name (required)")
	cmd.Flags().UintVarP(&l, "import-maxLen", "l", defaultImportLen, "maximum length of import path")
	_ = cmd.MarkFlagRequired("input")

	return cmd
}
