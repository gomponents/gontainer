package main

import (
	"decorators/container"
)

func main() {
	c := container.NewContainer()
	doer, _ := c.GetChainDoer()
	doer.Do()
}
