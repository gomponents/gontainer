package compiler

import (
	"github.com/gomponents/gontainer/pkg/dto/compiled"
	"github.com/gomponents/gontainer/pkg/dto/input"
)

type StepValidateInput struct {
	validator input.Validator
}

func NewStepValidateInput(validator input.Validator) *StepValidateInput {
	return &StepValidateInput{validator: validator}
}

func (s StepValidateInput) Do(i input.DTO, _ *compiled.DTO) error {
	return s.validator.Validate(i)
}
