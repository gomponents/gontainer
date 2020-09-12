package input

import (
	"fmt"
	"github.com/gomponents/gontainer/pkg/regex"
)

var (
	regexDecoratorsTag   = regex.MustCompileAz(regex.DecoratorTag)
	regexDecoratorMethod = regex.MustCompileAz(regex.DecoratorMethod)
)

// DefaultMetaValidators returns validators for DTO.Decorators.
func DefaultDecoratorsValidators() []func(DTO) error {
	return []func(DTO) error{
		ValidateDecoratorTag,
		ValidateDecoratorMethod,
		ValidateDecoratorArgs,
	}
}

func ValidateDecoratorTag(d DTO) error {
	for _, d := range d.Decorators {
		if !regexDecoratorsTag.MatchString(d.Tag) {
			return fmt.Errorf("invalid tag in definition of decorator `%s`", d.Tag)
		}
	}
	return nil
}

func ValidateDecoratorMethod(d DTO) error {
	for _, d := range d.Decorators {
		if !regexDecoratorMethod.MatchString(d.Decorator) {
			return fmt.Errorf("invalid decorator method `%s`", d.Decorator)
		}
	}
	return nil
}

func ValidateDecoratorArgs(d DTO) error {
	for _, d := range d.Decorators {
		for i, a := range d.Args {
			if !isPrimitiveType(a) {
				return fmt.Errorf("unsupported type `%T` of arg%d in decorator `%s`:`%s`", a, i, d.Tag, d.Decorator)
			}
		}
	}
	return nil
}
