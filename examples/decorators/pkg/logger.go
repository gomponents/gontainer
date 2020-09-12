package pkg

import (
	"fmt"
)

type Logger interface {
	Log(...interface{})
}

type BasicLogger struct{}

func (b BasicLogger) Log(a ...interface{}) {
	fmt.Println(a...)
}

func ServiceLogger(name string, service interface{}, logger Logger) interface{} {
	logger.Log(fmt.Sprintf("[DEBUG] Created service `%s` of type `%T`", name, service))
	return service
}
