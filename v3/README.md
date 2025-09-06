# Swagno v3 - OpenAPI 3.0 Support

This directory contains the OpenAPI 3.0 version of Swagno, providing support for generating OpenAPI 3.0.3 specifications.

## Overview

Swagno v3 is a library for creating OpenAPI 3.0 documentation in Go projects while maintaining the same philosophy as the original Swagno: **"no annotations, no files, no commands"** - meaning no annotations, no file exports, no command execution.

### Key Features

- ✅ **OpenAPI 3.0.3 Specification**: Full compliance with OpenAPI 3.0.3
- ✅ **No annotations required**: No need to write special comments in code
- ✅ **No file exports**: No need to create separate OpenAPI files
- ✅ **No command execution**: No additional steps required in build process
- ✅ **Programmatic approach**: All OpenAPI definitions are defined with Go code
- ✅ **Framework agnostic**: Works with Gin, Fiber, net/http and other frameworks
- ✅ **Enhanced Security**: Support for Bearer tokens, OpenID Connect, and more
- ✅ **Better Schema Support**: Improved schema definitions with nullable, readOnly, writeOnly
- ✅ **Multiple Servers**: Support for multiple server configurations
- ✅ **Request Bodies**: Proper request body handling with content types

## Quick Start

```go
package main

import (
    swagno3 "github.com/go-swagno/swagno/v3"
    "github.com/go-swagno/swagno/v3/components/endpoint"
    "github.com/go-swagno/swagno/v3/components/http/response"
    "github.com/go-swagno/swagno/v3/components/parameter"
)

func main() {
    // Create OpenAPI 3.0 instance
    openapi := swagno3.New(swagno3.Config{
        Title:       "My API",
        Version:     "v1.0.0",
        Description: "My awesome API using OpenAPI 3.0",
    })

    // Add servers
    openapi.AddServer("https://api.example.com/v1", "Production server")

    // Add security
    openapi.SetBearerAuth("JWT", "Bearer authentication")

    // Define endpoints
    endpoints := []*endpoint.EndPoint{
        endpoint.New(
            endpoint.GET,
            "/users",
            endpoint.WithSuccessfulReturns([]response.Response{
                response.New([]User{}, "200", "List of users"),
            }),
        ),
    }

    openapi.AddEndpoints(endpoints)
    jsonDoc := openapi.MustToJson()
}
```

## Migration from Swagger 2.0

### Key Differences

1. **Import Path**: Use `github.com/go-swagno/swagno/v3` instead of `github.com/go-swagno/swagno`
2. **Constructor**: Use `swagno3.New()` instead of `swagno.New()`
3. **Servers**: Use `AddServer()` instead of setting Host and BasePath
4. **Security**: Enhanced security schemes with Bearer, OpenID Connect support
5. **Schemas**: Located at `#/components/schemas/` instead of `#/definitions/`

### Migration Example

**Swagger 2.0 (old):**

```go
sw := swagno.New(swagno.Config{
    Title: "My API",
    Version: "v1.0.0",
    Host: "api.example.com",
    Path: "/v1",
})
sw.SetBasicAuth("Basic auth required")
```

**OpenAPI 3.0 (new):**

```go
openapi := swagno3.New(swagno3.Config{
    Title: "My API",
    Version: "v1.0.0",
})
openapi.AddServer("https://api.example.com/v1", "Production server")
openapi.SetBasicAuth("Basic auth required")
```

## Major Enhancements in v3

### 1. Multiple Servers Support

```go
openapi.AddServer("https://api.example.com/v1", "Production server")
openapi.AddServer("https://staging-api.example.com/v1", "Staging server")
openapi.AddServer("http://localhost:8080/v1", "Local development")
```

### 2. Enhanced Security Schemes

```go
// Bearer Authentication (JWT)
openapi.SetBearerAuth("JWT", "Bearer authentication using JWT tokens")

// API Key Authentication
openapi.SetApiKeyAuth("X-API-Key", "header", "API key authentication")

// OAuth2 with multiple flows
flows := &security.OAuthFlows{
    AuthorizationCode: &security.OAuthFlow{
        AuthorizationUrl: "https://example.com/oauth/authorize",
        TokenUrl:        "https://example.com/oauth/token",
        Scopes: map[string]string{
            "read":  "Read access",
            "write": "Write access",
        },
    },
}
openapi.SetOAuth2Auth(flows, "OAuth2 authentication")

// OpenID Connect
openapi.SetOpenIdConnectAuth("https://example.com/.well-known/openid_configuration")
```

### 3. Request Body Handling

OpenAPI 3.0 handles request bodies differently with proper content type support:

```go
endpoint.New(
    endpoint.POST,
    "/users",
    endpoint.WithBody(CreateUserRequest{}),
    // Request body automatically gets proper content-type handling
)
```

### 4. Enhanced Schema Definitions

```go
type User struct {
    ID       uint64     `json:"id"`
    Name     string     `json:"name"`
    Email    *string    `json:"email,omitempty"`    // Nullable in OpenAPI 3.0
    IsActive bool       `json:"is_active"`
    Metadata map[string]interface{} `json:"metadata,omitempty"`
}
```

### 5. Response Content Types

```go
// Responses now have proper content type definitions
response.New(User{}, "200", "User found")
// Automatically generates:
// {
//   "200": {
//     "description": "User found",
//     "content": {
//       "application/json": {
//         "schema": {
//           "$ref": "#/components/schemas/User"
//         }
//       }
//     }
//   }
// }
```

### 6. Parameter Enhancements

```go
parameter.IntParam("id", parameter.Path,
    parameter.WithRequired(),
    parameter.WithDescription("User ID"),
    parameter.WithExample(123),           // OpenAPI 3.0 examples
    parameter.WithMin(1),
    parameter.WithMax(999999),
    parameter.WithDeprecated(),           // Mark as deprecated
)
```

## Components Structure

The v3 package maintains the same modular structure as the original:

```
v3/
├── openapi.go              # Main OpenAPI struct and functions
├── auth.go                 # Security definitions
├── generate.go             # JSON generation
├── components/
│   ├── definition/         # Schema definitions (OpenAPI 3.0 schemas)
│   ├── endpoint/           # API endpoint definitions
│   ├── fields/             # Struct field parsing utilities
│   ├── http/response/      # HTTP response structures
│   ├── mime/               # MIME type definitions
│   ├── parameter/          # API parameters with OpenAPI 3.0 features
│   ├── security/           # Security structures and scopes
│   └── tag/                # OpenAPI tags
└── example/
    └── basic/              # Basic usage example
```

## Backward Compatibility

While the v3 package is designed for OpenAPI 3.0, it maintains API compatibility where possible. Most of your existing Swagno code will work with minimal changes when migrating to v3.

## Examples

Check the `example/` directory for comprehensive usage examples:

- `basic/main.go` - Complete API example with CRUD operations
- Shows proper use of security, parameters, request bodies, and responses

## Framework Integration

The v3 package will work with the same framework integrations as the original Swagno:

- **net/http**: Native Go HTTP server
- **Gin**: Popular Go web framework
- **Fiber**: Express-inspired web framework

Framework-specific handlers will be updated to support OpenAPI 3.0 in separate repositories.

## Future Enhancements

- [ ] Enhanced validation support
- [ ] Webhook support
- [ ] Link support
- [ ] Callback support
- [ ] OpenAPI extensions
- [ ] Schema composition (allOf, oneOf, anyOf)

## Contributing

This is part of the Swagno project. The v3 package follows the same principles and coding standards as the main project while providing full OpenAPI 3.0.3 specification support.
