package cmd

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_printer_printf(t *testing.T) {
	t.Run("Given error", func(t *testing.T) {
		defer func() {
			r := recover()
			if !assert.NotEmpty(t, r) {
				return
			}
			assert.Implements(t, (*error)(nil), r)
			assert.EqualError(t, r.(error), "my error")
		}()

		p := printer{mockWriter{err: fmt.Errorf("my error")}}
		p.printf("hello")
	})
}

type mockWriter struct {
	n   int
	err error
}

func (m mockWriter) Write([]byte) (n int, err error) {
	return m.n, m.err
}
