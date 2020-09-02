package compiler

import (
	"github.com/gomponents/gontainer/pkg/dto/compiled"
	"github.com/gomponents/gontainer/pkg/dto/input"
)

type Step interface {
	Do(*input.DTO, *compiled.DTO) error
}

type Compiler2 struct {
	steps []Step
}

func NewCompiler2(steps ...Step) *Compiler2 {
	return &Compiler2{steps: steps}
}

func (c Compiler2) Compile(i input.DTO) (compiled.DTO, error) {
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

type ImportAliases interface {
	GetAlias(string) string
}

type ImportPrefixes interface {
	RegisterPrefix(shortcut string, path string) error
}

type Tokenizer interface {
	RegisterFunction(goImport string, goFunc string, tokenFunc string)
}

type CompiledValidator interface {
	Validate(compiled.DTO) error
}
