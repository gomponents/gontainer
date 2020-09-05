package compiler

import (
	"fmt"
	"sort"

	"github.com/gomponents/gontainer/pkg/dto/compiled"
	"github.com/gomponents/gontainer/pkg/dto/input"
	"github.com/gomponents/gontainer/pkg/parameters"
)

type StepParams struct {
	paramResolver parameters.Resolver
}

func NewStepParams(paramResolver parameters.Resolver) *StepParams {
	return &StepParams{paramResolver: paramResolver}
}

func (s StepParams) Do(i input.DTO, r *compiled.DTO) error {
	var names []string
	for n, _ := range i.Params { //nolint:gosimple
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		v := i.Params[n]
		param, err := s.paramResolver.Resolve(v)
		if err != nil {
			return fmt.Errorf("cannot resolve param `%s`", n)
		}
		r.Params = append(
			r.Params,
			compiled.Param{
				Name:      n,
				Code:      param.Code,
				Raw:       param.Raw,
				DependsOn: param.DependsOn,
			},
		)
	}

	sort.SliceStable(r.Params, func(i, j int) bool {
		return r.Params[i].Name < r.Params[j].Name
	})

	return nil
}
