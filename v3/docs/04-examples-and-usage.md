# Examples and Usage Guide

This section provides comprehensive examples and usage patterns for Swagno v3 with OpenAPI 3.0.3 features.

## 1. Basic API Example

### Simple CRUD API

```go
package main

import (
    "fmt"

    swagno3 "github.com/go-swagno/swagno/v3"
    "github.com/go-swagno/swagno/v3/components/endpoint"
    "github.com/go-swagno/swagno/v3/components/http/response"
    "github.com/go-swagno/swagno/v3/components/parameter"
    "github.com/go-swagno/swagno/v3/components/security"
    "github.com/go-swagno/swagno/v3/components/tag"
)

// User model
type User struct {
    ID       uint64 `json:"id" example:"1" description:"User unique identifier"`
    Name     string `json:"name" example:"John Doe" description:"User full name"`
    Email    string `json:"email" example:"john@example.com" format:"email" description:"User email address"`
    IsActive *bool  `json:"is_active,omitempty" example:"true" description:"Whether user is active"`
}

// CreateUserRequest model
type CreateUserRequest struct {
    Name  string `json:"name" example:"John Doe" minLength:"1" maxLength:"100" description:"User full name"`
    Email string `json:"email" example:"john@example.com" format:"email" description:"User email address"`
}

// ErrorResponse model
type ErrorResponse struct {
    Error   bool   `json:"error" example:"true" description:"Indicates if this is an error"`
    Message string `json:"message" example:"User not found" description:"Error message"`
    Code    string `json:"code,omitempty" example:"USER_NOT_FOUND" description:"Error code"`
}

func main() {
    // Create OpenAPI 3.0 instance
    openapi := swagno3.New(swagno3.Config{
        Title:       "User Management API",
        Version:     "v1.0.0",
        Summary:     "A simple API for managing users",
        Description: "This API allows you to create, read, update, and delete users. It demonstrates OpenAPI 3.0.3 features with Swagno v3.",
        Contact: &swagno3.Contact{
            Name:  "API Support Team",
            Email: "support@example.com",
            URL:   "https://example.com/support",
        },
        License: &swagno3.License{
            Name: "MIT",
            URL:  "https://opensource.org/licenses/MIT",
        },
        TermsOfService: "https://example.com/terms",
        ExternalDocs: swagno3.NewExternalDocs(
            "https://docs.example.com/api",
            "Complete API Documentation",
        ),
    })

    // Add servers
    openapi.AddServer("https://api.example.com/v1", "Production server")
    openapi.AddServer("https://staging-api.example.com/v1", "Staging server")
    openapi.AddServer("http://localhost:8080/v1", "Development server")

    // Add security schemes
    openapi.SetBearerAuth("JWT", "Bearer authentication using JWT tokens")
    openapi.SetApiKeyAuth("X-API-Key", security.Header, "API key authentication")

    // Add tags
    openapi.AddTags(
        tag.New("users", "User management operations",
            tag.WithExternalDocs("https://docs.example.com/users", "User API docs"),
        ),
        tag.New("health", "Health check operations"),
    )

    // Define endpoints
    endpoints := []*endpoint.EndPoint{
        // GET /health - Health check
        endpoint.New(
            endpoint.GET,
            "/health",
            endpoint.WithTags("health"),
            endpoint.WithSummary("Health check"),
            endpoint.WithDescription("Check if the API is running"),
            endpoint.WithSuccessfulReturns([]response.Response{
                response.New(map[string]string{"status": "ok"}, "200", "API is healthy"),
            }),
        ),

        // GET /users - List users
        endpoint.New(
            endpoint.GET,
            "/users",
            endpoint.WithTags("users"),
            endpoint.WithSummary("List users"),
            endpoint.WithDescription("Retrieve a paginated list of all users"),
            endpoint.WithParams(
                parameter.IntParam("page", parameter.Query,
                    parameter.WithDescription("Page number for pagination"),
                    parameter.WithDefault(1),
                    parameter.WithMin(1),
                    parameter.WithExample(1),
                ),
                parameter.IntParam("limit", parameter.Query,
                    parameter.WithDescription("Number of users per page"),
                    parameter.WithDefault(10),
                    parameter.WithMin(1),
                    parameter.WithMax(100),
                    parameter.WithExample(10),
                ),
                parameter.StringParam("search", parameter.Query,
                    parameter.WithDescription("Search users by name or email"),
                    parameter.WithExample("john"),
                ),
                parameter.BoolParam("active", parameter.Query,
                    parameter.WithDescription("Filter by active status"),
                    parameter.WithExample(true),
                ),
            ),
            endpoint.WithSuccessfulReturns([]response.Response{
                response.New([]User{}, "200", "List of users"),
            }),
            endpoint.WithErrors([]response.Response{
                response.New(ErrorResponse{}, "400", "Bad Request"),
                response.New(ErrorResponse{}, "401", "Unauthorized"),
                response.New(ErrorResponse{}, "500", "Internal Server Error"),
            }),
            endpoint.WithSecurity([]map[security.SecuritySchemeName][]string{
                {security.BearerAuth: {}},
                {security.APIKeyAuth: {}},
            }),
        ),

        // POST /users - Create user
        endpoint.New(
            endpoint.POST,
            "/users",
            endpoint.WithTags("users"),
            endpoint.WithSummary("Create user"),
            endpoint.WithDescription("Create a new user with the provided information"),
            endpoint.WithBody(CreateUserRequest{}),
            endpoint.WithSuccessfulReturns([]response.Response{
                response.New(User{}, "201", "User created successfully"),
            }),
            endpoint.WithErrors([]response.Response{
                response.New(ErrorResponse{}, "400", "Bad Request - Invalid input"),
                response.New(ErrorResponse{}, "401", "Unauthorized"),
                response.New(ErrorResponse{}, "409", "Conflict - User already exists"),
                response.New(ErrorResponse{}, "500", "Internal Server Error"),
            }),
            endpoint.WithSecurity([]map[security.SecuritySchemeName][]string{
                {security.BearerAuth: {}},
            }),
        ),

        // GET /users/{id} - Get user by ID
        endpoint.New(
            endpoint.GET,
            "/users/{id}",
            endpoint.WithTags("users"),
            endpoint.WithSummary("Get user by ID"),
            endpoint.WithDescription("Retrieve a specific user by their unique identifier"),
            endpoint.WithParams(
                parameter.IntParam("id", parameter.Path,
                    parameter.WithRequired(),
                    parameter.WithDescription("User unique identifier"),
                    parameter.WithMin(1),
                    parameter.WithExample(123),
                ),
            ),
            endpoint.WithSuccessfulReturns([]response.Response{
                response.New(User{}, "200", "User found"),
            }),
            endpoint.WithErrors([]response.Response{
                response.New(ErrorResponse{}, "400", "Bad Request - Invalid ID"),
                response.New(ErrorResponse{}, "401", "Unauthorized"),
                response.New(ErrorResponse{}, "404", "User not found"),
                response.New(ErrorResponse{}, "500", "Internal Server Error"),
            }),
            endpoint.WithSecurity([]map[security.SecuritySchemeName][]string{
                {security.BearerAuth: {}},
                {security.APIKeyAuth: {}},
            }),
        ),

        // PUT /users/{id} - Update user
        endpoint.New(
            endpoint.PUT,
            "/users/{id}",
            endpoint.WithTags("users"),
            endpoint.WithSummary("Update user"),
            endpoint.WithDescription("Update an existing user with new information"),
            endpoint.WithParams(
                parameter.IntParam("id", parameter.Path,
                    parameter.WithRequired(),
                    parameter.WithDescription("User unique identifier"),
                    parameter.WithMin(1),
                    parameter.WithExample(123),
                ),
            ),
            endpoint.WithBody(CreateUserRequest{}),
            endpoint.WithSuccessfulReturns([]response.Response{
                response.New(User{}, "200", "User updated successfully"),
            }),
            endpoint.WithErrors([]response.Response{
                response.New(ErrorResponse{}, "400", "Bad Request - Invalid input"),
                response.New(ErrorResponse{}, "401", "Unauthorized"),
                response.New(ErrorResponse{}, "404", "User not found"),
                response.New(ErrorResponse{}, "500", "Internal Server Error"),
            }),
            endpoint.WithSecurity([]map[security.SecuritySchemeName][]string{
                {security.BearerAuth: {}},
            }),
        ),

        // DELETE /users/{id} - Delete user
        endpoint.New(
            endpoint.DELETE,
            "/users/{id}",
            endpoint.WithTags("users"),
            endpoint.WithSummary("Delete user"),
            endpoint.WithDescription("Delete an existing user by their ID"),
            endpoint.WithParams(
                parameter.IntParam("id", parameter.Path,
                    parameter.WithRequired(),
                    parameter.WithDescription("User unique identifier"),
                    parameter.WithMin(1),
                    parameter.WithExample(123),
                ),
            ),
            endpoint.WithSuccessfulReturns([]response.Response{
                response.New(struct{}{}, "204", "User deleted successfully"),
            }),
            endpoint.WithErrors([]response.Response{
                response.New(ErrorResponse{}, "400", "Bad Request - Invalid ID"),
                response.New(ErrorResponse{}, "401", "Unauthorized"),
                response.New(ErrorResponse{}, "404", "User not found"),
                response.New(ErrorResponse{}, "500", "Internal Server Error"),
            }),
            endpoint.WithSecurity([]map[security.SecuritySchemeName][]string{
                {security.BearerAuth: {}},
            }),
        ),
    }

    // Add endpoints to OpenAPI
    openapi.AddEndpoints(endpoints)

    // Generate and save documentation
    jsonDoc := openapi.MustToJson()
    openapi.ExportOpenAPIDocs("openapi.json")

    fmt.Println("OpenAPI 3.0.3 documentation generated successfully!")
}
```

