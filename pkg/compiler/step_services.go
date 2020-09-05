package compiler

import (
	"fmt"
	"regexp"
	"sort"

	"github.com/gomponents/gontainer/pkg/arguments"
	"github.com/gomponents/gontainer/pkg/dto/compiled"
	"github.com/gomponents/gontainer/pkg/dto/input"
	"github.com/gomponents/gontainer/pkg/imports"
	"github.com/gomponents/gontainer/pkg/regex"
	"github.com/gomponents/gontainer/pkg/syntax"
)

var (
	regexServiceType        = regexp.MustCompile(`\A` + regex.ServiceType + `\z`)
	regexServiceConstructor = regexp.MustCompile(`\A` + regex.ServiceConstructor + `\z`)
)

type StepServices struct {
	aliases     imports.Aliases
	argResolver arguments.Resolver
}

func NewStepServices(aliases imports.Aliases, argResolver arguments.Resolver) *StepServices {
	return &StepServices{aliases: aliases, argResolver: argResolver}
}

func (ss StepServices) Do(i input.DTO, r *compiled.DTO) error {
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
	r.Calls, err = ss.handleServiceCalls(s.Calls)
	if err != nil {
		return compiled.Service{}, err
	}
	r.Fields, err = ss.handleServiceFields(s.Fields)
	if err != nil {
		return compiled.Service{}, err
	}
	r.Tags = ss.handleServiceTags(s.Tags)
	r.Disposable = s.Disposable
	r.Todo = false

	return r, nil
}

func (ss StepServices) handleServiceType(serviceType string) string {
	if serviceType == "" {
		return ""
	}
	_, m := regex.Match(regexServiceType, serviceType)
	t := m["type"]
	if m["import"] != "" && m["import"] != `"."` {
		t = ss.aliases.GetAlias(sanitizeImport(m["import"])) + "." + t
	}
	return m["ptr"] + t
}

func (ss StepServices) handleServiceValue(serviceValue string) string {
	if serviceValue == "" {
		return ""
	}

	return syntax.CompileServiceValue(ss.aliases, serviceValue)
}

func (ss StepServices) handleServiceConstructor(serviceConstructor string) string {
	if serviceConstructor == "" {
		return ""
	}
	_, m := regex.Match(regexServiceConstructor, serviceConstructor)
	r := ""
	if m["import"] != "" && m["import"] != `"."` {
		r = ss.aliases.GetAlias(sanitizeImport(m["import"])) + "."
	}
	return r + m["fn"]
}

func (ss StepServices) handleServiceArgs(args []interface{}) ([]compiled.Arg, error) {
	var res []compiled.Arg
	for i, a := range args {
		arg, err := ss.argResolver.Resolve(a)
		if err != nil {
			return nil, fmt.Errorf("cannot resolve arg%d: %s", i, err.Error())
		}
		res = append(res, arg)
	}
	return res, nil
}

func (ss StepServices) handleServiceCalls(calls []input.Call) ([]compiled.Call, error) {
	var res []compiled.Call
	for _, raw := range calls {
		args, err := ss.handleServiceArgs(raw.Args)
		if err != nil {
			return nil, fmt.Errorf("call `%s`: %s", raw.Method, err.Error())
		}
		call := compiled.Call{
			Method:    raw.Method,
			Args:      args,
			Immutable: raw.Immutable,
		}
		res = append(res, call)
	}
	return res, nil
}

func (ss StepServices) handleServiceFields(fields map[string]interface{}) ([]compiled.Field, error) {
	var res []compiled.Field
	var names []string
	for n, _ := range fields { //nolint:gosimple
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		f := fields[n]
		arg, err := ss.argResolver.Resolve(f)
		if err != nil {
			return nil, fmt.Errorf("field `%s`: %s", n, err.Error())
		}
		field := compiled.Field{
			Name:  n,
			Value: arg,
		}
		res = append(res, field)
	}

	sort.SliceStable(res, func(i, j int) bool {
		return res[i].Name < res[j].Name
	})

	return res, nil
}

func (ss StepServices) handleServiceTags(tags []input.Tag) (r []compiled.Tag) {
	for _, t := range tags {
		r = append(r, compiled.Tag{
			Name:     t.Name,
			Priority: t.Priority,
		})
	}
	return
}
