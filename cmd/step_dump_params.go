package cmd

import (
	"github.com/gomponents/gontainer/pkg/dto/compiled"
	"github.com/gomponents/gontainer/pkg/dto/input"
	"github.com/landoop/tableprinter"
	"io"
)

type paramRow struct {
	Name    string `header:"Name"`
	Pattern string `header:"Param"`
}

type dumpParams struct {
	printer
}

func newDumpParams(w io.Writer) *dumpParams {
	return &dumpParams{printer: printer{w: w}}
}

func (d dumpParams) run(_ *input.DTO, o *compiled.DTO) error {
	var rows []paramRow
	for _, p := range o.Params {
		rows = append(rows, paramRow{
			Name:    p.Name,
			Pattern: p.Code,
		})
	}

	if len(rows) == 0 {
		d.printf("Could not find any parameters\n")
		return nil
	}

	p := tableprinter.New(d.printer.w)
	p.ColumnSeparator = "│"
	p.RowSeparator = "─"
	d.printf("\n")
	p.Print(rows)

	return nil
}
