package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	args := os.Args[1:]
	if len(args) < 4 {
		panic("Provide 3 arguments [cmd inputFile package constantName outputFile]")
	}

	inputFile := args[0]
	pkg := args[1]
	constant := args[2]
	outputFile := args[3]

	body, bodyErr := ioutil.ReadFile(inputFile)
	handleErr(bodyErr)

	output := `package {{pkg}}

var {{const}} = {{val}}
`

	output = strings.Replace(output, "{{pkg}}", pkg, 1)
	output = strings.Replace(output, "{{const}}", constant, 1)
	output = strings.Replace(output, "{{val}}", fmt.Sprintf("%+q", string(body)), 1)

	writeErr := ioutil.WriteFile(outputFile, []byte(output), 0644)
	handleErr(writeErr)
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
