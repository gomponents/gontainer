package cmd

import (
	"testing"
)

func TestNewLicenseCmd(t *testing.T) {
	assertCmd(
		t,
		NewLicenseCmd(),
		nil,
		license,
		"",
	)
}
