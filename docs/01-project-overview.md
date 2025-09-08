# Project Overview and Architecture

## Project Definition

**Swagno** is a library developed to create Swagger 2.0 documentation in Go projects. The main motto of the project: **"no annotations, no files, no commands"** - meaning no annotations, no file exports, no command execution.

### Key Features

- ✅ **No annotations required**: No need to write special comments in code
- ✅ **No file exports**: No need to create separate swagger files
- ✅ **No command execution**: No additional steps required in build process
- ✅ **Programmatic approach**: All swagger definitions are defined with Go code
- ✅ **Framework agnostic**: Works with Gin, Fiber, net/http and other frameworks

## Project Architecture

### Directory Structure

```
swagno/
├── swagger.go          # Main Swagger struct and basic functions
├── generate.go         # JSON generation and swagger definition creation
├── auth.go            # Security definitions (BasicAuth, OAuth2, ApiKey)
├── go.mod             # Go module definition
├── components/        # Swagger components
│   ├── definition/    # Model definitions and schema generation
│   ├── endpoint/      # API endpoint definitions
│   ├── fields/        # Struct field parsing tools
│   ├── http/response/ # HTTP response structures
│   ├── mime/          # MIME type definitions
│   ├── parameter/     # API parameters
│   ├── security/      # Security structures
│   └── tag/           # Swagger tags
└── example/           # Usage examples
    ├── fiber/         # Fiber framework example
    ├── gin/           # Gin framework example
    ├── http/          # net/http example
    └── models/        # Example model definitions
```

### Main Architecture Components

#### 1. Main Swagger Layer (`swagger.go`)

- **Swagger struct**: Holds all swagger documentation
- **Config struct**: Initial configuration
- **Info, Contact, License**: Metadata structures

#### 2. Generation Layer (`generate.go`)

- JSON swagger documentation generation
- Automatic creation of model definitions
- Converting endpoints to JSON format

#### 3. Security Layer (`auth.go`)

- Basic Authentication
- API Key Authentication
- OAuth2 Authentication

#### 4. Component Layer (`components/`)

- **Definition**: Generates Swagger model schemas from Go structs
- **Endpoint**: Defines API endpoints and converts them to JSON
- **Parameter**: Defines API parameters
- **Response**: Manages HTTP response structures
- **Fields**: Parses struct fields
- **MIME**: Content-Type definitions
- **Security**: Security scopes
- **Tag**: Endpoint grouping

## Swagger 2.0 Specification Compliance

The project is developed in compliance with the Swagger 2.0 specification:

### Supported Features ✅

- Basic Structure
- API Host and Base Path
- Paths and Operations
- Describing Parameters
- Describing Request Body
- Describing Responses
- Authentication (Basic, API Key, OAuth2)
- Enums
- Grouping Operations With Tags

### Features Under Development 🔄

- MIME Types (improvement needed)
- File Upload (improvement needed)

### Planned Features 🔜

- Swagger Extensions
- Swagger Validation

## Basic Workflow

1. **Configuration**: Creating Swagger instance with `swagno.New()`
2. **Endpoint Definition**: Defining API endpoints with `endpoint.New()`
3. **Model Definition**: Defining data models with Go structs
4. **Parameter Definition**: API parameters with `parameter.*Param()` functions
5. **Response Definition**: API responses with `response.New()`
6. **JSON Generation**: Swagger documentation generation with `MustToJson()`
7. **Service Integration**: Documentation serving with framework-specific handlers

## Design Principles

### 1. Builder Pattern

Most components use builder pattern:

```go
endpoint.New(
    endpoint.GET,
    "/users",
    endpoint.WithTags("users"),
    endpoint.WithParams(parameter.IntParam("id", parameter.Path)),
    endpoint.WithSuccessfulReturns([]response.Response{...}),
)
```

### 2. Functional Options Pattern

Functional options pattern for configuration:

```go
parameter.IntParam("id", parameter.Path,
    parameter.WithRequired(),
    parameter.WithMin(1),
    parameter.WithMax(1000),
)
```

### 3. Reflection-Based Schema Generation

Automatic swagger schema generation from Go structs:

```go
type User struct {
    ID   int    `json:"id"`
    Name string `json:"name,omitempty"`
}
// Swagger definition is automatically generated
```

### 4. Interface-Driven Design

Response and other components work through interfaces:

```go
type Response interface {
    Description() string
    ReturnCode() string
}
```

Thanks to this architecture, Swagno offers both a flexible and easy-to-use structure.
