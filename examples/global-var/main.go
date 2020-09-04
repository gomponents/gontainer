package main

import (
	"fmt"
	"time"
)

var (
	runtimeTime time.Time
)

type mysqlTime struct {
	*time.Time
}

func newMysqlTime(time *time.Time) *mysqlTime {
	return &mysqlTime{Time: time}
}

func (m mysqlTime) MysqlString() string {
	return m.In(time.UTC).Format("2006-01-02 15:04:05")
}

func main() {
	c := NewContainer()
	t, _ := c.GetTime()
	fmt.Printf("Time: %s\n", t.Format(time.RFC3339))

	m, _ := c.GetMysqlTime()
	fmt.Println(m.MysqlString())
}

func init() {
	runtimeTime = time.Now()
}
