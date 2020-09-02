package input

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"testing"
)

func Test_deepClone(t *testing.T) {
	t.Run("Shallow copy", func(t *testing.T) {
		y := `{"name": "Mary"}`
		var v map[interface{}]interface{}
		assert.NoError(t, yaml.Unmarshal([]byte(y), &v))
		v2 := v
		v2["name"] = "Jane"

		assert.Equal(t, "Jane", v["name"])
		assert.Equal(t, "Jane", v2["name"])
	})
	t.Run("Deep copy", func(t *testing.T) {
		y := `{"name": "Mary"}`
		var v map[interface{}]interface{}
		assert.NoError(t, yaml.Unmarshal([]byte(y), &v))
		v2 := deepClone(v).(map[interface{}]interface{})
		v2["name"] = "Jane"

		assert.Equal(t, "Mary", v["name"])
		assert.Equal(t, "Jane", v2["name"])
	})
}
