# Swagno v3 Documentation Summary

This document provides a quick reference and summary of all documentation available for Swagno v3 (OpenAPI 3.0.3 support).

## 📁 Documentation Structure

The documentation is organized into the following files:

1. **[README.md](README.md)** - Overview and quick start guide
2. **[01-project-overview.md](01-project-overview.md)** - Architecture and OpenAPI 3.0.3 features
3. **[02-core-modules.md](02-core-modules.md)** - Core modules and data structures
4. **[03-components-guide.md](03-components-guide.md)** - Detailed component explanations
5. **[04-examples-and-usage.md](04-examples-and-usage.md)** - Practical examples and framework integration
6. **[05-api-reference.md](05-api-reference.md)** - Complete API reference
7. **[06-migration-guide.md](06-migration-guide.md)** - Migration from Swagger 2.0 to OpenAPI 3.0

## 🚀 Quick Start

```go
import v3 "github.com/go-swagno/swagno/v3"

// Create OpenAPI instance
openapi := v3.New(v3.Config{
    Title:   "My API",
    Version: "v1.0.0",
})

// Add server
openapi.AddServer("https://api.example.com/v1", "Production server")

// Add security
openapi.SetBearerAuth("JWT", "Bearer authentication")

// Add endpoints
openapi.AddEndpoints(endpoints)

// Generate documentation
jsonDoc := openapi.MustToJson()
```

## 🔧 Key Features

### OpenAPI 3.0.3 Compliance

- Full specification compliance
- Multiple servers support
- Enhanced security schemes
- Request bodies with content types
- Response links and callbacks

### Enhanced Security

- Basic Authentication
- Bearer Authentication (JWT)
- API Key Authentication
- OAuth2 with multiple flows
- OpenID Connect

### Advanced Components

- Enhanced schema definitions
- Nullable type support
- ReadOnly/WriteOnly properties
- Parameter styles and examples
- Callbacks and links

## 📊 Architecture Overview

```
v3/
├── openapi.go              # Main OpenAPI struct
├── generate.go             # JSON generation
├── auth.go                 # Security schemes
├── external_docs.go        # External documentation
├── components/             # OpenAPI 3.0 components
│   ├── definition/         # Schema definitions
│   ├── endpoint/           # API endpoints
│   ├── parameter/          # Parameters
│   ├── http/response/      # Responses
│   ├── security/           # Security
│   └── tag/               # Tags
└── example/               # Usage examples
```

## 🔄 Migration from Swagger 2.0

### Key Changes

| Feature            | Swagger 2.0                   | OpenAPI 3.0                      |
| ------------------ | ----------------------------- | -------------------------------- |
| **Import**         | `github.com/go-swagno/swagno` | `github.com/go-swagno/swagno/v3` |
| **Constructor**    | `swagno.New()`                | `v3.New()`                       |
| **Servers**        | `host` + `basePath`           | `servers` array                  |
| **Schemas**        | `#/definitions/`              | `#/components/schemas/`          |
| **Request Bodies** | Part of parameters            | Separate `requestBody`           |

### Migration Steps

1. Update import path
2. Change constructor
3. Replace host/basePath with servers
4. Update security schemes
5. Enhance with OpenAPI 3.0 features

## 📚 Documentation Deep Dive

### For Beginners

Start with **[README.md](README.md)** and **[01-project-overview.md](01-project-overview.md)** to understand the basics and architecture.

### For Developers

Read **[02-core-modules.md](02-core-modules.md)** and **[03-components-guide.md](03-components-guide.md)** for detailed technical information.

### For Implementation

Check **[04-examples-and-usage.md](04-examples-and-usage.md)** for practical examples and framework integration patterns.

### For Migration

Follow **[06-migration-guide.md](06-migration-guide.md)** for step-by-step migration instructions from Swagger 2.0.

### For Reference

Use **[05-api-reference.md](05-api-reference.md)** as a complete function and type reference.

## 🛠 Code Examples

### Basic API

```go
openapi := v3.New(v3.Config{
    Title:   "User API",
    Version: "v1.0.0",
})

openapi.AddServer("https://api.example.com/v1", "Production")
openapi.SetBearerAuth("JWT", "Bearer authentication")

endpoint := endpoint.New(
    endpoint.GET,
    "/users/{id}",
    endpoint.WithTags("users"),
    endpoint.WithParams(
        parameter.IntParam("id", parameter.Path,
            parameter.WithRequired(),
            parameter.WithExample(123),
        ),
    ),
    endpoint.WithSuccessfulReturns([]response.Response{
        response.New(User{}, "200", "User found"),
    }),
)

openapi.AddEndpoint(endpoint)
```

### Advanced Features

```go
// Multiple security schemes
openapi.SetBearerAuth("JWT", "JWT authentication")
openapi.SetApiKeyAuth("X-API-Key", "header", "API key auth")

// OAuth2 with multiple flows
flows := &v3.OAuthFlows{
    AuthorizationCode: &v3.OAuthFlow{
        AuthorizationUrl: "https://example.com/oauth/authorize",
        TokenUrl:        "https://example.com/oauth/token",
        Scopes: map[string]string{
            "read":  "Read access",
            "write": "Write access",
        },
    },
}
openapi.SetOAuth2Auth(flows, "OAuth2 authentication")

// Enhanced parameters
parameter.StringParam("filter", parameter.Query,
    parameter.WithStyle("deepObject"),
    parameter.WithExplode(true),
    parameter.WithExamples(map[string]parameter.Example{
        "active": {Summary: "Active users", Value: "active=true"},
        "inactive": {Summary: "Inactive users", Value: "active=false"},
    }),
)

// Response with links
response.NewWithLinks(User{}, "201", "User created",
    map[string]endpoint.Link{
        "GetUser": endpoint.NewEnhancedLink().
            SetOperationId("getUserById").
            AddParameter("id", "$response.body#/id"),
    },
)
```

## 🔍 Key Differences from Original Swagno

### Enhanced Features

- **OpenAPI 3.0.3 Specification**: Full compliance with modern OpenAPI standard
- **Multiple Servers**: Support for multiple environment configurations
- **Enhanced Security**: Bearer tokens, OpenID Connect, enhanced OAuth2
- **Better Content Types**: Proper request body and response content handling
- **Callbacks & Links**: Support for webhooks and response relationships
- **Parameter Styles**: Advanced parameter serialization options

### Backward Compatibility

- **API Similarity**: Most functions work the same way for easy migration
- **Enhanced Options**: Existing patterns extended with new OpenAPI 3.0 features
- **Progressive Enhancement**: Can start simple and add advanced features gradually

## 🎯 Best Practices

### Documentation Quality

- Use descriptive summaries and descriptions
- Add examples to parameters and responses
- Include proper error responses
- Tag endpoints for organization

### Security

- Use appropriate security schemes for your API
- Apply security at operation level when needed
- Document security requirements clearly

### Schema Design

- Use struct tags for validation and examples
- Leverage nullable types (\*Type) for optional fields
- Add descriptions to model fields
- Use readonly/writeonly appropriately

### API Design

- Follow RESTful conventions
- Use appropriate HTTP status codes
- Include proper content types
- Add server configurations for different environments

## 📞 Support and Contribution

This documentation covers the complete Swagno v3 functionality. For additional help:

- Check the example files in `/example/` directory
- Review test files for usage patterns
- Refer to OpenAPI 3.0.3 specification for advanced features

The documentation is designed to be comprehensive yet accessible, covering everything from basic usage to advanced OpenAPI 3.0 features.
