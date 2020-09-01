# TEST REPOSITORY

[![Build Status](https://github.com/gomponents/gontainer/workflows/Tests/badge.svg?branch=master)](https://github.com/gomponents/gontainer/actions?query=workflow%3AGo)
[![Coverage Status](https://coveralls.io/repos/github/gomponents/gontainer/badge.svg?branch=master)](https://coveralls.io/github/gomponents/gontainer?branch=master)

# Gontainer

Depenendency Injection container for GO inspired by [Symfony](https://symfony.com/doc/current/components/dependency_injection.html).

## Command

Flag `-i` supports [glob](https://golang.org/pkg/path/filepath/#Glob) patterns.

```bash
gontainer build -i container.yml -i container_dev.yml [...] -o container.go
```

## Schema

```yaml
meta:
# additional options

parameters:
# list of parameters

services:
# list of services
```

### Meta

```yaml
meta:
    pkg: "main"                          # Package name, default "main".

    container_type: "gontainer"          # Type of declared container, default "gontainer".

    imports:                             # List of aliases.
        "viper": github.com/spf13/viper" # It allows to use shorter syntax in service definition,
                                         # e.g.: "viper.New" instead of "github.com/spf13/viper.New".

    functions:                           # List of functions to use in parameters.
        "env": "os.Getenv"               # It allows to inject values calculated in runtime,
                                         # e.g.: 'env("ENVIRONMENT")'.
```

### Parameters

Content between percent signs is a `%reference%` to another parameter or a `%function()%`.

```yaml
parameters:
    env: '%env("ENVIRONMENT")%' # os.Get("ENVIRONMENT")
    host: "localhost"           # "localhost"
    port: 80                    # 80
    hostport: "%host%:%port%"   # "localhost:80"
```

Gontainer has 3 default functions:

* `%env("NAME")%` - returns value of environment variable `NAME`.
* `%envInt("PORT")%` - returns value of environment variable `PORT` and converts to int.
* `%todo()%` - fake parameter, can be used during development to avoid compiler errors.

You can override all parameters in runtime (`container.OverrideParam`), it can be useful when combined with `%todo%`.
All content between parentheses must be valid GO code, because it is directly used in compiled DI container.
The following code

```yaml
meta:
    functions:
        "sum": "pkg.Sum"
parameters:
    six: '%sum(1, 2, 3)%'
```

will be compiled to `pkg.Sum(1, 2, 3)`.

### Imports

TODO

### Services

```yaml
parameters:
    host: "localhost"
    port: 3306
    user: "root"
    password: "root"

services:
    # db := db.NewDB(container.GetParameter("host"), ...
    # db.Debug(true)
    db:
        constructor: "pkg/db.NewDB"
        args: ["%host%", "%port%", "%user%", "%password%"]
        calls:
            - ["Debug", [true]] # see https://symfony.com/doc/current/service_container/calls.html

    # var storage storage.Storage
    # storage.Db = container.Get("db")
    storage:
        type: "pkg/storage.Storage"
        fields:
            Db: "@db"

    handlerOne:
        constructor: "pkg.NewHandler1"
        tags ["handler"]

    handlerTwo:
        constructor: "pkg.NewHandler2"
        tags ["handler"]

    handlerCollection:
        constructor: "pkg.NewHandlerCollection"
        args: ["!tagged handler"]
```

## Example

```yaml
meta:
  pkg: container
  container_type: MyContainer
  imports:
    gontainer: "github.com/gomponents/gontainer"

parameters:
  first_name: '%env("NAME")%'
  last_name: "Doe"
  age: '%envInt("AGE")%'
  salary: 30000
  position: "CTO"

services:
  personExample1:
    type: "*gontainer/example/pkg/Employee" # alias.Employee{}

  personExample2:
    type: "gontainer/example/pkg/Employee" # &alias.Employee{}

  person: # alias.NewPerson(...)
    constructor: "github.com/gomponents/gontainer/example/pkg.NewPerson"
    args: ["%first_name% %last_name%", "%age%"]

  employee: # alias.NewEmployee(container.Get("person'), ...)
    getter: "Employee"
    type: "*gontainer/example/pkg/Employee"
    constructor: "gontainer/example/pkg.NewEmployee"
    args:
      - "@person.(*gontainer/example/pkg.Person)"
      - "%salary%"
      - "%position%"
```

```go
    c := CreateContainer()
    employee, employeeErr := c.GetEmployee()
    person, personErr := c.Get("person")
```

### TODO

**Decorators**
```yaml
decorators:
    - tag: http-client
      decorator: myImport/pkg.MakeTracedHttpClient
      args: [@tracer]

# svc := pkg.MakeTracedHttpClient(svc, container.Get("tracer"))

    - instanceof: myImport/pkg.HttpClient
      decorator: myImport/pkg.MakeTracedHttpClient
      args: [@tracer]

# if _, ok := svc.(pkg.HttpClient); ok {
#     svc = pkg.MakeTracedHttpClient(svc.(pkg.HttpClient), container.Get("tracer"))
# }
```
