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

func TestValidateServicesCircularDeps(t *testing.T) {
	scenarios := []struct {
		services []Service
		error    string
	}{
		{
			services: []Service{
				{
					Name: "harryPotter",
					Tags: []Tag{{Name: "book"}},
					Calls: []Call{
						{Method: "SetLibrary", Args: []Arg{{DependsOnServices: []string{"library"}}}},
					},
				},
				{
					Name: "library",
					Fields: []Field{
						{Name: "Books", Value: Arg{DependsOnTags: []string{"book"}}},
					},
				},
			},
			error: "circular dependency in services: harryPotter -> library -> harryPotter",
		},
		{
			services: []Service{
				{
					Name: "harryPotter",
					Tags: []Tag{{Name: "book"}},
				},
				{
					Name: "library",
					Fields: []Field{
						{Name: "Books", Value: Arg{DependsOnTags: []string{"book"}}},
					},
				},
			},
		},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
			err := ValidateServicesCircularDeps(DTO{Services: s.services})
			if s.error == "" {
				assert.NoError(t, err)
				return
			}
			assert.EqualError(t, err, s.error)
		})
	}
}
