package compiler

import (
	"fmt"

	"github.com/gomponents/gontainer/pkg/arguments"
	"github.com/gomponents/gontainer/pkg/dto/compiled"
	"github.com/gomponents/gontainer/pkg/dto/input"
	"github.com/gomponents/gontainer/pkg/imports"
	"github.com/gomponents/gontainer/pkg/regex"
)

var (
	regexDecoratorMethod = regex.MustCompileAz(regex.DecoratorMethod)
)

type StepDecorators struct {
	aliases     imports.Aliases
	argResolver arguments.Resolver
}

func NewStepDecorators(aliases imports.Aliases, argResolver arguments.Resolver) *StepDecorators {
	return &StepDecorators{aliases: aliases, argResolver: argResolver}
}

func (s StepDecorators) Do(i input.DTO, r *compiled.DTO) error {
	for _, d := range i.Decorators {
		n, err := s.compileDecorator(d)
		if err != nil {
			return err
		}
		r.Decorators = append(r.Decorators, n)
	}
	return nil
}

func (s StepDecorators) compileDecorator(d input.Decorator) (compiled.Decorator, error) {
	r := compiled.Decorator{}
	r.Tag = d.Tag

	_, m := regex.Match(regexDecoratorMethod, d.Decorator)
	method := m["fn"]
	if m["import"] != "" && m["import"] != `"."` {
		method = s.aliases.GetAlias(m["import"]) + "." + method
	}
	r.Decorator = method

	for i, a := range d.Args {
		resolved, err := s.argResolver.Resolve(a)
		if err != nil {
			return compiled.Decorator{}, fmt.Errorf("decorator `%s`:`%s`: cannot resolve arg%d: %s", d.Tag, d.Decorator, i, err.Error())
		}
		r.Args = append(r.Args, resolved)
	}

	return r, nil
}
