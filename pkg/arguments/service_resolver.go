package arguments

import (
	"fmt"

	"github.com/gomponents/gontainer/pkg/dto/compiled"
	"github.com/gomponents/gontainer/pkg/regex"
)

type ServiceResolver struct {
}

func NewServiceResolver() *ServiceResolver {
	return &ServiceResolver{}
}

func (s ServiceResolver) Resolve(v interface{}) (compiled.Arg, error) {
	expr := v.(string)
	service := expr[1:]
	return compiled.Arg{
		Code:              fmt.Sprintf("container.MustGet(%+q)", service),
		Raw:               expr,
		DependsOnServices: []string{service},
	}, nil
}

var (
	serviceNameRegex = regex.MustCompileWrapped(regex.ServiceName)
)

func (s ServiceResolver) Supports(v interface{}) bool {
	expr, ok := v.(string)
	if !ok {
		return false
	}

	if len(expr) < 2 || []rune(expr)[0] != '@' {
		return false
	}

	return serviceNameRegex.MatchString(expr[1:])
}