## 2. Advanced Features Example

### E-commerce API with Advanced OpenAPI 3.0 Features

```go
package main

import (
    "time"

    swagno3 "github.com/go-swagno/swagno/v3"
    "github.com/go-swagno/swagno/v3/components/definition"
    "github.com/go-swagno/swagno/v3/components/endpoint"
    "github.com/go-swagno/swagno/v3/components/http/response"
    "github.com/go-swagno/swagno/v3/components/parameter"
    "github.com/go-swagno/swagno/v3/components/security"
)

// Product models
type Product struct {
    ID          uint64    `json:"id" example:"1" description:"Product ID"`
    Name        string    `json:"name" example:"Laptop" description:"Product name"`
    Description string    `json:"description" example:"High-performance laptop" description:"Product description"`
    Price       float64   `json:"price" example:"999.99" description:"Product price"`
    Category    Category  `json:"category" description:"Product category"`
    InStock     bool      `json:"in_stock" example:"true" description:"Stock availability"`
    CreatedAt   time.Time `json:"created_at" example:"2023-01-01T00:00:00Z" format:"date-time" description:"Creation timestamp"`
    UpdatedAt   time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z" format:"date-time" description:"Last update timestamp"`
}

type Category struct {
    ID   uint64 `json:"id" example:"1" description:"Category ID"`
    Name string `json:"name" example:"Electronics" description:"Category name"`
    Slug string `json:"slug" example:"electronics" description:"Category URL slug"`
}

