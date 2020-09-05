package imports

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const viperPkg = "github.com/spf13/viper"

func TestSimpleImports_RegisterPrefix(t *testing.T) {
	t.Run("Given scenario", func(t *testing.T) {
		i := NewSimpleImports()
		shortcut := "viper"
		expectedViper := "i0_spf13_viper"
		expectedRemote := "i1_viper_remote"

		if !assert.NoError(t, i.RegisterPrefix(shortcut, viperPkg)) {
			return
		}

		// aliases for "viper" and "github.com/spf13/viper" should be equal
		assert.Equal(t, expectedViper, i.GetAlias(shortcut))
		assert.Equal(t, expectedViper, i.GetAlias(viperPkg))

		// aliases for "viper/remote" and "github.com/spf13/viper/remote" should be equal
		assert.Equal(t, expectedRemote, i.GetAlias(shortcut+"/remote"))
		assert.Equal(t, expectedRemote, i.GetAlias(viperPkg+"/remote"))
	})

	t.Run("Given error", func(t *testing.T) {
		i := NewSimpleImports()
		if !assert.NoError(t, i.RegisterPrefix("viper", viperPkg)) {
			return
		}
		assert.EqualError(
			t,
			i.RegisterPrefix("viper", viperPkg),
			"shortcut `viper` is already registered",
		)
	})
}

func TestSimpleImports_GetImports(t *testing.T) {
	i := NewSimpleImports()
	i.GetAlias(viperPkg)
	assert.Equal(t, i.importsSlice, i.GetImports())
}
