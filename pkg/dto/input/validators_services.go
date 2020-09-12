package input

import (
	"fmt"
	"reflect"
	"regexp"
	"sort"

	"github.com/gomponents/gontainer-helpers/container"
	"github.com/gomponents/gontainer/pkg/regex"
)

var (
	regexServiceName        = regex.MustCompileAz(regex.ServiceName)
	regexServiceGetter      = regex.MustCompileAz(regex.ServiceGetter)
	regexServiceType        = regex.MustCompileAz(regex.ServiceType)
	regexServiceValue       = regex.MustCompileAz(regex.ServiceValue)
	regexServiceConstructor = regex.MustCompileAz(regex.ServiceConstructor)
	regexServiceCallName    = regex.MustCompileAz(regex.ServiceCallName)
	regexServiceFieldName   = regex.MustCompileAz(regex.ServiceFieldName)
	regexServiceTag         = regex.MustCompileAz(regex.ServiceTag)
)

type ValidateService func(Service) error

// DefaultServicesValidators returns validators for DTO.Services.
func DefaultServicesValidators() []func(DTO) error {
	return []func(DTO) error{
		ValidateServices,
	}
}

func ValidateServices(d DTO) error {
	validators := []ValidateService{
		ValidateConstructorType,
		ValidateServiceGetter,
		ValidateServiceType,
		ValidateServiceValue,
		ValidateServiceConstructor,
		ValidateServiceArgs,
		ValidateServiceCalls,
		ValidateServiceFields,
		ValidateServiceTags,
	}

	for n, s := range d.Services {
		if err := ValidateServiceName(n); err != nil {
			return err
		}
		if s.Todo {
			continue
		}
		for _, v := range validators {
			err := v(s)
			if err == nil {
				continue
			}
			return fmt.Errorf("service `%s`: %s", n, err.Error())
		}
	}
	return nil
}

func ValidateServiceName(n string) error {
	if !regexServiceName.MatchString(n) {
		return fmt.Errorf(
			"invalid service name `%s`",
			n,
		)
	}
	return nil
}

func ValidateConstructorType(s Service) error {
	if s.Constructor == "" && s.Value == "" {
		return fmt.Errorf("missing constructor or value")
	}

	if s.Constructor != "" && s.Value != "" {
		return fmt.Errorf("cannot define constructor and value together")
	}

	if len(s.Args) > 0 && s.Constructor == "" {
		return fmt.Errorf("arguments are not empty, but constructor is missing")
	}

	return nil
}

func ValidateServiceGetter(s Service) error {
	reserved := []string{"ValidateAllServices"}
	c := struct {
		*container.BaseContainer
		*container.BaseParamContainer
		*container.BaseTaggedContainer
	}{}
	r := reflect.TypeOf(c)
	for i := 0; i < r.NumMethod(); i++ {
		reserved = append(reserved, r.Method(i).Name)
	}

	for _, r := range reserved {
		if s.Getter == r {
			return fmt.Errorf("invalid getter, `%s` is reserved", r)
		}
	}

	if s.Getter != "" && s.Type == "" {
		return fmt.Errorf("getter is given, but type is missing")
	}
	return validateRegexField("getter", s.Getter, regexServiceGetter, true)
}

func ValidateServiceType(s Service) error {
	if s.Type != "" && s.Getter == "" {
		return fmt.Errorf("type is given, but getter is missing")
	}
	return validateRegexField("type", s.Type, regexServiceType, true)
}

func ValidateServiceValue(s Service) error {
	return validateRegexField("value", s.Value, regexServiceValue, true)
}

func ValidateServiceConstructor(s Service) error {
	return validateRegexField("constructor", s.Constructor, regexServiceConstructor, true)
}

func ValidateServiceArgs(s Service) error {
	// todo one common method for validating all args
	for i, a := range s.Args {
		if !isPrimitiveType(a) {
			return fmt.Errorf("unsupported type `%T` of arg%d", a, i)
		}
	}
	return nil
}

func ValidateServiceCalls(s Service) error {
	for j, c := range s.Calls {
		if err := validateRegexField(fmt.Sprintf("method name (call %d)", j), c.Method, regexServiceCallName, false); err != nil {
			return err
		}
		for i, a := range c.Args {
			if !isPrimitiveType(a) {
				return fmt.Errorf("unsupported type `%T` of (call `%s`, arg %d)", a, c.Method, i)
			}
		}
	}
	return nil
}

func ValidateServiceFields(s Service) error {
	var names []string
	for n, _ := range s.Fields { //nolint:gosimple
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		v := s.Fields[n]
		if err := validateRegexField("field", n, regexServiceFieldName, false); err != nil {
			return err
		}
		if !isPrimitiveType(v) {
			return fmt.Errorf("unsupported type `%T` of field `%s`", v, n)
		}
	}
	return nil
}

func ValidateServiceTags(s Service) error {
	exists := make(map[string]bool)
	for _, t := range s.Tags {
		if err := validateRegexField("tag", t.Name, regexServiceTag, false); err != nil {
			return err
		}
		if _, ok := exists[t.Name]; ok {
			return fmt.Errorf("duplicate tag `%s`", t.Name)
		}
		exists[t.Name] = true
	}
	return nil
}

func validateRegexField(field string, value string, expr *regexp.Regexp, optional bool) error {
	if optional && value == "" {
		return nil
	}
	if !expr.MatchString(value) {
		return fmt.Errorf("invalid %s `%s`", field, value)
	}
	return nil
}
