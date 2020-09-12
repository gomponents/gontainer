package compiler

import (
	"github.com/gomponents/gontainer/pkg/dto/compiled"
	"github.com/gomponents/gontainer/pkg/dto/input"
	"github.com/gomponents/gontainer/pkg/imports"
)

type StepDecorators struct {
	aliases imports.Aliases
}

func NewStepDecorators(aliases imports.Aliases) *StepDecorators {
	return &StepDecorators{aliases: aliases}
}

func (s StepDecorators) Do(_ input.DTO, r *compiled.DTO) error {
	// todo
	return nil
}
