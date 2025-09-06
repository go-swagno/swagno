# Project Overview and OpenAPI 3.0.3 Architecture

## Project Definition

**Swagno v3** is a library developed to create OpenAPI 3.0.3 documentation in Go projects. The main motto of the project remains: **"no annotations, no files, no commands"** - meaning no annotations, no file exports, no command execution.

### Key Features

- ✅ **OpenAPI 3.0.3 Compliance**: Full compliance with OpenAPI 3.0.3 specification
- ✅ **No annotations required**: No need to write special comments in code
- ✅ **No file exports**: No need to create separate OpenAPI files
- ✅ **No command execution**: No additional steps required in build process
- ✅ **Programmatic approach**: All OpenAPI definitions are defined with Go code
- ✅ **Framework agnostic**: Works with Gin, Fiber, net/http and other frameworks
- ✅ **Enhanced Security**: Bearer tokens, OpenID Connect, OAuth2 flows
- ✅ **Better Schema Support**: Enhanced schema definitions with OpenAPI 3.0 features
- ✅ **Multiple Servers**: Support for multiple server configurations
- ✅ **Request Bodies**: Proper request body handling with content types
- ✅ **Callbacks & Links**: Support for OpenAPI 3.0 callbacks and links

## Project Architecture

### Directory Structure

```
v3/
├── openapi.go          # Main OpenAPI struct and functions (replaces swagger.go)
├── generate.go         # JSON generation and OpenAPI definition creation
├── auth.go            # Security definitions (Bearer, OAuth2, OpenID Connect)
├── external_docs.go   # External documentation support
├── go.mod             # Go module definition for v3
├── components/        # OpenAPI 3.0 components
│   ├── definition/    # Schema definitions with OpenAPI 3.0 features
│   │   ├── definition.go      # Base schema definitions
│   │   └── enhanced_schema.go # Enhanced OpenAPI 3.0 schemas
│   ├── endpoint/      # API endpoint definitions with OpenAPI 3.0 features
│   │   ├── endpoints.go       # Core endpoint functionality
│   │   ├── callbacks_links.go # Callbacks and links support
│   │   └── operation_fixes.go # OpenAPI 3.0 operation enhancements
│   ├── fields/        # Struct field parsing tools
│   ├── http/          # HTTP structures with OpenAPI 3.0 enhancements
│   │   ├── enhanced_request_response.go # Enhanced request/response
│   │   └── response/  # HTTP response structures
│   ├── mime/          # MIME type definitions
│   ├── parameter/     # API parameters with OpenAPI 3.0 features
│   │   ├── parameter.go         # Base parameter functionality
│   │   └── openapi3_parameter.go # OpenAPI 3.0 specific features
│   ├── security/      # Security structures with OpenAPI 3.0 support
│   │   ├── security.go       # Security scheme types and constants
│   │   └── oauth_flows.go    # OAuth2 flows implementation
│   └── tag/           # OpenAPI tags
└── example/           # Usage examples for OpenAPI 3.0
    ├── basic/         # Basic OpenAPI 3.0 example
    ├── enhanced/      # Advanced OpenAPI 3.0 features
    └── models/        # Example model definitions
```

### Main Architecture Components

#### 1. Main OpenAPI Layer (`openapi.go`)

- **OpenAPI struct**: Holds all OpenAPI 3.0.3 documentation
- **Config struct**: Initial configuration for OpenAPI
- **Info, Contact, License**: Metadata structures
- **Server struct**: Server configuration (replaces host/basePath)
- **Components struct**: OpenAPI 3.0 components container

#### 2. Generation Layer (`generate.go`)

- JSON OpenAPI documentation generation
- Automatic creation of schema definitions in components/schemas
- Converting endpoints to OpenAPI 3.0 paths format
- Request body handling with proper content types

#### 3. Security Layer (`auth.go`)

- Basic Authentication
- API Key Authentication
- Bearer Authentication (JWT)
- OAuth2 Authentication with multiple flows
- OpenID Connect Authentication

#### 4. External Documentation (`external_docs.go`)

- External documentation references
- Validation for external docs URLs

#### 5. Component Layer (`components/`)

- **Definition**: Generates OpenAPI 3.0 schemas from Go structs with enhanced features
- **Endpoint**: Defines API endpoints with OpenAPI 3.0 features (callbacks, links)
- **Parameter**: OpenAPI 3.0 parameters with enhanced validation
- **Response**: HTTP response structures with content types
- **Fields**: Enhanced struct field parsing
- **MIME**: Content-Type definitions
- **Security**: OpenAPI 3.0 security schemes, types, and OAuth2 flows
- **Tag**: Endpoint grouping with external docs support

## OpenAPI 3.0.3 Specification Compliance

The project is developed in full compliance with the OpenAPI 3.0.3 specification:

