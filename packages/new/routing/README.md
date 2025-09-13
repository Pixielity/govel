# Webserver Package

A unified webserver package for the Govel framework that provides Laravel-like HTTP handling with support for multiple web frameworks including GoFiber, Gin, and Echo.

## Features

- **Multiple Framework Support**: GoFiber, Gin, Echo adapters
- **Laravel-inspired API**: Familiar request/response handling
- **Unified Interface**: Same API regardless of underlying framework
- **Middleware System**: Framework-agnostic middleware support
- **Builder Pattern**: Fluent configuration API
- **Type Safety**: Full Go type safety and interfaces

## Quick Start

```go
package main

import (
    "govel/packages/new/webserver/src"
    "govel/packages/new/webserver/src/enums"
)

func main() {
    server := webserver.NewBuilder().
        WithEngine(enums.GoFiber).
        WithPort(8080).
        Build()

    server.Get("/", func(req *webserver.Request) *webserver.Response {
        return webserver.Json(map[string]string{
            "message": "Hello, Govel!",
        })
    })

    server.Listen()
}
```

## Supported Frameworks

- **GoFiber** - High performance HTTP framework
- **Gin** - Popular Go web framework
- **Echo** - High performance, extensible web framework

## Documentation

See the `__examples__` directory for usage examples and the `__tests__` directory for comprehensive test cases.
