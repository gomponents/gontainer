package compiled

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateParamsReqParamsExist(t *testing.T) {
	scenarios := []struct {
		dto   DTO
		error string
	}{
		{
			dto: DTO{
				Params: []Param{
					{Name: "name", DependsOn: []string{"firstname"}},
				},
			},
			error: "parameter `name` requires param `firstname`, but it does not exist",
		},
		{
			dto: DTO{
				Params: []Param{
					{Name: "name", DependsOn: []string{"firstname"}},
					{Name: "firstname"},
				},
			},
		},
		{
			dto: DTO{
				Services: []Service{
					{
						Name: "db",
						Args: []Arg{
							{DependsOnParams: []string{"host"}},
						},
					},
				},
			},
			error: "service `db` requires param `host`, but it does not exist",
		},
		{
			dto: DTO{
				Services: []Service{{
					Name: "db",
					Args: []Arg{
						{DependsOnParams: []string{"host"}},
					},
				}},
				Params: []Param{{
					Name: "host",
				}},
			},
		},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
			err := ValidateParamsReqParamsExist(s.dto)
			if s.error == "" {
				assert.NoError(t, err)
				return
			}
			assert.EqualError(t, err, s.error)
		})
	}
}
