package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"

	"github.com/gomponents/gontainer/pkg/dto/compiled"
	"github.com/gomponents/gontainer/pkg/dto/input"
	"github.com/gomponents/gontainer/pkg/imports"
	"github.com/gomponents/gontainer/pkg/template"
)

type templatePrinter struct {
	printer
	aliases    imports.Aliases
	collection imports.Collection
	outputFile string
}

func newTemplatePrinter(
	w io.Writer,
	aliases imports.Aliases,
	collection imports.Collection,
	outputFile string,
) *templatePrinter {
	return &templatePrinter{
		printer:    printer{w: w},
		aliases:    aliases,
		collection: collection,
		outputFile: outputFile,
	}
}

func (t templatePrinter) run(_ *input.DTO, o *compiled.DTO) error {
	tpl, tplErr := template.NewBuilder(t.aliases, t.collection).Build(*o)
	if tplErr != nil {
		return fmt.Errorf("unexpected error: %s", tplErr)
	}
	of := filepath.Clean(t.outputFile)

	if err := ioutil.WriteFile(of, []byte(tpl), 0644); err != nil {
		return err
	}

	t.printf("Successfully generated container to file: `%s`\n", of)

	return nil
}
