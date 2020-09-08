package cmd

import (
	"fmt"
	"testing"

	"github.com/gomponents/gontainer/pkg/consts"
)

func TestNewGetUpdateCmd(t *testing.T) {
	assertCmd(
		t,
		NewGetUpdateCmd(),
		nil,
		fmt.Sprintf("go get -u %s\n", consts.GontainerHelperPath),
		"",
	)
}
