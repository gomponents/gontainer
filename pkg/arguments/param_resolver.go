package arguments

import (
	"github.com/gomponents/gontainer/pkg/dto/compiled"
	"github.com/gomponents/gontainer/pkg/parameters"
)

type RawParamResolver interface {
	Resolve(interface{}) (parameters.Expr, error)
}

type ParamResolver struct {
	resolver RawParamResolver
}

func NewParamResolver(resolver RawParamResolver) *ParamResolver {
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