type CreateProductRequest struct {
    Name        string  `json:"name" example:"Laptop" minLength:"1" maxLength:"200" description:"Product name"`
    Description string  `json:"description" example:"High-performance laptop" maxLength:"1000" description:"Product description"`
    Price       float64 `json:"price" example:"999.99" minimum:"0" description:"Product price"`
    CategoryID  uint64  `json:"category_id" example:"1" description:"Category ID"`
}

// Order models
type Order struct {
    ID          uint64      `json:"id" example:"1" description:"Order ID"`
    UserID      uint64      `json:"user_id" example:"1" description:"User ID"`
    Items       []OrderItem `json:"items" description:"Order items"`
    Total       float64     `json:"total" example:"1999.98" description:"Total order amount"`
    Status      OrderStatus `json:"status" description:"Order status"`
    CreatedAt   time.Time   `json:"created_at" format:"date-time" description:"Order creation time"`
}

type OrderItem struct {
    ProductID uint64  `json:"product_id" example:"1" description:"Product ID"`
    Quantity  int     `json:"quantity" example:"2" minimum:"1" description:"Item quantity"`
    Price     float64 `json:"price" example:"999.99" description:"Unit price"`
}

type OrderStatus string

const (
    OrderStatusPending   OrderStatus = "pending"
    OrderStatusPaid      OrderStatus = "paid"
    OrderStatusShipped   OrderStatus = "shipped"
    OrderStatusDelivered OrderStatus = "delivered"
    OrderStatusCancelled OrderStatus = "cancelled"
)

