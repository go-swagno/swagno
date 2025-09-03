# Swagno v3 Documentation

This folder contains detailed documentation for Swagno v3 (OpenAPI 3.0.3 support).

## Documentation Files

1. **[01-project-overview.md](01-project-overview.md)** - Project overview and OpenAPI 3.0.3 architecture
2. **[02-core-modules.md](02-core-modules.md)** - Core modules and basic OpenAPI 3.0 structures
3. **[03-components-guide.md](03-components-guide.md)** - Detailed explanation of all OpenAPI 3.0 components
4. **[04-examples-and-usage.md](04-examples-and-usage.md)** - Usage examples and migration from Swagger 2.0
5. **[05-api-reference.md](05-api-reference.md)** - API reference and function descriptions
6. **[06-migration-guide.md](06-migration-guide.md)** - Migration guide from Swagger 2.0 to OpenAPI 3.0

## Quick Start

Swagno v3 is a library used to create OpenAPI 3.0.3 documentation in Go projects. It requires no annotations, no file exports, and no command execution.

```go
// Basic usage
openapi := v3.New(v3.Config{Title: "My API", Version: "v1.0.0"})
openapi.AddEndpoints(endpoints)
jsonDoc := openapi.MustToJson()
```

## Key Features

- ✅ **OpenAPI 3.0.3 Specification**: Full compliance with OpenAPI 3.0.3
- ✅ **No annotations required**: No need to write special comments in code
- ✅ **No file exports**: No need to create separate OpenAPI files
- ✅ **No command execution**: No additional steps required in build process
- ✅ **Enhanced Security**: Support for Bearer tokens, OpenID Connect, and more
- ✅ **Better Schema Support**: Improved schema definitions with nullable, readOnly, writeOnly
- ✅ **Multiple Servers**: Support for multiple server configurations
- ✅ **Request Bodies**: Proper request body handling with content types

## Migration from Swagger 2.0

If you're migrating from the original Swagno (Swagger 2.0), please refer to **[06-migration-guide.md](06-migration-guide.md)** for detailed migration instructions.

## Framework Integration

Swagno v3 works with all major Go web frameworks:

- **net/http**: Native Go HTTP server
- **Gin**: Popular Go web framework
- **Fiber**: Express-inspired web framework
- **Echo**: High performance, extensible, minimalist Go web framework

Please refer to the documentation files above for detailed information and examples.
