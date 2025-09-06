# Migration Guide: Swagger 2.0 to OpenAPI 3.0

This guide provides detailed instructions for migrating from Swagno (Swagger 2.0) to Swagno v3 (OpenAPI 3.0.3).

## 1. Overview of Changes

### Major Differences

| Feature                   | Swagger 2.0 (Swagno)          | OpenAPI 3.0 (Swagno v3)              |
| ------------------------- | ----------------------------- | ------------------------------------ |
| **Specification Version** | `swagger: "2.0"`              | `openapi: "3.0.3"`                   |
| **Import Path**           | `github.com/go-swagno/swagno` | `github.com/go-swagno/swagno/v3`     |
| **Constructor**           | `swagno.New()`                | `swagno3.New()`                      |
| **Server Configuration**  | `host` + `basePath`           | `servers` array                      |
| **Schema Location**       | `#/definitions/`              | `#/components/schemas/`              |
| **Request Bodies**        | Part of parameters            | Separate `requestBody`               |
| **Security Schemes**      | Limited options               | Enhanced with Bearer, OpenID Connect |
| **Content Types**         | `consumes`/`produces`         | `content` with media types           |

## 2. Installation and Import Changes

### Old (Swagger 2.0)

```go
import "github.com/go-swagno/swagno"
```

### New (OpenAPI 3.0)

```go
import swagno3 "github.com/go-swagno/swagno/v3"
```

### Go Module Update

Update your `go.mod` file:

```bash
go get github.com/go-swagno/swagno/v3
```

## 3. Basic Setup Migration

### Old (Swagger 2.0)

```go
sw := swagno.New(swagno.Config{
    Title:       "My API",
    Version:     "v1.0.0",
    Description: "API description",
    Host:        "api.example.com",
    Path:        "/v1",
})
```

### New (OpenAPI 3.0)

```go
openapi := swagno3.New(swagno3.Config{
    Title:       "My API",
    Version:     "v1.0.0",
    Description: "API description",
})

// Replace host + basePath with servers
openapi.AddServer("https://api.example.com/v1", "Production server")
openapi.AddServer("https://staging.example.com/v1", "Staging server")
```

## 4. Security Migration

### Basic Authentication

#### Old (Swagger 2.0)

```go
sw.SetBasicAuth("Basic authentication required")
```

#### New (OpenAPI 3.0)

```go
openapi.SetBasicAuth("Basic authentication required")
```

_No changes required for Basic Auth_

### API Key Authentication

#### Old (Swagger 2.0)

```go
sw.SetApiKeyAuth("X-API-Key", "header", "API key authentication")
```

#### New (OpenAPI 3.0)

```go
openapi.SetApiKeyAuth("X-API-Key", security.Header, "API key authentication")
```

_No changes required for API Key Auth_

### Bearer Authentication (New in OpenAPI 3.0)

#### Old (Swagger 2.0)

```go
// Not directly supported, had to use custom security definition
sw.SecurityDefinitions["bearerAuth"] = securityDefinition{
    Type:        "apiKey",
    Name:        "Authorization",
    In:          "header",
    Description: "Bearer token",
}
```

#### New (OpenAPI 3.0)

```go
openapi.SetBearerAuth("JWT", "Bearer authentication using JWT tokens")
```

### OAuth2 Authentication

#### Old (Swagger 2.0)

```go
// OAuth2 scopes were handled differently in Swagger 2.0
// Usually defined as part of security definitions
```

#### New (OpenAPI 3.0)

