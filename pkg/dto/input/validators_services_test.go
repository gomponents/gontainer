package input

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateServiceName(t *testing.T) {
	scenarios := []struct {
		name  string
		error string
	}{
		{
			name: "my.service",
		},
		{
			name:  "my..service",
			error: "service name must match pattern `" + regexServiceName.String() + "`, `my..service` given",
		},
		{
			name:  "%service%",
			error: "service name must match pattern `" + regexServiceName.String() + "`, `%service%` given",
		},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
			err := ValidateServiceName(s.name)
			if s.error == "" {
				assert.NoError(t, err)
				return
			}

			assert.EqualError(t, err, s.error)
		})
	}
}

func TestValidateConstructorType(t *testing.T) {
	scenarios := []struct {
		service Service
		error   string
	}{
		{
			service: Service{
				Constructor: "NewService",
			},
		},
		{
			error: "missing constructor, value or type",
		},
		{
			service: Service{
				Value:       "MyValue",
				Constructor: "MyConstructor",
			},
			error: "cannot define constructor and value together",
		},
		{
			service: Service{
				Constructor: "MyConstructor",
				Type:        "MyType",
			},
			error: "defined type will not be used, provide getter",
		},
		{
			service: Service{
				Value: "MyValue",
				Type:  "MyType",
			},
			error: "defined type will not be used, provide getter",
		},
		{
			service: Service{
				Type: "MyType",
				Args: []interface{}{"param"},
			},
			error: "arguments are not empty, but constructor is missing",
		},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
			err := ValidateConstructorType(s.service)
			if s.error == "" {
				assert.NoError(t, err)
				return
			}

			assert.EqualError(t, err, s.error)
		})
	}
}

func TestValidateServiceGetter(t *testing.T) {
	scenarios := []struct {
		getter string
		type_  string
		error  string
	}{
		{
			getter: "",
		},
		{
			getter: "GetName",
			error:  "getter is given, but type is missing",
		},
		{
			getter: "GetName",
			type_:  "MyType",
		},
		{
			getter: "getName",
			type_:  "MyType",
		},
		{
			getter: "0getName",
			error:  "getter must match `" + regexServiceGetter.String() + "`, `0getName` given",
			type_:  "MyType",
		},
		{
			getter: "Get Name",
			error:  "getter must match `" + regexServiceGetter.String() + "`, `Get Name` given",
			type_:  "MyType",
		},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
			err := ValidateServiceGetter(Service{Getter: s.getter, Type: s.type_})
			if s.error == "" {
				assert.NoError(t, err)
				return
			}

			assert.EqualError(t, err, s.error)
		})
	}
}

func TestValidateServiceType(t *testing.T) {
	scenarios := []struct {
		type_ string
		error string
	}{
		{
			type_: "",
		},
		{
			type_: "my/import/foo.Bar",
		},
		{
			type_: "*my/import/foo.Bar",
		},
		{
			type_: "**my/import/foo.Bar",
			error: "type must match `" + regexServiceType.String() + "`, `**my/import/foo.Bar` given",
		},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
			err := ValidateServiceType(Service{Type: s.type_})
			if s.error == "" {
				assert.NoError(t, err)
				return
			}

			assert.EqualError(t, err, s.error)
		})
	}
}

func TestValidateServiceValue(t *testing.T) {
	scenarios := []struct {
		val   string
		error string
	}{
		{
			val: "",
		},
		{
			val: "my/import/foo.Bar",
		},
		{
			val: "my/import/foo.MyStruct{}.Bar",
		},
		{
			val:   "my/import/foo.MyStruct{.Bar",
			error: "value must match `" + regexServiceValue.String() + "`, `my/import/foo.MyStruct{.Bar` given",
		},
		{
			val:   "my/import/foo",
			error: "value must match `" + regexServiceValue.String() + "`, `my/import/foo` given",
		},
		{
			val:   "*my/import/foo.Bar",
			error: "value must match `" + regexServiceValue.String() + "`, `*my/import/foo.Bar` given",
		},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
			err := ValidateServiceValue(Service{Value: s.val})
			if s.error == "" {
				assert.NoError(t, err)
				return
			}

			assert.EqualError(t, err, s.error)
		})
	}
}

func TestValidateServiceConstructor(t *testing.T) {
	scenarios := []struct {
		constructor string
		error       string
	}{
		{
			constructor: "",
		},
		{
			constructor: "my/import/foo.NewBar",
		},
		{
			constructor: "NewBar",
		},
		{
			constructor: "gopkg.in/yaml.v2.NewSth",
		},
		{
			constructor: "my/import/foo..NewBar",
		},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
			err := ValidateServiceConstructor(Service{Constructor: s.constructor})
			if s.error == "" {
				assert.NoError(t, err)
				return
			}

			assert.Equal(t, err, s.error)
		})
	}
}

