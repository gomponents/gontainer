package compiled

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_circularDepFinder_doFind(t *testing.T) {
	a := map[string][]string{
		"foo":    {"bar"},
		"bar":    {"foobar"},
		"foobar": {"hey", "foo"},
		"hey":    {},
	}

	finder := newCircularDepFinder(func(id string) []string {
		deps, _ := a[id] //nolint:gosimple
		return deps
	})

	assert.Equal(t, []string{"bar", "foobar", "foo", "bar"}, finder.find("bar"))
	assert.Equal(t, []string{"foo", "bar", "foobar", "foo"}, finder.find("foo"))
	assert.Nil(t, finder.find("hey"))
	assert.Nil(t, finder.find("test"))
}