// Webhook models
type WebhookPayload struct {
    Event     string      `json:"event" example:"order.created" description:"Event name"`
    Timestamp time.Time   `json:"timestamp" format:"date-time" description:"Event timestamp"`
    Data      interface{} `json:"data" description:"Event data"`
}

func main() {
    // Create OpenAPI instance
    openapi := swagno3.New(swagno3.Config{
        Title:       "E-commerce API",
        Version:     "v2.0.0",
        Summary:     "Advanced e-commerce API with OpenAPI 3.0 features",
        Description: "A comprehensive e-commerce API demonstrating advanced OpenAPI 3.0.3 features including callbacks, links, multiple security schemes, and complex schemas.",
        Contact: &swagno3.Contact{
            Name:  "E-commerce API Team",
            Email: "api@ecommerce.com",
            URL:   "https://ecommerce.com/api-support",
        },
        License: &swagno3.License{
            Name: "Apache 2.0",
            URL:  "https://www.apache.org/licenses/LICENSE-2.0.html",
        },
    })

    // Add multiple servers with variables
    openapi.AddServer("https://{environment}.ecommerce.com/{version}", "Configurable server")

    // Advanced security schemes
    openapi.SetBearerAuth("JWT", "JWT bearer token authentication")
    openapi.SetApiKeyAuth("X-API-Key", security.Header, "API key for public endpoints")

    // OAuth2 with multiple flows
    flows := security.NewOAuthFlows().
        WithAuthorizationCode(
            "https://auth.ecommerce.com/oauth/authorize",
            "https://auth.ecommerce.com/oauth/token",
            map[string]string{
                "read:products": "Read products",
                "write:products": "Create and update products",
                "read:orders": "Read orders",
                "write:orders": "Create and update orders",
                "admin": "Administrative access",
            },
        ).
        WithClientCredentials(
            "https://auth.ecommerce.com/oauth/token",
            map[string]string{
                "webhooks": "Webhook access",
            },
        )

    flows.AuthorizationCode.SetRefreshUrl("https://auth.ecommerce.com/oauth/refresh")

    openapi.SetOAuth2Auth(flows, "OAuth2 authentication")

    // Define advanced endpoints
    endpoints := []*endpoint.EndPoint{
        // GET /products - List products with advanced filtering
        endpoint.New(
            endpoint.GET,
            "/products",
            endpoint.WithTags("products"),
            endpoint.WithSummary("List products"),
            endpoint.WithDescription("Retrieve products with advanced filtering and pagination"),
            endpoint.WithParams(
                // Pagination
                parameter.IntParam("page", parameter.Query,
                    parameter.WithDefault(1),
                    parameter.WithMin(1),
                    parameter.WithExample(1),
                ),
                parameter.IntParam("limit", parameter.Query,
                    parameter.WithDefault(20),
                    parameter.WithMin(1),
                    parameter.WithMax(100),
                ),

                // Filtering
                parameter.StringParam("category", parameter.Query,
                    parameter.WithDescription("Filter by category slug"),
                    parameter.WithExample("electronics"),
                ),
                parameter.FloatParam("min_price", parameter.Query,
                    parameter.WithDescription("Minimum price filter"),
                    parameter.WithMin(0),
                    parameter.WithExample(10.00),
                ),
                parameter.FloatParam("max_price", parameter.Query,
                    parameter.WithDescription("Maximum price filter"),
                    parameter.WithMin(0),
                    parameter.WithExample(1000.00),
                ),
                parameter.BoolParam("in_stock", parameter.Query,
                    parameter.WithDescription("Filter by stock availability"),
                ),

                // Sorting
                parameter.StringParam("sort", parameter.Query,
                    parameter.WithDescription("Sort field"),
                    parameter.WithEnum([]interface{}{"name", "price", "created_at"}),
                    parameter.WithDefault("name"),
                ),
                parameter.StringParam("order", parameter.Query,
                    parameter.WithDescription("Sort order"),
                    parameter.WithEnum([]interface{}{"asc", "desc"}),
                    parameter.WithDefault("asc"),
                ),

                // Advanced array parameter
                parameter.ArrayParam("tags", parameter.Query,
                    parameter.WithDescription("Filter by tags"),
                    parameter.WithStyle("form"),
                    parameter.WithExplode(true),
                    parameter.WithItems(parameter.StringItems()),
                ),
            ),
            endpoint.WithSuccessfulReturns([]response.Response{
                response.New([]Product{}, "200", "Products retrieved successfully"),
            }),
            endpoint.WithSecurity([]map[string][]string{
                {"bearerAuth": {}},
                {"apiKeyAuth": {}},
                {"oauth2": {"read:products"}},
            }),
        ),

        // POST /products - Create product with callbacks
        endpoint.New(
            endpoint.POST,
            "/products",
            endpoint.WithTags("products"),
            endpoint.WithSummary("Create product"),
            endpoint.WithDescription("Create a new product and trigger webhooks"),
            endpoint.WithBody(CreateProductRequest{}),
            endpoint.WithSuccessfulReturns([]response.Response{
                response.NewWithLinks(Product{}, "201", "Product created successfully",
                    map[string]endpoint.Link{
                        "GetProduct": endpoint.NewEnhancedLink().
                            SetOperationId("getProductById").
                            AddParameter("id", "$response.body#/id").
                            SetDescription("Get the created product"),
                        "UpdateProduct": endpoint.NewEnhancedLink().
                            SetOperationId("updateProduct").
                            AddParameter("id", "$response.body#/id").
                            SetDescription("Update the created product"),
                    }),
            }),
            endpoint.WithCallbacks(map[string]endpoint.Callback{
                "productCreated": createProductWebhookCallback(),
            }),
            endpoint.WithSecurity([]map[string][]string{
                {"oauth2": {"write:products"}},
            }),
        ),

        // GET /orders/{id} - Get order with links
        endpoint.New(
            endpoint.GET,
            "/orders/{id}",
            endpoint.WithTags("orders"),
            endpoint.WithSummary("Get order by ID"),
            endpoint.WithDescription("Retrieve an order with linked operations"),
            endpoint.WithParams(
                parameter.IntParam("id", parameter.Path,
                    parameter.WithRequired(),
                    parameter.WithDescription("Order ID"),
                    parameter.WithExample(123),
                ),
            ),
            endpoint.WithSuccessfulReturns([]response.Response{
                response.NewWithLinks(Order{}, "200", "Order found",
                    map[string]endpoint.Link{
                        "CancelOrder": endpoint.NewEnhancedLink().
                            SetOperationId("cancelOrder").
                            AddParameter("id", "$response.body#/id").
                            SetDescription("Cancel this order"),
                        "TrackOrder": endpoint.NewEnhancedLink().
                            SetOperationId("trackOrder").
                            AddParameter("id", "$response.body#/id").
                            SetDescription("Track this order"),
                    }),
            }),
            endpoint.WithSecurity([]map[string][]string{
                {"oauth2": {"read:orders"}},
            }),
        ),

        // POST /webhooks/register - Register webhook endpoint
        endpoint.New(
            endpoint.POST,
            "/webhooks/register",
            endpoint.WithTags("webhooks"),
            endpoint.WithSummary("Register webhook"),
            endpoint.WithDescription("Register a webhook endpoint for receiving events"),
            endpoint.WithBody(WebhookRegistration{}),
            endpoint.WithSecurity([]map[string][]string{
                {"oauth2": {"webhooks"}},
            }),
        ),
    }

    openapi.AddEndpoints(endpoints)
    openapi.ExportOpenAPIDocs("ecommerce-openapi.json")
}

