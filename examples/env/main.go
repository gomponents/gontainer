package main

import (
	"fmt"
	"os"
)

type Person struct {
	Name       string
	Age        int
	EyesColor  string
	MotherName string
}

func (p *Person) SetEyesColor(s string) {
	p.EyesColor = s
}

// WithMotherName returns new pointer. We must to return pointer to have the same type as already declared,
// see services.person.value in gontainer.yml
//
// v := &Person{}
// v = v.WithMotherName("Lily")
func (p Person) WithMotherName(s string) *Person {
	p.MotherName = s
	return &p
}

func main() {
	c := NewContainer()
	p, err := c.GetPerson()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}
	fmt.Printf("%s is %d years old; %s eyes; mother's name is %s\n", p.Name, p.Age, p.EyesColor, p.MotherName)
}
