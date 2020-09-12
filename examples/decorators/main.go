package main

import (
	"fmt"

	"decorators/container"
)

func main() {
	c := container.NewContainer()
	coolGreeter, _ := c.GetCoolGreeter()
	coolGreeter.Greet()
	fmt.Println()
}
