package main

import (
	"fmt"
	"time"
)

var (
	runtimeTime time.Time
)

func main() {
	c := NewContainer()
	t, _ := c.GetTime()
	fmt.Printf("Time: %s\n", t.Format(time.RFC3339))
}

func init() {
	runtimeTime = time.Now()
}
