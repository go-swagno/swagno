# Swagno Project Documentation

This folder contains detailed documentation for the Swagno project.

## Documentation Files

1. **[01-project-overview.md](01-project-overview.md)** - Project overview and architecture
2. **[02-core-modules.md](02-core-modules.md)** - Core modules and basic structures
3. **[03-components-guide.md](03-components-guide.md)** - Detailed explanation of all components
4. **[04-examples-and-usage.md](04-examples-and-usage.md)** - Usage examples and framework integrations
5. **[05-api-reference.md](05-api-reference.md)** - API reference and function descriptions

## Quick Start

Swagno is a library used to create Swagger 2.0 documentation in Go projects. It requires no annotations, no file exports, and no command execution.

```go
// Basic usage
sw := swagno.New(swagno.Config{Title: "My API", Version: "v1.0.0"})
sw.AddEndpoints(endpoints)
jsonDoc := sw.MustToJson()
```

Please refer to the documentation files above for detailed information.
