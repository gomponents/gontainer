package cmd

import (
	"fmt"
	"io"

	"github.com/gomponents/gontainer/pkg/dto/compiled"
	"github.com/gomponents/gontainer/pkg/dto/input"
)

const indent = "    "

type step interface {
	run(*input.DTO, *compiled.DTO) error
}

type stepRunner struct {
	steps []step
}

func newStepRunner(steps ...step) *stepRunner {
	return &stepRunner{steps: steps}
}

func (s stepRunner) run() error {
	i := input.CreateDefaultDTO()
	o := compiled.DTO{}
	for _, sr := range s.steps {
		if err := sr.run(&i, &o); err != nil {
			return err
		}
	}
	return nil
}

type printer struct {
	w io.Writer
}

func (p printer) printf(f string, a ...interface{}) {
	_, err := p.w.Write([]byte(fmt.Sprintf(f, a...)))
	if err != nil {
		panic(err)
	}
}
