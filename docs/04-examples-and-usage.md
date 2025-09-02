# Usage Examples and Framework Integrations

This section explores how to use Swagno with different Go web frameworks and various usage scenarios.

## 1. Framework Integrations

### 1.1. Net/HTTP Integration (`example/http/server.go`)

Basic HTTP framework integration:

```go
package main

import (
    "fmt"
    "net/http"
    "github.com/go-swagno/swagno"
    "github.com/go-swagno/swagno/components/endpoint"
    "github.com/go-swagno/swagno/components/http/response"
    "github.com/go-swagno/swagno/components/mime"
    "github.com/go-swagno/swagno/components/parameter"
    "github.com/go-swagno/swagno/example/models"
    "github.com/go-swagno/swagno-http/swagger"
)

func main() {
    // Create swagger instance
    sw := swagno.New(swagno.Config{
        Title: "Testing API",
        Version: "v1.0.0"
    })

    // Define endpoints
    endpoints := []*endpoint.EndPoint{
        endpoint.New(
            endpoint.GET,
            "/product",
            endpoint.WithTags("product"),
            endpoint.WithSuccessfulReturns([]response.Response{
                response.New(models.Product{}, "200", "OK")
            }),
            endpoint.WithErrors([]response.Response{
                response.New(models.UnsuccessfulResponse{}, "400", "Bad Request")
            }),
            endpoint.WithDescription("Get all products"),
            endpoint.WithProduce([]mime.MIME{mime.JSON, mime.XML}),
            endpoint.WithConsume([]mime.MIME{mime.JSON}),
            endpoint.WithSummary("Retrieve product list"),
        ),
        endpoint.New(
            endpoint.GET,
            "/product/{id}",
            endpoint.WithTags("product"),
            endpoint.WithParams(
                parameter.IntParam("id", parameter.Path, parameter.WithRequired())
            ),
            endpoint.WithSuccessfulReturns([]response.Response{
                response.New(models.SuccessfulResponse{}, "200", "Request Accepted")
            }),
            endpoint.WithErrors([]response.Response{
                response.New(models.UnsuccessfulResponse{}, "400", "Bad Request")
            }),
            endpoint.WithProduce([]mime.MIME{mime.JSON, mime.XML}),
        ),
    }

    // Add endpoints
    sw.AddEndpoints(endpoints)

    // Add swagger handler
    http.HandleFunc("/swagger/", swagger.SwaggerHandler(sw.MustToJson()))

    fmt.Println("Server is running on http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}
```

**Features:**

- Basic HTTP routing
- Swagger UI served at `/swagger/` path
- JSON and XML content-type support

### 1.2. Gin Framework Integration (`example/gin/server.go`)

Integration with Gin web framework:

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/go-swagno/swagno"
    "github.com/go-swagno/swagno-gin/swagger"
    "github.com/go-swagno/swagno/components/endpoint"
    "github.com/go-swagno/swagno/components/http/response"
    "github.com/go-swagno/swagno/components/mime"
    "github.com/go-swagno/swagno/components/parameter"
    "github.com/go-swagno/swagno/example/models"
)

