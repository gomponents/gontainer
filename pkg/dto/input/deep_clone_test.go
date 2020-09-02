package input

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
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

	scenarios := []struct {
		input interface{}
	}{
		{
			input: []int{1, 2, 3},
		},
		{
			input: []interface{}{"hello", 1, math.Pi},
		},
		{
			input: map[interface{}]interface{}{
				"name": "Mary",
				3:      math.Pi,
			},
		},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
			output := deepClone(s.input)
			assert.Equal(t, s.input, output)
		})
	}
}