func TestValidateServiceArgs(t *testing.T) {
	scenarios := []struct {
		args  []interface{}
		error string
	}{
		{
			args: nil,
		},
		{
			args: []interface{}{
				[]interface{}{},
			},
			error: "unsupported type `[]interface {}` of arg0",
		},
		{
			args: []interface{}{
				"hello",
				1,
				false,
				math.Pi,
			},
		},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
			err := ValidateServiceArgs(Service{Args: s.args})
			if s.error == "" {
				assert.NoError(t, err)
				return
			}

			assert.EqualError(t, err, s.error)
		})
	}
}

func TestValidateServiceCalls(t *testing.T) {
	scenarios := []struct {
		calls []Call
		error string
	}{
		{
			calls: []Call{
				{
					Method:    "WithDB",
					Args:      []interface{}{"@db"},
					Immutable: true,
				},
			},
		},
		{
			calls: []Call{
				{
					Method: "Incorrect method name",
				},
			},
			error: "method name (call 0) must match `" + regexServiceCallName.String() + "`, `Incorrect method name` given",
		},
		{
			calls: []Call{
				{
					Method: "SetID",
					Args:   []interface{}{struct{}{}},
				},
			},
			error: "unsupported type `struct {}` of (call `SetID`, arg 0)",
		},
	}
	for i, s := range scenarios {
		t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
			err := ValidateServiceCalls(Service{Calls: s.calls})
			if s.error == "" {
				assert.NoError(t, err)
				return
			}
			assert.EqualError(t, err, s.error)
		})
	}
}

func TestValidateServiceFields(t *testing.T) {
	scenarios := []struct {
		fields map[string]interface{}
		error  string
	}{
		{
			fields: map[string]interface{}{
				"foo": "bar",
			},
		},
		{
			fields: map[string]interface{}{
				"Incorrect field name": "bar",
			},
			error: "field must match `" + regexServiceFieldName.String() + "`, `Incorrect field name` given",
		},
		{
			fields: map[string]interface{}{
				"foo": struct{}{},
			},
			error: "unsupported type `struct {}` of field `foo`",
		},
		{
			fields: map[string]interface{}{
				"Bar": []uint{},
			},
			error: "unsupported type `[]uint` of field `Bar`",
		},
	}
	for i, s := range scenarios {
		t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
			err := ValidateServiceFields(Service{Fields: s.fields})
			if s.error == "" {
				assert.NoError(t, err)
				return
			}
			assert.EqualError(t, err, s.error)
		})
	}
}

func TestValidateServiceTags(t *testing.T) {
	scenarios := []struct {
		tags  []Tag
		error string
	}{
		{
			tags: []Tag{
				{
					Name:     "logger",
					Priority: 0,
				},
				{
					Name:     "logger",
					Priority: 10,
				},
			},
			error: "duplicate tag `logger`",
		},
		{
			tags: []Tag{
				{
					Name:     "logger",
					Priority: 0,
				},
			},
		},
		{
			tags: []Tag{
				{
					Name:     "white space",
					Priority: 0,
				},
			},
			error: "tag must match `" + regexServiceTag.String() + "`, `white space` given",
		},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
			err := ValidateServiceTags(Service{Tags: s.tags})
			if s.error == "" {
				assert.NoError(t, err)
				return
			}
			assert.EqualError(t, err, s.error)
		})
	}
}

func TestValidateServices(t *testing.T) {
	t.Run("Given incorrect service name", func(t *testing.T) {
		d := DTO{
			Services: map[string]Service{
				"white space": {},
			},
		}
		err := ValidateServices(d)
		assert.EqualError(t, err, "service name must match pattern `"+regexServiceName.String()+"`, `white space` given")
	})
	t.Run("Given incorrect type", func(t *testing.T) {
		d := DTO{
			Services: map[string]Service{
				"db": {Type: "@!#$"},
			},
		}
		assert.EqualError(t, ValidateServices(d), "service `db`: type must match `"+regexServiceType.String()+"`, `@!#$` given")
	})
	t.Run("Given todo", func(t *testing.T) {
		d := DTO{
			Services: map[string]Service{
				"db": {Todo: true, Type: "@!#$"}, // incorrect type, but service is marked as todo
			},
		}
		assert.NoError(t, ValidateServices(d))
	})
}