// Helper function to create webhook callback
func createProductWebhookCallback() endpoint.Callback {
    callback := endpoint.NewEnhancedCallback()

    // Webhook callback for product creation
    pathItem := &endpoint.PathItem{
        Post: &endpoint.Operation{
            Description: "Webhook called when a product is created",
            RequestBody: &endpoint.RequestBody{
                Description: "Product creation event",
                Required:    true,
                Content: map[string]endpoint.MediaType{
                    "application/json": {
                        Schema: definition.SchemaFromStruct(WebhookPayload{}),
                    },
                },
            },
            Responses: map[string]endpoint.JsonResponse{
                "200": {
                    Description: "Webhook received successfully",
                },
            },
        },
    }

    callback.AddExpression("{$request.body#/webhookUrl}", pathItem)
    return callback
}

type WebhookRegistration struct {
    URL    string   `json:"url" example:"https://api.client.com/webhooks" format:"uri" description:"Webhook URL"`
    Events []string `json:"events" example:"[\"product.created\", \"order.updated\"]" description:"Events to subscribe to"`
    Secret string   `json:"secret" example:"webhook_secret_key" description:"Secret for webhook signature"`
}
```

## 3. Framework Integration Examples

### Gin Framework Integration

```go
package main

import (
    "net/http"

    "github.com/gin-gonic/gin"
    swagno3 "github.com/go-swagno/swagno/v3"
)

