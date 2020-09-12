package pkg

import (
	"fmt"
)

type Greeter struct{}

func (Greeter) Greet() {
	fmt.Print("Hello")
}

type CoolGreeter struct {
	*Greeter
	name string
}

func (cg CoolGreeter) Greet() {
	cg.Greeter.Greet()
	fmt.Printf(", bro! My name is %s.", cg.name)
}

func NewCoolGreeter(name string, greeter *Greeter) *CoolGreeter {
	return &CoolGreeter{
		name:    name,
		Greeter: greeter,
	}
}