func main() {
    sw := swagno.New(swagno.Config{Title: "Testing API", Version: "v1.0.0"})

    endpoints := []*endpoint.EndPoint{
        endpoint.New(
            endpoint.GET,
            "/product",
            endpoint.WithTags("product"),
            endpoint.WithSuccessfulReturns([]response.Response{
                response.New(models.Product{}, "200", "OK")
            }),
            endpoint.WithErrors([]response.Response{
                response.New(models.UnsuccessfulResponse{}, "400", "Bad Request")
            }),
            endpoint.WithDescription("Get all products"),
            endpoint.WithProduce([]mime.MIME{mime.JSON, mime.XML}),
            endpoint.WithConsume([]mime.MIME{mime.JSON}),
            endpoint.WithSummary("Retrieve product list"),
        ),
        endpoint.New(
            endpoint.GET,
            "/product/{id}",
            endpoint.WithTags("product"),
            endpoint.WithParams(
                parameter.IntParam("id", parameter.Path, parameter.WithRequired())
            ),
            endpoint.WithSuccessfulReturns([]response.Response{
                response.New(models.SuccessfulResponse{}, "201", "Request Accepted")
            }),
            endpoint.WithErrors([]response.Response{
                response.New(models.UnsuccessfulResponse{}, "400", "Bad Request")
            }),
            endpoint.WithProduce([]mime.MIME{mime.JSON, mime.XML}),
        ),
    }

    sw.AddEndpoints(endpoints)

    app := gin.Default()
    // Gin-specific swagger handler
    app.GET("/swagger/*any", swagger.SwaggerHandler(sw.MustToJson()))

    app.Run()
}
```

**Gin Features:**

- Compatible with Gin routing system
- Wildcard route (`/*any`) usage
- Gin middleware support

### 1.3. Fiber Framework Integration (`example/fiber/server.go`)

Integration with Fiber web framework:

```go
package main

import (
    "fmt"
    "github.com/go-swagno/swagno"
    "github.com/go-swagno/swagno-fiber/swagger"
    "github.com/go-swagno/swagno/components/endpoint"
    "github.com/go-swagno/swagno/components/http/response"
    "github.com/go-swagno/swagno/components/mime"
    "github.com/go-swagno/swagno/components/parameter"
    "github.com/go-swagno/swagno/example/models"
    "github.com/gofiber/fiber/v2"
)

func main() {
    sw := swagno.New(swagno.Config{Title: "Testing API", Version: "v1.0.0"})

    endpoints := []*endpoint.EndPoint{
        endpoint.New(
            endpoint.GET,
            "/product",
            endpoint.WithTags("product"),
            endpoint.WithSuccessfulReturns([]response.Response{
                response.New(models.Product{}, "200", "OK")
            }),
            endpoint.WithDescription("Get all products"),
            endpoint.WithProduce([]mime.MIME{mime.JSON, mime.XML}),
            endpoint.WithConsume([]mime.MIME{mime.JSON}),
            endpoint.WithSummary("Retrieve product list"),
        ),
        endpoint.New(
            endpoint.POST,
            "/product",
            endpoint.WithTags("product"),
            endpoint.WithBody(models.ProductPost{}),
            endpoint.WithSuccessfulReturns([]response.Response{
                response.New(models.Product{}, "200", "OK")
            }),
            endpoint.WithDescription("Create new product"),
            endpoint.WithProduce([]mime.MIME{mime.JSON, mime.XML}),
            endpoint.WithConsume([]mime.MIME{mime.JSON}),
        ),
        endpoint.New(
            endpoint.GET,
            "/product/{id}",
            endpoint.WithTags("product"),
            endpoint.WithParams(
                parameter.IntParam("id", parameter.Path, parameter.WithRequired())
            ),
            endpoint.WithSuccessfulReturns([]response.Response{
                response.New(models.SuccessfulResponse{}, "201", "Request Accepted")
            }),
            endpoint.WithErrors([]response.Response{
                response.New(models.UnsuccessfulResponse{}, "400", "Bad Request")
            }),
            endpoint.WithProduce([]mime.MIME{mime.JSON, mime.XML}),
        ),
    }

    sw.AddEndpoints(endpoints)

    app := fiber.New()
    // Fiber-specific swagger handler with prefix
    swagger.SwaggerHandler(app, sw.MustToJson(), swagger.WithPrefix("/swagger"))

    fmt.Println(app.Listen(":8080"))
}
```

**Fiber Features:**

- Fiber context support
- Swagger path configuration with prefix
- POST endpoint with request body example

## 2. Model Examples (`example/models/`)

### 2.1. Basic Model Structures (`models/products.go`)

```go
package models

import "time"

// Main product model
type Product struct {
    Id             uint64                    `json:"id"`
    Name           string                    `json:"name"`
    MerchantId     uint64                    `json:"merchant_id"`
    CategoryId     *uint64                   `json:"category_id,omitempty"`  // Optional field
    Tags           []uint64                  `json:"tags"`                   // Array of integers
    Images         []*string                 `json:"image_ids"`              // Array of string pointers
    ImagesPtr      *[]string                 `json:"image_ids_ptr"`          // Pointer to string array
    Sizes          []Sizes                   `json:"sizes"`                  // Array of structs
    SizePtrs       []*Sizes                  `json:"size_ptrs"`              // Array of struct pointers
    SaleDate       time.Time                 `json:"sale_date"`              // Time field
    EndDate        *time.Time                `json:"end_date"`               // Optional time field
    Complex        ComplexSuccessfulResponse `json:"complex"`                // Nested struct
    Interface      interface{}               `json:"interface"`              // Interface field
    OmitEmpty      string                    `json:"omitemptytest,omitempty"`
    RequiredField  interface{}               `json:"required_field,omitempty" required:"true"` // Forced required
    EmbeddedStruct EmbeddedStruct            `json:"embedded_struct"`        // Embedded struct
}

// Embedded struct example
type EmbeddedStruct struct {
    Sizes                    // Anonymous embedded
    OtherField int `json:"other_field"`
}

// Sub struct
type Sizes struct {
    Size string `json:"size"`
}

// Model for POST request
type ProductPost struct {
    Name       string  `json:"name" example:"John Smith"`
    MerchantId uint64  `json:"merchant_id" example:"123456"`
    CategoryId *uint64 `json:"category_id,omitempty" example:"123"`
}

// Error response model
type ErrorResponse struct {
    Error   bool   `json:"error"`
    Message string `json:"message"`
}

// Pagination response example
type MerchantPageResponse struct {
    Brochures    []Product `json:"products"`
    MerchantName string    `json:"merchantName"`
}
```

**Model Features:**

- **Optional Fields**: Using `*uint64` pointer types
- **Arrays**: Different array types (`[]uint64`, `[]*string`, `[]Struct`)
- **Time Handling**: Special handling for `time.Time` and `*time.Time`
- **Embedded Structs**: Anonymous embedding support
- **Example Tags**: Example values for Swagger UI
- **Required Tags**: Manual required field marking
- **Interface Fields**: `interface{}` generic type support

### 2.2. Complex Nested Structures (`models/complex.go`)

```go
package models

// Simple object
type Object struct {
    Name string `json:"name" example:"John Smith"`
}

// Nested structure
type Nested struct {
    Objects *[]Object `json:"objects,omitempty"`  // Pointer to array
    Strings *[]string `json:"strings,omitempty"`  // Pointer to string array
}

// Deeply nested structure
type Deeply struct {
    Nested Nested `json:"nested"`
}

// Complex response
type ComplexSuccessfulResponse struct {
    Data *Deeply `json:"deeply"`  // Pointer to deeply nested
}
```

**Complex Structure Features:**

- **Deep Nesting**: 3+ levels of nested structures
- **Pointer Arrays**: `*[]Object` structures
- **Optional Nesting**: Nullable nested objects

### 2.3. Generic Utility Models (`models/generic.go`)

```go
package models

// POST request body
type PostBody struct {
    Name string `json:"name" example:"John Smith"`
    ID   uint64 `json:"id" example:"123456"`
}

// Empty successful response
type EmptySuccessfulResponse struct{}

// Successful response
type SuccessfulResponse struct {
    ID string `json:"ID" example:"1234-1234-1234-1234"`
}

// Generic error responses
type UnsuccessfulResponse struct {
    ErrorField1 string `json:"error_msg1"`
}

type PageNotFound struct {
    ErrorMsg2 string `json:"error_msg2"`
}
```

## 3. Advanced Usage Scenarios

### 3.1. Authentication Usage

```go
// Basic Auth
sw.SetBasicAuth("Basic authentication required")

// API Key Auth
sw.SetApiKeyAuth("X-API-Key", "header", "API key required for access")

// OAuth2 Auth
scopes := security.Scopes(
    security.Scope("read:products", "Read product data"),
    security.Scope("write:products", "Create and update products"),
    security.Scope("admin", "Administrative access"),
)
sw.SetOAuth2Auth(
    "oauth2",
    "password",
    "",
    "http://localhost:8080/oauth/token",
    scopes,
    "OAuth2 authentication"
)
```

### 3.2. Endpoint Security Definition

```go
endpoint.New(
    endpoint.POST,
    "/admin/products",
    endpoint.WithTags("admin", "products"),
    endpoint.WithSecurity([]map[string][]string{
        {"oauth2": {"admin"}},           // OAuth2 admin scope required
        {"X-API-Key": {}},               // Or API key required
    }),
    // ... other options
)
```

### 3.3. Complex Parameter Usage

```go
endpoint.New(
    endpoint.GET,
    "/products",
    endpoint.WithParams(
        // Query parameters
        parameter.StrParam("search", parameter.Query,
            parameter.WithDescription("Search term for products"),
            parameter.WithMinLen(2),
            parameter.WithMaxLen(100),
        ),
        parameter.IntParam("page", parameter.Query,
            parameter.WithDefault(1),
            parameter.WithMin(1),
            parameter.WithMax(1000),
        ),
        parameter.IntParam("limit", parameter.Query,
            parameter.WithDefault(10),
            parameter.WithMin(1),
            parameter.WithMax(100),
        ),
        // Enum parameter
        parameter.StrEnumParam("sort", parameter.Query,
            []string{"name", "price", "date"},
            parameter.WithDefault("name"),
        ),
        // Header parameter
        parameter.StrParam("X-Client-Version", parameter.Header,
            parameter.WithRequired(),
            parameter.WithPattern("^\\d+\\.\\d+\\.\\d+$"),
        ),
    ),
    // ... other options
)
```

### 3.4. File Upload Endpoint

```go
endpoint.New(
    endpoint.POST,
    "/products/{id}/image",
    endpoint.WithTags("products", "upload"),
    endpoint.WithParams(
        parameter.IntParam("id", parameter.Path, parameter.WithRequired()),
        parameter.FileParam("image",
            parameter.WithRequired(),
            parameter.WithDescription("Product image file"),
        ),
        parameter.StrParam("description", parameter.Form,
            parameter.WithDescription("Image description"),
        ),
    ),
    endpoint.WithConsume([]mime.MIME{mime.MULTIFORM}),
    endpoint.WithSuccessfulReturns([]response.Response{
        response.New(UploadResponse{}, "200", "Image uploaded successfully"),
    }),
)
```

### 3.5. Multiple Response Types

```go
endpoint.New(
    endpoint.GET,
    "/products/{id}",
    endpoint.WithSuccessfulReturns([]response.Response{
        response.New(Product{}, "200", "Product found"),
        response.New(ProductWithReviews{}, "200", "Product with reviews"),
    }),
    endpoint.WithErrors([]response.Response{
        response.New(ErrorResponse{}, "400", "Invalid product ID"),
        response.New(ErrorResponse{}, "404", "Product not found"),
        response.New(ErrorResponse{}, "500", "Internal server error"),
    }),
)
```

### 3.6. Tag Usage

```go
// Define tags
sw.AddTags(
    tag.New("products", "Product management operations",
        tag.WithExternalDocs("Product API Guide", "https://docs.example.com/products"),
    ),
    tag.New("users", "User management operations"),
    tag.New("admin", "Administrative operations"),
)

// Use in endpoints
endpoint.New(
    endpoint.GET,
    "/admin/products",
    endpoint.WithTags("products", "admin"),  // Multiple tags
    // ...
)
```

### 3.7. Swagger Export

```go
// Get as JSON string
jsonDoc := sw.MustToJson()

// Save to file
sw.ExportSwaggerDocs("swagger.json")

// With error handling
jsonDoc, err := sw.ToJson()
if err != nil {
    log.Fatal("Swagger generation failed:", err)
}
```

## 4. Best Practices

### 4.1. Model Design

```go
// ✅ Good: Clear field definitions
type User struct {
    ID       uint64     `json:"id"`
    Name     string     `json:"name" example:"John Doe"`
    Email    string     `json:"email" example:"john@example.com"`
    IsActive *bool      `json:"is_active,omitempty"`  // Optional boolean
    Created  time.Time  `json:"created"`
}

// ❌ Bad: Ambiguous interface{} usage
type BadUser struct {
    Data interface{} `json:"data"`  // Appears as "Ambiguous Type" in Swagger
}
```

### 4.2. Endpoint Organization

```go
// ✅ Good: Logical groups
var userEndpoints = []*endpoint.EndPoint{
    endpoint.New(endpoint.GET, "/users", ...),
    endpoint.New(endpoint.POST, "/users", ...),
    endpoint.New(endpoint.GET, "/users/{id}", ...),
}

var productEndpoints = []*endpoint.EndPoint{
    endpoint.New(endpoint.GET, "/products", ...),
    endpoint.New(endpoint.POST, "/products", ...),
}

sw.AddEndpoints(userEndpoints)
sw.AddEndpoints(productEndpoints)
```

### 4.3. Error Handling

```go
// ✅ Good: Consistent error responses
type APIError struct {
    Code    int    `json:"code" example:"400"`
    Message string `json:"message" example:"Validation failed"`
    Details string `json:"details,omitempty" example:"Name is required"`
}

// Use same error format in all endpoints
endpoint.WithErrors([]response.Response{
    response.New(APIError{}, "400", "Bad Request"),
    response.New(APIError{}, "404", "Not Found"),
    response.New(APIError{}, "500", "Internal Server Error"),
})
```

These examples and best practices will help you effectively use Swagno's powerful features.
