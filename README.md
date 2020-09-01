# TEST REPOSITORY

[![Build Status](https://github.com/gomponents/gontainer/workflows/Tests/badge.svg?branch=master)](https://github.com/gomponents/gontainer/actions?query=workflow%3AGo)
[![Coverage Status](https://coveralls.io/repos/github/gomponents/gontainer/badge.svg?branch=master)](https://coveralls.io/github/gomponents/gontainer?branch=master)

# Gontainer

Depenendency Injection container for GO inspired by Symfony.

## Cmd

```bash
ENV=dev
gontainer build -i container.yml -i container_${DEV}.yml -o container.go
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
```
decorators:
    - tag: http-client
      decorator: myImport/pkg.MakeTracedHttpClient
      args: [@tracer]

# svc := pkg.MakeTracedHttpClient(svc, container.Get("tracer"))

    - instanceof: myImport/pkg.HttpClient
      decorator: myImport/pkg.MakeTracedHttpClient
      args: [@tracer]

# if _, ok := svc.(pkg.HttpClient); ok {
#     svc = pkg.MakeTracedHttpClient(svc, container.Get("tracer"))
# }
```
