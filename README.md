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

### Schema

```yaml
meta:
# additional options

parameters:
# list of parameters

services:
# list of services
```

#### Meta

```
meta:
    pkg: "main"                          # Package name, default "main".

    container_type: "gontainer"          # Type of declared container, default "gontainer".

    imports:                             # List of aliases.
        "viper": github.com/spf13/viper" # It allows to use shorter syntax in service definition,
                                         # e.g.: "viper.New" instead of "github.com/spf13/viper.New".

    functions:                           # List of functions to use in parameters.
        "env": "os.Getenv"               # It allows to inject values calculated in real-time,
                                         # e.g.: 'env("ENVIRONMENT")'
```

#### Imports



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
