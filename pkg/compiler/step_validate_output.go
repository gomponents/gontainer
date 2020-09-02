package compiler

import (
	"github.com/gomponents/gontainer/pkg/dto/compiled"
	"github.com/gomponents/gontainer/pkg/dto/input"
)

type StepValidateOutput struct {
	validator CompiledValidator
}

func NewStepValidateOutput(validator CompiledValidator) *StepValidateOutput {
	return &StepValidateOutput{validator: validator}
}

func (s StepValidateOutput) Do(_ *input.DTO, r *compiled.DTO) error {
	return s.validator.Validate(*r)
}
