# API

Compiled container implements the following interface.

```go
package main

type Container interface {
    // Get returns service by id
    Get(id string) (interface{}, error)

    // MustGet returns service by id
    MustGet(id string) interface{}
}
```