```go
import "github.com/go-swagno/swagno/v3/components/security"

flows := security.NewOAuthFlows().
    WithPassword(
        "http://localhost:8080/oauth/token",
        map[string]string{
            "read":  "Read access",
            "write": "Write access",
        },
    ).
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

### OpenID Connect (New in OpenAPI 3.0)

#### New Feature

```go
openapi.SetOpenIdConnectAuth(
    "https://example.com/.well-known/openid_configuration",
    "OpenID Connect authentication",
)
```

## 5. Endpoint Migration

### Basic Endpoint Structure

#### Old (Swagger 2.0)

```go
endpoint.New(
    endpoint.GET,
    "/users/{id}",
    endpoint.WithTags("users"),
    endpoint.WithParams(parameter.IntParam("id", parameter.Path, parameter.WithRequired())),
    endpoint.WithSuccessfulReturns([]response.Response{response.New(User{}, "200", "OK")}),
    endpoint.WithProduce([]mime.MIME{mime.JSON}),
    endpoint.WithConsume([]mime.MIME{mime.JSON}),
)
```

#### New (OpenAPI 3.0)

```go
endpoint.New(
    endpoint.GET,
    "/users/{id}",
    endpoint.WithTags("users"),
    endpoint.WithParams(parameter.IntParam("id", parameter.Path, parameter.WithRequired())),
    endpoint.WithSuccessfulReturns([]response.Response{response.New(User{}, "200", "User found")}),
    // Produce/Consume are automatically handled based on response content types
)
```

### Request Body Handling

#### Old (Swagger 2.0)

```go
endpoint.New(
    endpoint.POST,
    "/users",
    endpoint.WithBody(CreateUserRequest{}),
    endpoint.WithConsume([]mime.MIME{mime.JSON}),
)
```

#### New (OpenAPI 3.0)

```go
// Simple approach (same as before)
endpoint.New(
    endpoint.POST,
    "/users",
    endpoint.WithBody(CreateUserRequest{}),
    // Content types are automatically inferred
)

// Advanced approach with multiple content types
endpoint.New(
    endpoint.POST,
    "/users",
    endpoint.WithRequestBody(endpoint.RequestBody{
        Description: "User creation data",
        Required:    true,
        Content: map[string]endpoint.MediaType{
            "application/json": {
                Schema: definition.SchemaFromStruct(CreateUserRequest{}),
            },
            "application/xml": {
                Schema: definition.SchemaFromStruct(CreateUserRequest{}),
            },
        },
    }),
)
```

## 6. Parameter Migration

### Basic Parameters

#### Old (Swagger 2.0)

```go
parameter.IntParam("id", parameter.Path,
    parameter.WithRequired(),
    parameter.WithDescription("User ID"),
    parameter.WithMin(1),
    parameter.WithMax(999999),
)
```

#### New (OpenAPI 3.0)

```go
parameter.IntParam("id", parameter.Path,
    parameter.WithRequired(),
    parameter.WithDescription("User ID"),
    parameter.WithExample(123),          // Enhanced examples support
    parameter.WithMin(1),
    parameter.WithMax(999999),
)
```

### Enhanced Parameter Features (New in OpenAPI 3.0)

```go
// Parameter styles
parameter.ArrayParam("tags", parameter.Query,
    parameter.WithStyle("form"),           // New: Parameter styles
    parameter.WithExplode(true),           // New: Explode behavior
    parameter.WithDescription("Filter by tags"),
)

// Multiple examples
parameter.StringParam("format", parameter.Query,
    parameter.WithExamples(map[string]parameter.Example{  // New: Multiple examples
        "json": {
            Summary: "JSON format",
            Value:   "json",
        },
        "xml": {
            Summary: "XML format",
            Value:   "xml",
        },
    }),
)
```

## 7. Response Migration

### Basic Responses

#### Old (Swagger 2.0)

```go
response.New(User{}, "200", "OK")
```

#### New (OpenAPI 3.0)

```go
response.New(User{}, "200", "User found successfully")
// More descriptive messages recommended
```

### Enhanced Response Features (New in OpenAPI 3.0)

```go
// Basic response (headers would be added using enhanced response features)
response.New([]User{}, "200", "Users retrieved")

// Response with links
response.NewWithLinks(User{}, "201", "User created",
    map[string]endpoint.Link{
        "GetUser": endpoint.NewEnhancedLink().
            SetOperationId("getUserById").
            AddParameter("id", "$response.body#/id"),
    },
)

