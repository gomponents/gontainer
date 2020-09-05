package cmd

import (
	"testing"
)

func TestNewVersionCmd(t *testing.T) {
	assertCmd(
		t,
		NewVersionCmd("my-version", "my-commit", "my-date"),
		nil,
		"gontainer has version my-version built from my-commit on my-date\n",
		"",
	)
	assertCmd(
		t,
		NewVersionCmd("my-version2", "my-commit2", "my-date2"),
		nil,
		"gontainer has version my-version2 built from my-commit2 on my-date2\n",
		"",
	)
}
