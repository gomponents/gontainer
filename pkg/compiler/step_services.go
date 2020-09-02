package compiler

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/gomponents/gontainer/pkg/dto/compiled"
	"github.com/gomponents/gontainer/pkg/dto/input"
	"github.com/gomponents/gontainer/pkg/regex"
)

var (
	regexServiceType        = regexp.MustCompile(`\A` + regex.ServiceType + `\z`)
	regexServiceValue       = regexp.MustCompile(`\A` + regex.ServiceValue + `\z`)
	regexServiceConstructor = regexp.MustCompile(`\A` + regex.ServiceConstructor + `\z`)
)

type StepServices struct {
	aliases     ImportAliases
	argResolver ArgResolver
}

func (ss StepServices) Do(i *input.DTO, r *compiled.DTO) error {
	var names []string
	for n, _ := range i.Services { //nolint:gosimple
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		s := i.Services[n]
		c, err := ss.handleService(n, s)
		if err != nil {
			return fmt.Errorf("service `%s`: %s", n, err.Error())
		}
		r.Services = append(r.Services, c)
	}
	sort.SliceStable(
		r.Services,
		func(i, j int) bool {
			return r.Services[i].Name < r.Services[j].Name
		},
	)
	return nil
}

func (ss StepServices) handleService(name string, s input.Service) (compiled.Service, error) {
	if s.Todo {
		return compiled.Service{
			Name: name,
			Todo: true,
		}, nil
	}

	var err error

	r := compiled.Service{}
	r.Name = name
	r.Getter = s.Getter
	r.Type = ss.handleServiceType(s.Type)
	r.Value = ss.handleServiceValue(s.Value)
	r.Constructor = ss.handleServiceConstructor(s.Constructor)
	r.Args, err = ss.handleServiceArgs(s.Args)
	if err != nil {
		return compiled.Service{}, err
	}

	return r, nil
}

func (ss StepServices) handleServiceType(serviceType string) string {
	if serviceType == "" {
		return ""
	}
	_, m := regex.Match(regexServiceType, serviceType)
	t := m["type"]
	if m["import"] != "" {
		t = ss.aliases.GetAlias(sanitizeImport(m["import"])) + "." + t
	}
	return m["ptr"] + t
}

func (ss StepServices) handleServiceValue(serviceValue string) string {
	if serviceValue == "" {
		return ""
	}

	_, m := regex.Match(regexServiceValue, serviceValue)

	if m["v1"] != "" {
		parts := make([]string, 0)
		if m["import"] != "" {
			parts = append(parts, ss.aliases.GetAlias(sanitizeImport(m["import"])))
		}
		if m["struct"] != "" {
			parts = append(parts, m["struct"]+"{}")
		}
		return strings.Join(append(parts, m["value"]), ".")
	}

	parts := make([]string, 0)
	if m["import2"] != "" {
		parts = append(parts, ss.aliases.GetAlias(sanitizeImport(m["import2"])))
	}
	return m["ptr2"] + strings.Join(append(parts, m["struct2"]), ".") + "{}"
}

func (ss StepServices) handleServiceConstructor(serviceConstructor string) string {
	if serviceConstructor == "" {
		return ""
	}
	_, m := regex.Match(regexServiceConstructor, serviceConstructor)
	r := ""
	if m["import"] != "" {
		r = ss.aliases.GetAlias(sanitizeImport(m["import"])) + "."
	}
	return r + m["fn"]
}

func (ss StepServices) handleServiceArgs(args []interface{}) ([]compiled.Arg, error) {
	var res []compiled.Arg
	for i, a := range args {
		arg, err := ss.argResolver.Resolve(a)
		if err != nil {
			return nil, fmt.Errorf("cannot resolve arg%d", i)
		}
		res = append(res, arg)
	}
	return res, nil
}
