package cmd

import (
	"fmt"
	"testing"

	"github.com/gomponents/gontainer/pkg/consts"
)

func TestNewGetUpgradeCmd(t *testing.T) {
	assertCmd(
		t,
		NewGetUpgradeCmd(),
		nil,
		fmt.Sprintf("go get -u %s\n", consts.GontainerHelperPath),
		"",
	)
}
