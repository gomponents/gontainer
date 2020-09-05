package runner

import (
	"github.com/gomponents/gontainer/pkg/dto/compiled"
	"github.com/gomponents/gontainer/pkg/dto/input"
)

// todo

type Runner interface {
	Run() error
}

type Step interface {
	Run(*input.DTO, *compiled.DTO) error
}

type ChainStepRunner struct {
	steps []Step
}

func NewChainStepRunner(steps ...Step) *ChainStepRunner {
	return &ChainStepRunner{steps: steps}
}

func (c ChainStepRunner) Run() error {
	i := input.CreateDefaultDTO()
	o := compiled.DTO{}

	for _, s := range c.steps {
		if err := s.Run(&i, &o); err != nil {
			return err
		}
	}

	return nil
}

func NewBuildRunner() *ChainStepRunner {
	return NewChainStepRunner(
		ReadConfigStep{},
		// CompileConfigStep{},
	)
}

type ReadConfigStep struct {
	//patterns []string
}

func (r ReadConfigStep) Run(i *input.DTO, _ *compiled.DTO) error {
	panic("implement me")
}
