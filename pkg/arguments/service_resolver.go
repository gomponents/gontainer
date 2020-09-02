package arguments

import (
	"fmt"
	"github.com/gomponents/gontainer/pkg/regex"
	"regexp"

	"github.com/gomponents/gontainer/pkg/dto/compiled"
)

type ServiceResolver struct {
}

func NewServiceResolver() *ServiceResolver {
	return &ServiceResolver{}
}

func (s ServiceResolver) Resolve(v interface{}) (compiled.Arg, error) {
	expr, _ := v.(string)
	service := expr[1:]
	return compiled.Arg{
		Code:              fmt.Sprintf("container.MustGet(%+q)", service),
		Raw:               expr,
		DependsOnServices: []string{service},
	}, nil
}

var (
	serviceNameRegex = regexp.MustCompile(`\A` + regex.ServiceName + `\z`)
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
