package compiler

import (
	"regexp"
	"sort"

	"github.com/gomponents/gontainer/pkg/dto/compiled"
	"github.com/gomponents/gontainer/pkg/dto/input"
	"github.com/gomponents/gontainer/pkg/regex"
)

var (
	regexServiceType        = regexp.MustCompile("^" + regex.ServiceType + "$")
	regexServiceValue       = regexp.MustCompile("^" + regex.ServiceValue + "$")
	regexServiceConstructor = regexp.MustCompile("^" + regex.ServiceConstructor + "$")
)

type StepServices struct {
}

func (ss StepServices) Do(i *input.DTO, r *compiled.DTO) error {
	for n, s := range i.Services {
		r.Services = append(r.Services, ss.handleService(n, s))
	}
	sort.SliceStable(
		r.Services,
		func(i, j int) bool {
			return r.Services[i].Name < r.Services[j].Name
		},
	)
	return nil
}

func (ss StepServices) handleService(name string, s input.Service) compiled.Service {
	if s.Todo {
		return compiled.Service{
			Name: name,
			Todo: true,
		}
	}

	// todo
	return compiled.Service{}
}
