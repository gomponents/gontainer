package arguments

import (
	"github.com/gomponents/gontainer/pkg/dto/compiled"
	"github.com/gomponents/gontainer/pkg/imports"
	"github.com/gomponents/gontainer/pkg/regex"
	"github.com/gomponents/gontainer/pkg/syntax"
)

var (
	argRegex = regex.MustCompileWrapped(regex.ArgValue)
)

type ValueResolver struct {
	aliases imports.Aliases
}

func NewValueResolver(aliases imports.Aliases) *ValueResolver {
	return &ValueResolver{aliases: aliases}
}

func (v ValueResolver) Resolve(p interface{}) (compiled.Arg, error) {
	s := p.(string)
	_, m := regex.Match(argRegex, s)

	return compiled.Arg{
		Code: syntax.CompileServiceValue(v.aliases, m["argval"]),
		Raw:  s,
	}, nil
}

func (v ValueResolver) Supports(p interface{}) bool {
	s, ok := p.(string)
	if !ok {
		return false
	}
	return argRegex.MatchString(s)
}
