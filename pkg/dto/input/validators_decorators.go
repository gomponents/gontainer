package input

import (
	"fmt"
	"github.com/gomponents/gontainer/pkg/regex"
)

var (
	regexDecoratorsTag = regex.MustCompileAz(regex.DecoratorTag)
)

// DefaultMetaValidators returns validators for DTO.Decorators.
func DefaultDecoratorsValidators() []func(DTO) error {
	return []func(DTO) error{
		ValidateDecoratorTag,
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
