package compiled

import (
	"fmt"
	"sort"
	"strings"
)

func DefaultServicesValidators() []func(DTO) error {
	return []func(DTO) error{
		ValidateServicesReqServicesExist,
		ValidateServicesCircularDeps,
	}
}

func ValidateServicesReqServicesExist(d DTO) error {
	list := make(map[string]bool)
	for _, s := range d.Services {
		list[s.Name] = true
	}
	for _, s := range d.Services {
		for _, a := range getAllServiceArgs(s) {
			for _, n := range a.DependsOnServices {
				if _, ok := list[n]; !ok {
					return fmt.Errorf("service `%s` requires service `%s`, but it does not exist", s.Name, n)
				}
			}
		}
	}
	return nil
}

func ValidateServicesCircularDeps(d DTO) error {
	taggedMapping := make(map[string][]string)
	serviceMapping := make(map[string]Service)
	for _, s := range d.Services {
		serviceMapping[s.Name] = s
		for _, t := range s.Tags {
			if _, ok := taggedMapping[t.Name]; !ok {
				taggedMapping[t.Name] = make([]string, 0)
			}
			taggedMapping[t.Name] = append(taggedMapping[t.Name], s.Name)
		}
	}

	depsByArg := func(a Arg) []string {
		var res []string
		res = append(res, a.DependsOnServices...)
		for _, t := range a.DependsOnTags {
			deps, _ := taggedMapping[t] //nolint:gosimple
			res = append(res, deps...)
		}
		return res
	}

	depsByService := func(s Service) []string {
		var r []string
		for _, a := range getAllServiceArgs(s) {
			r = append(r, depsByArg(a)...)
		}
		return r
	}

	finder := newCircularDepFinder(func(id string) []string {
		r := depsByService(serviceMapping[id])
		sort.Strings(r)
		return r
	})

	services := append(d.Services) //nolint:staticcheck
	sort.Slice(services, func(i, j int) bool {
		return services[i].Name < services[j].Name
	})

	for _, s := range services {
		if circular, deps := finder.find(s.Name); circular {
			return fmt.Errorf("circular dependency in services: %s", strings.Join(deps, " -> "))
		}
	}

	return nil
}