func main() {
    // Create OpenAPI documentation
    openapi := createOpenAPIDoc()

    // Create Gin router
    r := gin.Default()

    // Add OpenAPI documentation endpoint
    r.GET("/openapi.json", func(c *gin.Context) {
        c.Header("Content-Type", "application/json")
        c.String(http.StatusOK, string(openapi.MustToJson()))
    })

    // Add Swagger UI
    r.StaticFile("/docs", "./swagger-ui.html")

    // Add API routes
    api := r.Group("/api/v1")
    {
        api.GET("/users", getUsers)
        api.POST("/users", createUser)
        api.GET("/users/:id", getUserByID)
        api.PUT("/users/:id", updateUser)
        api.DELETE("/users/:id", deleteUser)
    }

    r.Run(":8080")
}

func createOpenAPIDoc() *swagno3.OpenAPI {
    openapi := swagno3.New(swagno3.Config{
        Title:   "Gin API Example",
        Version: "v1.0.0",
    })

    openapi.AddServer("http://localhost:8080/api/v1", "Development server")

    // Add endpoints...

    return openapi
}

// Handler implementations...
func getUsers(c *gin.Context) { /* implementation */ }
func createUser(c *gin.Context) { /* implementation */ }
func getUserByID(c *gin.Context) { /* implementation */ }
func updateUser(c *gin.Context) { /* implementation */ }
func deleteUser(c *gin.Context) { /* implementation */ }
```

### Fiber Framework Integration

```go
package main

