# GoVel Framework

A modern Go web framework inspired by Laravel, focusing on developer experience, clean architecture, and elegant API design.

## Overview

GoVel brings the elegance and developer-friendly features of Laravel to the Go ecosystem. Built with clean architecture principles, interface segregation, and comprehensive lifecycle management.

## Features

- **Laravel-Inspired API**: Familiar directory structure and helper methods
- **Interface Segregation Principle (ISP)**: Clean, focused interfaces for better maintainability
- **Lifecycle Management**: Complete application lifecycle with hooks
- **Graceful Shutdown**: Proper cleanup and resource management
- **Modular Architecture**: Package-based organization for scalability

## Architecture

### Core Packages

- **application**: Core application package with lifecycle management
- **interfaces**: Clean, segregated interfaces following ISP

### Design Principles

1. **Interface Segregation**: Small, focused interfaces that do one thing well
2. **Dependency Injection**: Clean dependency management
3. **Lifecycle Hooks**: Extensible pre/post execution hooks
4. **Graceful Operations**: Proper startup and shutdown sequences

## Quick Start

```go
package main

import "govel/packages/application"

func main() {
    // Create new GoVel application
    application := application.New()
    
    // Configure application
    application.Configure()
    
    // Start the application
    application.Start()
}
```

## License

MIT License - see LICENSE file for details.
