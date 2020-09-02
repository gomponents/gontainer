package arguments

import (
	"fmt"

	"github.com/gomponents/gontainer/pkg/dto/compiled"
	"github.com/gomponents/gontainer/pkg/parameters"
)

type Subresolver interface {
	Resolve(interface{}) (compiled.Arg, error)
	Supports(interface{}) bool
}

type Resolver struct {
	subresolvers []Subresolver
}

func NewResolver(subresolvers ...Subresolver) *Resolver {
	return &Resolver{subresolvers: subresolvers}
}

func (s Resolver) Resolve(i interface{}) (compiled.Arg, error) {
	for _, r := range s.subresolvers {
		if r.Supports(i) {
			result, err := r.Resolve(i)
			if err == nil {
				result.Raw = i
			}

			return result, err
		}
	}

	return compiled.Arg{}, fmt.Errorf("cannot resolve argument `%s`", i)
}

func NewDefaultResolver(resolver parameters.Resolver) *Resolver {
	return NewResolver(
		NewServiceResolver(),
		NewTaggedResolver(),
		NewParamResolver(resolver),
	)
}