import (
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/go-swagno/swagno/v3"
)

func main() {
    app := fiber.New()

    // Enable CORS
    app.Use(cors.New())

    // Create OpenAPI documentation
    openapi := createOpenAPIDoc()

    // Add OpenAPI endpoint
    app.Get("/openapi.json", func(c *fiber.Ctx) error {
        c.Set("Content-Type", "application/json")
        return c.SendString(string(openapi.MustToJson()))
    })

    // API routes
    api := app.Group("/api/v1")
    api.Get("/users", getUsers)
    api.Post("/users", createUser)
    api.Get("/users/:id", getUserByID)
    api.Put("/users/:id", updateUser)
    api.Delete("/users/:id", deleteUser)

    app.Listen(":3000")
}

func createOpenAPIDoc() *swagno3.OpenAPI {
    openapi := swagno3.New(swagno3.Config{
        Title:   "Fiber API Example",
        Version: "v1.0.0",
    })

    openapi.AddServer("http://localhost:3000/api/v1", "Development server")

    // Add endpoints...

    return openapi
}

// Handler implementations...
func getUsers(c *fiber.Ctx) error { /* implementation */ return nil }
func createUser(c *fiber.Ctx) error { /* implementation */ return nil }
func getUserByID(c *fiber.Ctx) error { /* implementation */ return nil }
func updateUser(c *fiber.Ctx) error { /* implementation */ return nil }
func deleteUser(c *fiber.Ctx) error { /* implementation */ return nil }
```

### Net/HTTP Integration

```go
package main

import (
    "net/http"
    "log"

    swagno3 "github.com/go-swagno/swagno/v3"
)

func main() {
    // Create OpenAPI documentation
    openapi := createOpenAPIDoc()

    // Create HTTP mux
    mux := http.NewServeMux()

    // Add OpenAPI endpoint
    mux.HandleFunc("/openapi.json", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.Write(openapi.MustToJson())
    })

    // Add API endpoints
    mux.HandleFunc("/api/v1/users", usersHandler)
    mux.HandleFunc("/api/v1/users/", userHandler)

    log.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", mux))
}

func createOpenAPIDoc() *swagno3.OpenAPI {
    openapi := swagno3.New(swagno3.Config{
        Title:   "Net/HTTP API Example",
        Version: "v1.0.0",
    })

    openapi.AddServer("http://localhost:8080/api/v1", "Development server")

    // Add endpoints...

    return openapi
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        // Get users implementation
    case http.MethodPost:
        // Create user implementation
    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}

func userHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        // Get user by ID implementation
    case http.MethodPut:
        // Update user implementation
    case http.MethodDelete:
        // Delete user implementation
    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}
```

## 4. Complex Schema Examples

### Polymorphic Models with Discriminator

```go
// Base interface
type Animal struct {
    Type string `json:"type" description:"Animal type" discriminator:"type"`
    Name string `json:"name" description:"Animal name"`
    Age  int    `json:"age" description:"Animal age"`
}

// Specific implementations
type Dog struct {
    Animal
    Breed  string `json:"breed" description:"Dog breed"`
    IsGood bool   `json:"is_good" description:"Is good dog" example:"true"`
}

type Cat struct {
    Animal
    Lives      int  `json:"lives" description:"Number of lives remaining" example:"9"`
    IndoorOnly bool `json:"indoor_only" description:"Indoor cat only"`
}

// Usage in endpoint
endpoint.New(
    endpoint.POST,
    "/animals",
    endpoint.WithBody(Animal{}), // Will include discriminator mapping
    endpoint.WithSuccessfulReturns([]response.Response{
        response.New(Animal{}, "201", "Animal created"),
    }),
)
```

### Complex Nested Structures

```go
type Organization struct {
    ID          uint64       `json:"id" description:"Organization ID"`
    Name        string       `json:"name" description:"Organization name"`
    Address     Address      `json:"address" description:"Organization address"`
    Departments []Department `json:"departments" description:"Organization departments"`
    Settings    Settings     `json:"settings" description:"Organization settings"`
}

