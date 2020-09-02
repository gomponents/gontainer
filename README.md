# TEST REPOSITORY

[![Build Status](https://github.com/gomponents/gontainer/workflows/Tests/badge.svg?branch=master)](https://github.com/gomponents/gontainer/actions?query=workflow%3AGo)
[![Coverage Status](https://coveralls.io/repos/github/gomponents/gontainer/badge.svg?branch=master)](https://coveralls.io/github/gomponents/gontainer?branch=master)

# Gontainer

Depenendency Injection container for GO inspired by [Symfony](https://symfony.com/doc/current/components/dependency_injection.html).

## TL;DR

**Describe dependencies in YAML**

```yaml
parameters:              # No need to hardcode configuration values, e.g.:
    db.host: "localhost" # '%env("APP_DB_HOST")%'
    db.port: 3306        # '%envInt("APP_DB_PORT")%'

services:
    db:
        constructor: "pkg.NewDB"
        args: ["%db.host%", "%db.port%"]
    storage:
        constructor: "pkg.NewStorage"
        args: ["@db"]
        getter: "GetStorage"
        type: "*pkg.Storage"
```

**Voilà!**

```go
c := NewContainer()
s, err := c.GetStorage()
```

## Command

Flag `-i` supports [glob](https://golang.org/pkg/path/filepath/#Glob) patterns.

```bash
gontainer build -i container.yml -i container_dev.yml [...] -o container.go
```

Files are being processed from the left to the right, it means in the above example
at first `container.yml` will be parsed, then values from `container_dev.yml`
will override already loaded values.

## Schema

```yaml
meta:
# additional options

parameters:
# list of parameters

services:
# list of services
```

## Meta

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

## Parameters

Content between percent signs is a `%reference%` to another parameter or a `%function()%` (`%sth% != %sth()%`).

```yaml
parameters:
    env: '%env("ENVIRONMENT")%' # os.Get("ENVIRONMENT")
    host: "localhost"           # "localhost"
    port: 80                    # 80
    hostport: "%host%:%port%"   # "localhost:80" // ToString(container.GetParam("host")) + ":" + ToString(container.GetParam("port"))
```

Gontainer has 3 default functions:

* `%env("NAME")%` - returns value of environment variable `NAME`.
* `%envInt("PORT")%` - returns value of environment variable `PORT` and converts to int.
* `%todo()%` - fake parameter, can be used during development to avoid compiler erro%sth% != %sth()%rs.

You can override all parameters in runtime (`container.OverrideParam`), it can be useful when combined with `%todo()%`.

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

## Services

Fields, arguments of constructors and calls accept the same syntax as parameters and in addition:

* reference to any other service, e.g.: `@service`
* reference to group of tagged services, e.g.: `!tagged my.tag`

### Create service using constructor

```yaml
parameters:
    db.host: "localhost"
    db.port: 3306

services:
    # db := db.NewDB(container.GetParameter("db.host"), ...
    db:
        constructor: "pkg/db.NewDB"
        args: ["%host%", "%port%"]
```

### Setter injection

```yaml
services:
    # db := db.NewDB(container.GetParameter("db.host"), ...
    # db.Debug(true)
    db:
        constructor: "pkg/db.NewDB"
        args: ["%host%", "%port%"]
        calls:
            - ["Debug", [true]] # see https://symfony.com/doc/current/service_container/calls.html
                                # see https://symfony.com/blog/new-in-symfony-4-3-configuring-services-with-immutable-setters
```

### Direct injection

```yaml
services:
    # var myStorage storage.Storage
    # myStorage.Db = container.Get("db")
    # myStorage.Debug = true
    myStorage:
        type: "pkg/storage.Storage"
        fields:
            Db: "@db"
            Debug: true
```

### Tags

```yaml
services:
    handlerOne:
        constructor: "pkg.NewHandler1"
        tags: ["handler"]

    handlerTwo:
        constructor: "pkg.NewHandler2"
        tags: [{"name": "handler", "priority": 100}]

    # handlerCollection := pkg.NewHandlerCollection(container.MustGetByTag("handler"))
    handlerCollection:
        constructor: "pkg.NewHandlerCollection"
        args: ["!tagged handler"]
```

## TODO

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

**Create service by type**

remove:

```yaml
service:
    foo:
        type: "MyType"
```

and replace by:

```yaml
service:
    foo:
        value: "MyType{}"
```

**!value**

```yaml
args: ["!value config.Cfg.Some.Global.Config", "!value &myStruct{}"]
```