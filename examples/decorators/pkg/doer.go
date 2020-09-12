package pkg

import (
	"fmt"
)

type Doer interface {
	Do()
}

type ChainDoer struct {
	doers []Doer
}

func (c ChainDoer) Do() {
	for _, d := range c.doers {
		d.Do()
	}
}

type SimpleDoer struct {
	i int
}

func NewSimpleDoer(i int) *SimpleDoer {
	return &SimpleDoer{i: i}
}

func (s SimpleDoer) Do() {
	fmt.Printf("doing %d...\n", s.i)
}

type DecoratedSimpleDoer struct {
	*SimpleDoer
}

func NewDecoratedSimpleDoer(_ string, simpleDoer *SimpleDoer) *DecoratedSimpleDoer {
	return &DecoratedSimpleDoer{SimpleDoer: simpleDoer}
}

func (d DecoratedSimpleDoer) Do() {
	fmt.Print("[decorated] ")
	d.SimpleDoer.Do()
}