type Address struct {
    Street   string  `json:"street" description:"Street address"`
    City     string  `json:"city" description:"City"`
    State    string  `json:"state" description:"State/Province"`
    Country  string  `json:"country" description:"Country code" pattern:"^[A-Z]{2}$"`
    ZipCode  string  `json:"zip_code" description:"Postal code"`
    Location *LatLng `json:"location,omitempty" description:"GPS coordinates"`
}

type LatLng struct {
    Latitude  float64 `json:"latitude" description:"Latitude" minimum:"-90" maximum:"90"`
    Longitude float64 `json:"longitude" description:"Longitude" minimum:"-180" maximum:"180"`
}

type Department struct {
    ID        uint64     `json:"id" description:"Department ID"`
    Name      string     `json:"name" description:"Department name"`
    Manager   *Employee  `json:"manager,omitempty" description:"Department manager"`
    Employees []Employee `json:"employees" description:"Department employees"`
}

type Employee struct {
    ID       uint64 `json:"id" description:"Employee ID"`
    Name     string `json:"name" description:"Employee name"`
    Email    string `json:"email" format:"email" description:"Employee email"`
    Position string `json:"position" description:"Job position"`
    Salary   Money  `json:"salary" description:"Employee salary"`
}

type Money struct {
    Amount   float64 `json:"amount" description:"Monetary amount" minimum:"0"`
    Currency string  `json:"currency" description:"Currency code" pattern:"^[A-Z]{3}$" example:"USD"`
}

type Settings struct {
    Theme        string            `json:"theme" enum:"light,dark" description:"UI theme"`
    Notifications bool             `json:"notifications" description:"Enable notifications"`
    Features     map[string]bool   `json:"features" description:"Feature flags"`
    Metadata     map[string]interface{} `json:"metadata,omitempty" description:"Additional metadata"`
}
```

## 5. Advanced Parameter Examples

### Complex Query Parameters

```go
// Array parameters with different styles
parameter.ArrayParam("tags", parameter.Query,
    parameter.WithStyle("form"),      // ?tags=tag1&tags=tag2
    parameter.WithExplode(true),
    parameter.WithDescription("Filter by tags"),
)

parameter.ArrayParam("categories", parameter.Query,
    parameter.WithStyle("spaceDelimited"), // ?categories=cat1 cat2 cat3
    parameter.WithExplode(false),
    parameter.WithDescription("Filter by categories"),
)

parameter.ArrayParam("filters", parameter.Query,
    parameter.WithStyle("pipeDelimited"), // ?filters=filter1|filter2|filter3
    parameter.WithExplode(false),
    parameter.WithDescription("Apply filters"),
)

// Object parameter with deep object style
parameter.ObjectParam("search", parameter.Query,
    parameter.WithStyle("deepObject"), // ?search[name]=john&search[age]=30
    parameter.WithExplode(true),
    parameter.WithDescription("Search parameters"),
)

// Header parameter with examples
parameter.StringParam("Accept-Language", parameter.Header,
    parameter.WithDescription("Preferred language"),
    parameter.WithExamples(map[string]parameter.Example{
        "english": {
            Summary: "English",
            Value:   "en-US",
        },
        "spanish": {
            Summary: "Spanish",
            Value:   "es-ES",
        },
        "french": {
            Summary: "French",
            Value:   "fr-FR",
        },
    }),
)

// Path parameter with pattern
parameter.StringParam("slug", parameter.Path,
    parameter.WithRequired(),
    parameter.WithPattern("^[a-z0-9-]+$"),
    parameter.WithDescription("URL-friendly identifier"),
    parameter.WithExample("my-product-slug"),
)
```

This comprehensive examples guide demonstrates how to use Swagno v3 for creating sophisticated OpenAPI 3.0.3 documentation with all the modern features while maintaining simplicity and readability.
