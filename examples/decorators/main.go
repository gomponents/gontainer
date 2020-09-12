package main

import (
	"fmt"
	"os"

	"decorators/container"
)

func main() {
	c := container.New()
	doer, err := c.GetChainDoer()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	doer.Do()
}
