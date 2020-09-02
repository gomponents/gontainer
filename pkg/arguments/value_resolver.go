package arguments

import (
	"regexp"

	"github.com/gomponents/gontainer/pkg/dto/compiled"
	"github.com/gomponents/gontainer/pkg/regex"
)

var (
	argRegex = regexp.MustCompile(`\A` + regex.ArgValue + `\z`)
)

type ValueResolver struct {
}

func (v ValueResolver) Resolve(interface{}) (compiled.Arg, error) {
	// todo
	panic("implement me")
}

func (v ValueResolver) Supports(p interface{}) bool {
	s, ok := p.(string)
	if !ok {
		return false
	}
	return argRegex.MatchString(s)
}
