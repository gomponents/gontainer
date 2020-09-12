package main

import (
	"decorators/container"
)

func main() {
	c := container.New()
	doer, _ := c.GetChainDoer()
	doer.Do()
}
