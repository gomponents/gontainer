package compiled

import (
	"fmt"
	"sort"
	"strings"
)

func DefaultParamsValidators() []func(DTO) error {
	return []func(DTO) error{
		ValidateParamsReqParamsExist,
		ValidateParamsCircularDeps,
	}
}

func ValidateParamsReqParamsExist(d DTO) error {
	list := make(map[string]bool)
	for _, p := range d.Params {
		list[p.Name] = true
	}
	for _, p := range d.Params {
		for _, n := range p.DependsOn {
			if _, ok := list[n]; !ok {
				return fmt.Errorf("parameter `%s` requires param `%s`, but it does not exist", p.Name, n)
			}
		}
	}
	for _, s := range d.Services {
		for _, a := range getAllServiceArgs(s) {
			for _, n := range a.DependsOnParams {
				if _, ok := list[n]; !ok {
					return fmt.Errorf("service `%s` requires param `%s`, but it does not exist", s.Name, n)
				}
			}
		}
	}
	return nil
}

func ValidateParamsCircularDeps(d DTO) error {
	mapping := make(map[string]Param)
	//nolint:staticcheck
	params := append(d.Params)
	for _, p := range d.Params {
		mapping[p.Name] = p
	}
	sort.SliceStable(params, func(i, j int) bool {
		return params[i].Name < params[j].Name
	})

	finder := newCircularDepFinder(func(id string) []string {
		//nolint:gosimple
		param, _ := mapping[id]
		deps := param.DependsOn
		sort.Strings(deps)
		return deps
	})

	for _, p := range params {
		if deps := finder.find(p.Name); deps != nil {
			return fmt.Errorf("circular dependency in params: %s", strings.Join(deps, " -> "))
		}
	}

	return nil
}
