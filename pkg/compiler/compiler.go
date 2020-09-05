package compiler

import (
	"github.com/gomponents/gontainer/pkg/dto/compiled"
	"github.com/gomponents/gontainer/pkg/dto/input"
)

type Step interface {
	Do(input.DTO, *compiled.DTO) error
}

type Compiler struct {
	steps []Step
}

func NewCompiler(steps ...Step) *Compiler {
	return &Compiler{steps: steps}
}

func (c Compiler) Compile(i input.DTO) (compiled.DTO, error) {
	r := compiled.DTO{}
	for _, s := range c.steps {
		if err := s.Do(i, &r); err != nil {
			return compiled.DTO{}, err
		}
	}
	return r, nil
}
