package main

import (
	"fmt"
)

type Person struct {
	Name string
	Age  int
}

func main() {
	c := NewContainer()
	p, err := c.GetPerson()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s is %d years old\n", p.Name, p.Age)
}
