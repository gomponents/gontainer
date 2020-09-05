# TEST REPOSITORY

[![Build Status](https://github.com/gomponents/gontainer/workflows/Tests/badge.svg?branch=master)](https://github.com/gomponents/gontainer/actions?query=workflow%3AGo)
[![Coverage Status](https://coveralls.io/repos/github/gomponents/gontainer/badge.svg?branch=master)](https://coveralls.io/github/gomponents/gontainer?branch=master)

# Gontainer

Depenendency Injection container for GO inspired by [Symfony](https://symfony.com/doc/current/components/dependency_injection.html).

## TL;DR

**Describe dependencies in YAML**

```yaml
meta:
    imports:
        "pkg": "github.com/gontainer/repo/pkg"

parameters:              # No need to hardcode configuration values, e.g.:
    db.host: "localhost" # '%env("APP_DB_HOST")%'
    db.port: 3306        # '%envInt("APP_DB_PORT")%'

services:
    db:
        constructor: "pkg.NewDB" # equivalent for "github.com/gontainer/repo/pkg.NewDB", additionally
                                 # import can be surrounded by `"` to make it more explicit
                                 # e.g.:
                                 # - "pkg".NewDB
                                 # - "github.com/gontainer/repo/pkg".NewDB
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

## Brief

**Gontainer** builds DI container based on input YAML files.
Code is generated automatically, but internally it uses [reflect](https://golang.org/pkg/reflect/).
Whenever docs shows source code, given code is just an equivalent what is really going on inside to make docs easier to understand.

**Example**

```yaml
services:
    # db := db.NewDB(container.GetParameter("db.host"), ...
    db:
        constructor: "pkg/db.NewDB"
        args: ["%host%", "%port%"]
```

In the above example, generated code will differ than `db := db.NewDB(container.GetParameter("db.host"), ...`,
because internally it uses reflection (GO is statically typed and conversion of parameter is required),
however result will work as described using GO code.

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

    container_type: "Gontainer"          # Type of declared container, default "Gontainer".

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
* `%envInt("PORT")%` - returns converted to int value of environment variable `PORT`.
* `%todo()%` - fake parameter, can be used during development to avoid compiler errors (e.g. `service "db" requires param "db.host", but it does not exist`).

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
* reference to value, e.g.: `!value &pkg.MyStruct{}`

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

Constructor must return 1 or 2 values. Second (optional) value must be an instance of error.

**Sample constructors**

```go
type Server struct {
    Port int
}

// NewServer is just a constructor.
func NewServer(port int) *Server {
    return &Server{Port: port}
}

// NewServerWithError is a constructor, but returns an error whenever port is equal to 0.
func NewServerWithError(port int) (*Server, error) {
    if port == 0 {
        return nil, fmt.Errorf("port cannot be equal to 0")
    }
    return NewServer(port), nil
}
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
    # myStorage := storage.Storage{}
    # myStorage.Db = container.Get("db")
    # myStorage.Debug = true
    myStorage:
        value: "pkg/storage.Storage{}"
        fields:
            Db: "@db"
            Debug: true
```

### Inject global variable

```yaml
services:
    # db := pkg.NewDB(config.GlobalConfig.DB)
    db:
        constructor: "pkg.NewDB"
        args: ['!value "config".GlobalConfig.DB'] # compiler doesn't know whether `config` or `config.GlobalConfig`
                                                  # is expected import name, to avoid issues, surround import by `"` characters

```

### More about `!value`

GO doesn't really force you to use constructors, this is the reason why **Gontainer** gives you choice,
you can create service by constructor either by value.

```yaml
services:
    db1:
        contructor: "pkg.NewDB"
        args: ["localhost", 3306]
    db2:
        value: "&pkg.DB"
        fields:
            host: "localhost"
            port: 3306
```

Sometimes it makes sense to use value directly in argument or field, to do it, prefix your value by `!value `, e.g.:

```yaml
services:
    httpClient:
        constructor: "pkg.NewHttpClient"
        args: ['!value "config".GlobalConfig.HttpClient']
```

Values allow to use:

* constant or variable: `MyConfig`, `"pkg.config".MyConfig`
* field of global variable: `"pkg".MyConfig.Some.Field`, `".".MyConfig.Some.Field`
(`"."` is an alias to current package name, it is required in this specific case, otherwise compiler considers `MyConfig` as an import)
* struct: `MyStruct{}`, `my/import.MyStruct{}`
* pointers: `&"pkg.config".MyConfig`, `&"pkg".MyConfig.Some.Field`, `&MyStruct{}`, `&my/import.MyStruct{}`

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

# svc := pkg.MakeTracedHttpClient(svc, serviceName, container.Get("tracer"))
```
