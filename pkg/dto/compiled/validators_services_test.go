package compiled

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateServicesReqServicesExist(t *testing.T) {
	scenarios := []struct {
		services []Service
		error    string
	}{
		{
			services: []Service{
				{Name: "db"},
			},
		},
		{
			services: []Service{
				{Name: "storage", Fields: []Field{{Name: "DB", Value: Arg{DependsOnServices: []string{"db"}}}}},
			},
			error: "service `storage` requires service `db`, but it does not exist",
		},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
			err := ValidateServicesReqServicesExist(DTO{Services: s.services})
			if s.error == "" {
				assert.NoError(t, err)
				return
			}
			assert.EqualError(t, err, s.error)
		})
	}
}
