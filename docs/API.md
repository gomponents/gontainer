# API

The following interfaces are not declared, given piece of code is prepared for docs purposes.

> Talk is cheap. Show me the code.
-- <cite>Linus Torvalds</cite>

```go
package main

import (
    "github.com/gomponents/gontainer-helpers/container"
)

type ServiceContainer interface {
    // Get returns service by id.
    Get(id string) (interface{}, error)

    // MustGet returns service by id.
    MustGet(id string) interface{}

    // Has returns true when service is registered.
    Has(id string) bool

    // Register registers new service. It returns error when service is already registered.
    //
    // container.Register("logger", ServiceDefinition{
    //     Provider: func() (interface{}, error) {
    //         return logrus.New(), nil
    //     },
    //     Disposable: false,
    // })
    Register(id string, s container.ServiceDefinition) error

    // Override registers new service. When service is already registered will be replaced by new one.
    Override(id string, s container.ServiceDefinition)

    // GetAllServiceIDs returns IDs of all registered services.
    GetAllServiceIDs() []string

    // RegisterDecorator registers new decorator.
    RegisterDecorator(d container.Decorator)

    // ValidateAllServices tries to fetch all registered services and returns map of all errors which have been occurred.
    // Key of given map is an ID of a service.
    ValidateAllServices() (errors map[string]error)
}

type ParamContainer interface {
    // GetParam returns parameter by id.
    GetParam(id string) (interface{}, error)

    // MustGetParam returns parameter by id.
    MustGetParam(id string) interface{}

    // todo
}

type TaggedContainer interface {
    // todo
}

type Container interface {
    ServiceContainer
    ParamContainer
    TaggedContainer
}
```