// Response with multiple content types
response.NewWithContent("200", "User found", map[string]response.MediaType{
    "application/json": {
        Schema: definition.SchemaFromStruct(User{}),
    },
    "application/xml": {
        Schema: definition.SchemaFromStruct(User{}),
    },
})
```

## 8. Schema Migration

### Basic Schema Handling

#### Old (Swagger 2.0)

```go
type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name,omitempty"`
    Email string `json:"email"`
}
```

#### New (OpenAPI 3.0)

```go
type User struct {
    ID    int     `json:"id" example:"123" description:"User ID"`
    Name  string  `json:"name,omitempty" example:"John Doe" description:"User name"`
    Email *string `json:"email,omitempty" example:"john@example.com" format:"email" description:"User email"`
    // Note: *string automatically becomes nullable: true in OpenAPI 3.0
}
```

### Enhanced Schema Features (New in OpenAPI 3.0)

```go
type User struct {
    ID        int       `json:"id" readonly:"true" description:"User ID"`
    Password  string    `json:"password" writeonly:"true" description:"User password"`
    Name      string    `json:"name" description:"User name" minLength:"1" maxLength:"100"`
    CreatedAt time.Time `json:"created_at" format:"date-time" readonly:"true"`
    IsActive  *bool     `json:"is_active,omitempty"` // Nullable boolean
}
```

## 9. Tag Migration

#### Old (Swagger 2.0)

```go
sw.AddTags(
    tag.New("users", "User operations"),
)
```

#### New (OpenAPI 3.0)

```go
openapi.AddTags(
    tag.New("users", "User operations",
        tag.WithExternalDocs("https://docs.example.com/users", "User docs"), // New feature
    ),
)
```

## 10. Complete Migration Example

### Before (Swagger 2.0)

```go
package main

import (
    "github.com/go-swagno/swagno"
    "github.com/go-swagno/swagno/components/endpoint"
    "github.com/go-swagno/swagno/components/http/response"
    "github.com/go-swagno/swagno/components/parameter"
    "github.com/go-swagno/swagno/components/mime"
)

func main() {
    sw := swagno.New(swagno.Config{
        Title:   "User API",
        Version: "v1.0.0",
        Host:    "api.example.com",
        Path:    "/v1",
    })

    sw.SetBasicAuth("Basic auth required")

    endpoints := []*endpoint.EndPoint{
        endpoint.New(
            endpoint.GET,
            "/users/{id}",
            endpoint.WithParams(parameter.IntParam("id", parameter.Path, parameter.WithRequired())),
            endpoint.WithSuccessfulReturns([]response.Response{response.New(User{}, "200", "OK")}),
            endpoint.WithProduce([]mime.MIME{mime.JSON}),
        ),
        endpoint.New(
            endpoint.POST,
            "/users",
            endpoint.WithBody(CreateUserRequest{}),
            endpoint.WithConsume([]mime.MIME{mime.JSON}),
            endpoint.WithProduce([]mime.MIME{mime.JSON}),
            endpoint.WithSuccessfulReturns([]response.Response{response.New(User{}, "201", "Created")}),
        ),
    }

    sw.AddEndpoints(endpoints)
    jsonDoc := sw.MustToJson()
}
```

### After (OpenAPI 3.0)

```go
package main

import (
    swagno3 "github.com/go-swagno/swagno/v3"
    "github.com/go-swagno/swagno/v3/components/endpoint"
    "github.com/go-swagno/swagno/v3/components/http/response"
    "github.com/go-swagno/swagno/v3/components/parameter"
    "github.com/go-swagno/swagno/v3/components/security"
    "github.com/go-swagno/swagno/v3/components/tag"
)

