package arguments

import (
	"github.com/gomponents/gontainer/pkg/dto/compiled"
	"github.com/gomponents/gontainer/pkg/regex"
	"github.com/gomponents/gontainer/pkg/syntax"
	"regexp"
)

var (
	argRegex = regexp.MustCompile(`\A` + regex.ArgValue + `\z`)
)

// todo rename it
type ImportAliases interface {
	GetAlias(string) string
}

type ValueResolver struct {
	aliases ImportAliases
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
