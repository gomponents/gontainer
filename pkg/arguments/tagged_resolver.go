package arguments

import (
	"fmt"

	"github.com/gomponents/gontainer/pkg/dto/compiled"
	"github.com/gomponents/gontainer/pkg/regex"
)

var (
	taggedRegex = regex.MustCompileWrapped(regex.ArgTagged)
)

type TaggedResolver struct {
}

func NewTaggedResolver() *TaggedResolver {
	return &TaggedResolver{}
}

func (t TaggedResolver) Resolve(p interface{}) (compiled.Arg, error) {
	s := p.(string)
	_, m := regex.Match(taggedRegex, s)

	return compiled.Arg{
		Code:          fmt.Sprintf("container.MustGetByTag(%+q)", m["tag"]),
		Raw:           s,
		DependsOnTags: []string{m["tag"]},
	}, nil
}

func (t TaggedResolver) Supports(p interface{}) bool {
	s, ok := p.(string)
	if !ok {
		return false
	}
	return taggedRegex.MatchString(s)
}
