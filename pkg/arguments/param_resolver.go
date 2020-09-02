package arguments

import (
	"github.com/gomponents/gontainer/pkg/dto/compiled"
	"github.com/gomponents/gontainer/pkg/parameters"
)

type ParamResolver struct {
	resolver parameters.Resolver
}

func NewParamResolver(resolver parameters.Resolver) *ParamResolver {
	return &ParamResolver{resolver: resolver}
}

func (p ParamResolver) Resolve(v interface{}) (compiled.Arg, error) {
	param, err := p.resolver.Resolve(v)
	return compiled.Arg{
		Code:            param.Code,
		Raw:             param.Raw,
		DependsOnParams: param.DependsOn,
	}, err
}

func (p ParamResolver) Supports(interface{}) bool {
	return true
}
