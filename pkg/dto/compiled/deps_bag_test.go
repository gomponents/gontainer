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

	find := func(id string) []string {
		_, deps := finder.find(id)
		return deps
	}

	assert.Equal(t, []string{"bar", "foobar", "foo", "bar"}, find("bar"))
	assert.Equal(t, []string{"foo", "bar", "foobar", "foo"}, find("foo"))
	assert.Empty(t, find("hey"))
	assert.Empty(t, find("test"))
}
