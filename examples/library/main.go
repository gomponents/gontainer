package main

import (
	"os"

	"github.com/landoop/tableprinter"
	"library/container"
	"library/pkg"
)

//go:generate gontainer build -i container/gontainer.yml -o container/container.go

func main() {
	c := container.NewLibraryContainer()
	l, _ := c.GetLibrary()
	printLibrary(l)
}

type bookRow struct {
	Title  string `header:"Title"`
	Author string `header:"Author"`
}

func printLibrary(library pkg.Library) {
	p := tableprinter.New(os.Stdout)

	var rows []bookRow
	for _, b := range library.Books {
		rows = append(
			rows,
			bookRow{
				Title:  b.Title,
				Author: b.Author.Name,
			},
		)
	}

	p.Print(rows)
}
