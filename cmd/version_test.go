package cmd

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewVersionCmd(t *testing.T) {
	cmd := NewVersionCmd("my-version", "my-commit", "my-date")
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	if !assert.NoError(t, cmd.Execute()) {
		return
	}
	contents, readErr := ioutil.ReadAll(b)
	if !assert.NoError(t, readErr) {
		return
	}
	assert.Equal(t, "gontainer has version my-version built from my-commit on my-date\n", string(contents))
}
