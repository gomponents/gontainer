package compiler

import (
	"github.com/gomponents/gontainer/pkg/dto/compiled"
	"github.com/gomponents/gontainer/pkg/dto/input"
)

type Step interface {
	Do(*input.DTO, *compiled.DTO) error
}

type Compiler struct {
	steps []Step
}

func NewCompiler2(steps ...Step) *Compiler {
	return &Compiler{steps: steps}
}

func (c Compiler) Compile(i input.DTO) (compiled.DTO, error) {
	r := compiled.DTO{}
	cpInput := i.Clone()
	for _, s := range c.steps {
		if err := s.Do(&cpInput, &r); err != nil {
			return compiled.DTO{}, err
		}
	}
	return r, nil
}

type InputValidator interface {
	Validate(input.DTO) error
}

// todo rename it
type ImportAliases interface {
	GetAlias(string) string
}

// todo rename it
type ImportPrefixes interface {
	RegisterPrefix(shortcut string, path string) error
}

type Tokenizer interface {
	RegisterFunction(goImport string, goFunc string, tokenFunc string)
}

type CompiledValidator interface {
	Validate(compiled.DTO) error
}

type ArgResolver interface {
	Resolve(interface{}) (compiled.Arg, error)
}
