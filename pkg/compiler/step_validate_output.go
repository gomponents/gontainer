package compiler

import (
	"github.com/gomponents/gontainer/pkg/dto/compiled"
	"github.com/gomponents/gontainer/pkg/dto/input"
)

type StepValidateOutput struct {
	validator compiled.Validator
}

func NewStepValidateOutput(validator compiled.Validator) *StepValidateOutput {
	return &StepValidateOutput{validator: validator}
}

func (s StepValidateOutput) Do(_ *input.DTO, r *compiled.DTO) error {
	return s.validator.Validate(*r)
}
