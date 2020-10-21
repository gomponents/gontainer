package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())

	c := NewContainer()
	// c.MustGetParam("rand") returns different value each time,
	// because of `cache_params: false`
	fmt.Println(
		c.MustGetParam("rand"),
		c.MustGetParam("rand"),
		c.MustGetParam("rand"),
	)
}