func main() {
    openapi := swagno3.New(swagno3.Config{
        Title:       "User API",
        Version:     "v1.0.0",
        Description: "User management API",
        Contact: &swagno3.Contact{
            Name:  "API Support",
            Email: "support@example.com",
        },
    })

    openapi.AddServer("https://api.example.com/v1", "Production server")
    openapi.SetBasicAuth("Basic authentication required")
    openapi.SetBearerAuth("JWT", "JWT bearer authentication")

    openapi.AddTags(
        tag.New("users", "User management operations"),
    )

    endpoints := []*endpoint.EndPoint{
        endpoint.New(
            endpoint.GET,
            "/users/{id}",
            endpoint.WithTags("users"),
            endpoint.WithSummary("Get user by ID"),
            endpoint.WithDescription("Retrieve a user by their ID"),
            endpoint.WithParams(
                parameter.IntParam("id", parameter.Path,
                    parameter.WithRequired(),
                    parameter.WithDescription("User ID"),
                    parameter.WithExample(123),
                ),
            ),
            endpoint.WithSuccessfulReturns([]response.Response{
                response.New(User{}, "200", "User found"),
            }),
            endpoint.WithSecurity([]map[string][]string{
                {"bearerAuth": {}},
                {"basicAuth": {}},
            }),
        ),
        endpoint.New(
            endpoint.POST,
            "/users",
            endpoint.WithTags("users"),
            endpoint.WithSummary("Create user"),
            endpoint.WithDescription("Create a new user"),
            endpoint.WithBody(CreateUserRequest{}),
            endpoint.WithSuccessfulReturns([]response.Response{
                response.New(User{}, "201", "User created successfully"),
            }),
            endpoint.WithSecurity([]map[string][]string{
                {"bearerAuth": {}},
            }),
        ),
    }

    openapi.AddEndpoints(endpoints)
    jsonDoc := openapi.MustToJson()
    openapi.ExportOpenAPIDocs("openapi.json")
}
```

## 11. Migration Checklist

### ✅ Required Changes

- [ ] Update import path to `github.com/go-swagno/swagno/v3`
- [ ] Change constructor from `swagno.New()` to `swagno3.New()`
- [ ] Replace `Host` and `Path` config with `AddServer()` calls
- [ ] Update Go module with `go get github.com/go-swagno/swagno/v3`

### ✅ Recommended Enhancements

- [ ] Add Bearer authentication if using JWT tokens
- [ ] Enhance parameter definitions with examples and descriptions
- [ ] Add tags to endpoints for better organization
- [ ] Improve response descriptions
- [ ] Add contact and license information
- [ ] Consider adding external documentation links
- [ ] Use enhanced schema features (nullable, readonly, writeonly)

### ✅ Optional Advanced Features

- [ ] Implement callbacks for webhook documentation
- [ ] Use response links for related operations
- [ ] Add multiple content types for responses
- [ ] Use parameter styles for better array/object handling
- [ ] Add operation-level security requirements
- [ ] Implement multiple OAuth2 flows

## 12. Common Migration Issues

### Issue 1: Schema References

**Problem:** References change from `#/definitions/` to `#/components/schemas/`

**Solution:** This is handled automatically by Swagno v3. No action required.

### Issue 2: Request Body Handling

**Problem:** Request bodies are now separate from parameters

**Solution:** Use `WithBody()` for simple cases, or `WithRequestBody()` for advanced content type handling.

### Issue 3: Content Type Handling

**Problem:** `consumes` and `produces` are replaced with content-type specific schemas

**Solution:** Swagno v3 handles this automatically. For advanced cases, use the new content-type specific response methods.

### Issue 4: Security Scheme Names

**Problem:** Security scheme references might need updates

**Solution:** Ensure security scheme names match between definition and usage:

```go
// Define
openapi.SetBearerAuth("JWT", "Bearer authentication")

// Use
endpoint.WithSecurity([]map[string][]string{
    {"bearerAuth": {}}, // Must match the scheme name
})
```

## 13. Testing Your Migration

### Validation Steps

1. **Generate Documentation:**

   ```go
   jsonDoc := openapi.MustToJson()
   openapi.ExportOpenAPIDocs("openapi.json")
   ```

2. **Validate with OpenAPI Tools:**

   ```bash
   # Install OpenAPI validator
   npm install -g @apidevtools/swagger-parser

   # Validate your OpenAPI document
   swagger-parser validate openapi.json
   ```

3. **Test with Swagger UI:**

   - Load your `openapi.json` in Swagger UI
   - Verify all endpoints are displayed correctly
   - Test the "Try it out" functionality

4. **Compare Schemas:**
   - Compare the generated schema with your Swagger 2.0 version
   - Ensure all models are correctly represented
   - Verify that nullable fields are properly marked

By following this migration guide, you can successfully upgrade from Swagno (Swagger 2.0) to Swagno v3 (OpenAPI 3.0.3) while taking advantage of the enhanced features and improved specification compliance.
