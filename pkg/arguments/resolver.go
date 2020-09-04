package arguments

import (
	"fmt"

	"github.com/gomponents/gontainer/pkg/dto/compiled"
	"github.com/gomponents/gontainer/pkg/imports"
	"github.com/gomponents/gontainer/pkg/parameters"
)

type Resolver interface {
	Resolve(interface{}) (compiled.Arg, error)
}

type SubResolver interface {
	Resolve(interface{}) (compiled.Arg, error)
	Supports(interface{}) bool
}

type ChainResolver struct {
	subResolvers []SubResolver
}

func NewChainResolver(subResolvers ...SubResolver) *ChainResolver {
	return &ChainResolver{subResolvers: subResolvers}
}

func (s ChainResolver) Resolve(i interface{}) (compiled.Arg, error) {
	for _, r := range s.subResolvers {
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

func NewDefaultResolver(r parameters.Resolver, a imports.Aliases) *ChainResolver {
	return NewChainResolver(
		NewServiceResolver(),
		NewTaggedResolver(),
		NewValueResolver(a),
		NewParamResolver(r),
	)
}
