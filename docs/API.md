# API

The following interfaces are not declared, given piece of code is prepared for docs purposes.

> Talk is cheap. Show me the code.
>
> -- <cite>Linus Torvalds</cite>

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

## Issues

All functions lock internal semaphore. To avoid deadlocks do not use container inside the container's callbacks.

### Deadlocks

In the following example `c.MustGet("foobar")` locks internal mutex.
`c.MustGet("foo")` tries to acquire the lock on the same, already locked mutex, which finally causes deadlock.

```go
c := container.NewAtomicContainer(container.NewContainer(nil))
// ...
c.Override("foobar", container.ServiceDefinition{
    Provider: func() (i interface{}, e error) {
        return c.MustGet("foo").(string) + "bar", nil
    },
})
c.MustGet("foobar")
```

To solve the above issue move `c.MustGet("foo")` out of the scope of the callback.

```go
c := container.NewAtomicContainer(container.NewContainer(nil))
// ...
foo := c.MustGet("foo").(string)
c.Override("foobar", container.ServiceDefinition{
    Provider: func() (i interface{}, e error) {
        return foo + "bar", nil
    },
})
c.MustGet("foobar")
```
