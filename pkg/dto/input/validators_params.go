package input

import (
	"fmt"
	"sort"

	"github.com/gomponents/gontainer/pkg/regex"
)

var (
	regexParamName = regex.MustCompileAz(regex.ParamName)
)

// DefaultParamsValidators returns validators for DTO.Params.
func DefaultParamsValidators() []func(DTO) error {
	return []func(DTO) error{
		ValidateParams,
	}
}

func ValidateParams(d DTO) error {
	var names []string
	for k, _ := range d.Params { //nolint:gosimple
		names = append(names, k)
	}
	sort.Strings(names)

	for _, n := range names {
		if !regexParamName.MatchString(n) {
			return fmt.Errorf("invalid parameter name `%s`", n)
		}

		v := d.Params[n]

		if !isPrimitiveType(v) {
			return fmt.Errorf("unsupported type `%T` of parameter `%s`", v, n)
		}
	}
	return nil
}