### Supported Features ✅

#### Core OpenAPI 3.0 Features

- OpenAPI 3.0.3 Structure
- Multiple Servers Configuration
- Components and References
- Enhanced Security Schemes
- Request Bodies with Content Types
- Response Content Types
- Enhanced Parameter Definitions

#### Security

- Basic Authentication
- API Key Authentication
- Bearer Token Authentication (JWT)
- OAuth2 with Multiple Flows (Authorization Code, Client Credentials, Password, Implicit)
- OpenID Connect Authentication

#### Schema Features

- Enhanced Schema Definitions
- Nullable Types
- ReadOnly and WriteOnly Properties
- Discriminator with Mapping
- Examples and Default Values
- Enhanced Validation Rules

#### Advanced Features

- Callbacks (Webhook definitions)
- Links (Response linking)
- External Documentation
- Server Variables
- Operation-level Security

### Features Under Development 🔄

- Schema Composition (allOf, oneOf, anyOf)
- Advanced Callback Expressions
- Custom Extensions

### Planned Features 🔜

- OpenAPI Validation
- Schema Inheritance
- Advanced Link Templates

## Basic Workflow

1. **Configuration**: Creating OpenAPI instance with `swagno3.New()`
2. **Server Setup**: Adding servers with `AddServer()`
3. **Security Configuration**: Setting up authentication schemes
4. **Endpoint Definition**: Defining API endpoints with `endpoint.New()`
5. **Model Definition**: Defining data models with Go structs
6. **Parameter Definition**: API parameters with enhanced OpenAPI 3.0 features
7. **Response Definition**: API responses with content types
8. **JSON Generation**: OpenAPI documentation generation with `MustToJson()`
9. **Service Integration**: Documentation serving with framework-specific handlers

## Design Principles

### 1. Builder Pattern

Most components use builder pattern for OpenAPI 3.0 features:

```go
endpoint.New(
    endpoint.GET,
    "/users",
    endpoint.WithTags("users"),
    endpoint.WithParams(parameter.IntParam("id", parameter.Path)),
    endpoint.WithSuccessfulReturns([]response.Response{...}),
    endpoint.WithSecurity([]map[string][]string{{"bearerAuth": {}}}),
)
```

### 2. Functional Options Pattern

Enhanced functional options for OpenAPI 3.0 features:

```go
parameter.IntParam("id", parameter.Path,
    parameter.WithRequired(),
    parameter.WithMin(1),
    parameter.WithMax(1000),
    parameter.WithExample(123),
    parameter.WithDescription("User ID"),
)
```

### 3. Enhanced Schema Generation

Automatic OpenAPI 3.0 schema generation from Go structs:

```go
type User struct {
    ID   int     `json:"id"`
    Name string  `json:"name,omitempty"`
    Email *string `json:"email,omitempty"` // Nullable in OpenAPI 3.0
}
// OpenAPI 3.0 schema is automatically generated with proper nullable handling
```

### 4. Interface-Driven Design

Components work through well-defined interfaces:

```go
type Response interface {
    Description() string
    ReturnCode() string
    ContentType() string // Enhanced for OpenAPI 3.0
}
```

### 5. Multiple Server Support

OpenAPI 3.0's multiple server configuration:

```go
openapi.AddServer("https://api.example.com/v1", "Production server")
openapi.AddServer("https://staging.example.com/v1", "Staging server")
openapi.AddServer("http://localhost:8080/v1", "Development server")
```

### 6. Enhanced Security Schemes

OpenAPI 3.0 security with multiple options:

```go
// Bearer Authentication
openapi.SetBearerAuth("JWT", "Bearer authentication using JWT tokens")

// OAuth2 with multiple flows
flows := security.NewOAuthFlows().
    WithAuthorizationCode(
        "https://example.com/oauth/authorize",
        "https://example.com/oauth/token",
        map[string]string{
            "read":  "Read access",
            "write": "Write access",
        },
    )
openapi.SetOAuth2Auth(flows, "OAuth2 authentication")
```

## Key Differences from Swagger 2.0

### 1. Structure Changes

- `openapi: "3.0.3"` instead of `swagger: "2.0"`
- `servers` array instead of `host` and `basePath`
- `components` instead of top-level `definitions`
- Request bodies are separate from parameters

### 2. Enhanced Security

- Bearer token authentication
- OpenID Connect support
- Enhanced OAuth2 flows

### 3. Better Content Type Handling

- Request bodies with content types
- Response content types
- Multiple content type support

### 4. Advanced Features

- Callbacks for webhook definitions
- Links for response relationships
- Enhanced parameter validation

Thanks to this enhanced architecture, Swagno v3 offers full OpenAPI 3.0.3 compliance while maintaining the same ease of use and flexibility as the original Swagno.
