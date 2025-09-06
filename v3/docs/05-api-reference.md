# API Reference

This section provides a comprehensive reference for all Swagno v3 APIs and functions.

## 1. Core OpenAPI Functions

### `swagno3.New(config Config) *OpenAPI`

Creates a new OpenAPI 3.0.3 instance.

**Parameters:**

- `config`: Configuration for the OpenAPI document

**Returns:**

- `*OpenAPI`: New OpenAPI instance

**Example:**

```go
openapi := swagno3.New(swagno3.Config{
    Title:       "My API",
    Version:     "v1.0.0",
    Description: "API description",
    Contact: &swagno3.Contact{
        Name:  "Support",
        Email: "support@example.com",
    },
})
```

### `(o *OpenAPI) AddEndpoints(endpoints []*endpoint.EndPoint)`

Adds multiple endpoints to the OpenAPI specification.

**Parameters:**

- `endpoints`: Array of endpoint definitions

**Example:**

```go
openapi.AddEndpoints([]*endpoint.EndPoint{
    endpoint1,
    endpoint2,
    endpoint3,
})
```

### `(o *OpenAPI) AddEndpoint(endpoint *endpoint.EndPoint)`

Adds a single endpoint to the OpenAPI specification.

**Parameters:**

- `endpoint`: Single endpoint definition

### `(o *OpenAPI) AddServer(url, description string)`

Adds a server to the OpenAPI specification.

**Parameters:**

- `url`: Server URL
- `description`: Server description

**Example:**

```go
openapi.AddServer("https://api.example.com/v1", "Production server")
openapi.AddServer("https://staging.example.com/v1", "Staging server")
```

### `(o *OpenAPI) AddTags(tags ...tag.Tag)`

Adds tags to the OpenAPI specification.

**Parameters:**

- `tags`: Variable number of tag definitions

**Example:**

```go
openapi.AddTags(
    tag.New("users", "User operations"),
    tag.New("products", "Product operations"),
)
```

### `(o *OpenAPI) ToJson() ([]byte, error)`

Converts the OpenAPI specification to JSON format.

**Returns:**

- `[]byte`: JSON representation
- `error`: Error if conversion fails

### `(o *OpenAPI) MustToJson() []byte`

Converts the OpenAPI specification to JSON format. Panics on error.

**Returns:**

- `[]byte`: JSON representation

### `(o *OpenAPI) ExportOpenAPIDocs(filename string) string`

Exports the OpenAPI specification to a file and returns the JSON string.

**Parameters:**

- `filename`: Output file name

**Returns:**

- `string`: JSON representation

## 2. Configuration Types

### `Config`

Configuration structure for creating OpenAPI instances.

```go
type Config struct {
    Title          string        // API title (required)
    Version        string        // API version (required)
    Summary        string        // API summary
    Description    string        // API description
    Servers        []Server      // Server configurations
    License        *License      // License information
    Contact        *Contact      // Contact information
    TermsOfService string        // Terms of service URL
    ExternalDocs   *ExternalDocs // External documentation
}
```

### `Contact`

Contact information structure.

```go
type Contact struct {
    Name  string `json:"name,omitempty"`
    URL   string `json:"url,omitempty"`
    Email string `json:"email,omitempty"`
}
```

### `License`

License information structure.

```go
type License struct {
    Name string `json:"name"`           // Required
    URL  string `json:"url,omitempty"`
}
```

### `Server`

Server configuration structure.

```go
type Server struct {
    URL         string                    `json:"url"`
    Description string                    `json:"description,omitempty"`
    Variables   map[string]ServerVariable `json:"variables,omitempty"`
}

type ServerVariable struct {
    Enum        []string `json:"enum,omitempty"`
    Default     string   `json:"default"`
    Description string   `json:"description,omitempty"`
}
```

## 3. Security Functions

### `(o *OpenAPI) SetBasicAuth(description ...string)`

Sets up Basic authentication.

**Parameters:**

- `description`: Optional description

**Example:**

```go
openapi.SetBasicAuth("Basic authentication required")
```

### `(o *OpenAPI) SetBearerAuth(bearerFormat string, description ...string)`

Sets up Bearer token authentication.

**Parameters:**

- `bearerFormat`: Bearer token format (e.g., "JWT")
- `description`: Optional description

**Example:**

```go
openapi.SetBearerAuth("JWT", "Bearer authentication using JWT tokens")
```

### `(o *OpenAPI) SetApiKeyAuth(name string, in security.SecuritySchemeIn, description ...string)`

Sets up API key authentication.

**Parameters:**

- `name`: API key name
- `in`: Location (can use constants `security.Query`, `security.Header`, `security.Cookie` or string literals "query", "header", "cookie")
- `description`: Optional description

**Example:**

```go
import "github.com/go-swagno/swagno/v3/components/security"

// Using constants (recommended)
openapi.SetApiKeyAuth("X-API-Key", security.Header, "API key authentication")

// Using string literals (also works due to type conversion)
openapi.SetApiKeyAuth("X-API-Key", "header", "API key authentication")
```

### `(o *OpenAPI) SetOAuth2Auth(flows *security.OAuthFlows, description ...string)`

Sets up OAuth2 authentication.

**Parameters:**

- `flows`: OAuth2 flows configuration from security package
- `description`: Optional description

**Example:**

```go
import "github.com/go-swagno/swagno/v3/components/security"

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

### `(o *OpenAPI) SetOpenIdConnectAuth(url string, description ...string)`

Sets up OpenID Connect authentication.

**Parameters:**

- `url`: OpenID Connect discovery URL
- `description`: Optional description

**Example:**

```go
openapi.SetOpenIdConnectAuth(
    "https://example.com/.well-known/openid_configuration",
    "OpenID Connect authentication",
)
```

## 4. Endpoint Functions

### `endpoint.New(method MethodType, path string, options ...EndPointOption) *EndPoint`

Creates a new endpoint definition.

**Parameters:**

- `method`: HTTP method (GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS)
- `path`: Endpoint path
- `options`: Variable number of endpoint options

**Returns:**

- `*EndPoint`: New endpoint definition

**Example:**

```go
ep := endpoint.New(
    endpoint.GET,
    "/users/{id}",
    endpoint.WithTags("users"),
    endpoint.WithSummary("Get user by ID"),
    endpoint.WithParams(
        parameter.IntParam("id", parameter.Path, parameter.WithRequired()),
    ),
)
```

### Endpoint Options

#### `endpoint.WithTags(tags ...string)`

Sets endpoint tags.

#### `endpoint.WithSummary(summary string)`

Sets endpoint summary.

#### `endpoint.WithDescription(description string)`

Sets endpoint description.

#### `endpoint.WithOperationId(operationId string)`

Sets endpoint operation ID.

#### `endpoint.WithDeprecated()`

Marks endpoint as deprecated.

#### `endpoint.WithParams(params ...*parameter.Parameter)`

Adds parameters to endpoint.

#### `endpoint.WithBody(body interface{})`

Sets request body schema.

#### `endpoint.WithSuccessfulReturns(responses []response.Response)`

Sets successful response definitions.

#### `endpoint.WithErrors(responses []response.Response)`

Sets error response definitions.

#### `endpoint.WithSecurity(security []map[string][]string)`

Sets endpoint security requirements.

#### `endpoint.WithExternalDocs(url, description string)`

Sets external documentation.

#### `endpoint.WithCallbacks(callbacks map[string]Callback)`

Sets webhook callbacks.

#### `endpoint.WithServers(servers []OperationServer)`

Sets operation-specific servers.

## 5. Parameter Functions

### String Parameters

#### `parameter.StringParam(name string, location Location, options ...ParameterOption) Parameter`

Creates a string parameter.

**Parameters:**

- `name`: Parameter name
- `location`: Parameter location (Path, Query, Header, Cookie)
- `options`: Parameter options

**Example:**

```go
param := parameter.StringParam("name", parameter.Query,
    parameter.WithRequired(),
    parameter.WithDescription("User name"),
    parameter.WithExample("John Doe"),
    parameter.WithPattern("^[a-zA-Z ]+$"),
    parameter.WithMinLength(1),
    parameter.WithMaxLength(100),
)
```

### Integer Parameters

#### `parameter.IntParam(name string, location Location, options ...ParameterOption) Parameter`

Creates an integer parameter.

**Example:**

```go
param := parameter.IntParam("id", parameter.Path,
    parameter.WithRequired(),
    parameter.WithDescription("User ID"),
    parameter.WithExample(123),
    parameter.WithMin(1),
    parameter.WithMax(999999),
)
```

### Float Parameters

#### `parameter.FloatParam(name string, location Location, options ...ParameterOption) Parameter`

Creates a float parameter.

### Boolean Parameters

#### `parameter.BoolParam(name string, location Location, options ...ParameterOption) Parameter`

Creates a boolean parameter.

### Array Parameters

#### `parameter.ArrayParam(name string, location Location, options ...ParameterOption) Parameter`

Creates an array parameter.

**Example:**

```go
param := parameter.ArrayParam("tags", parameter.Query,
    parameter.WithDescription("Filter by tags"),
    parameter.WithItems(parameter.StringItems()),
    parameter.WithStyle("form"),
    parameter.WithExplode(true),
)
```

### Parameter Options

#### `parameter.WithRequired()`

Marks parameter as required.

#### `parameter.WithDescription(description string)`

Sets parameter description.

#### `parameter.WithExample(example interface{})`

Sets parameter example.

#### `parameter.WithDefault(defaultValue interface{})`

Sets parameter default value.

#### `parameter.WithDeprecated()`

Marks parameter as deprecated.

#### `parameter.WithPattern(pattern string)`

Sets string pattern validation.

#### `parameter.WithMinLength(min int64)`

Sets minimum string length.

#### `parameter.WithMaxLength(max int64)`

Sets maximum string length.

#### `parameter.WithMin(min float64)`

Sets minimum numeric value.

#### `parameter.WithMax(max float64)`

Sets maximum numeric value.

#### `parameter.WithEnum(values []interface{})`

Sets enumeration values.

#### `parameter.WithStyle(style string)`

Sets parameter style (OpenAPI 3.0 feature).

**Styles:**

- `"simple"`: Default for path and header parameters
- `"form"`: Default for query and cookie parameters
- `"matrix"`: Matrix style for path parameters
- `"label"`: Label style for path parameters
- `"spaceDelimited"`: Space-separated values for arrays
- `"pipeDelimited"`: Pipe-separated values for arrays
- `"deepObject"`: Deep object style for objects

#### `parameter.WithExplode(explode bool)`

Sets parameter explode behavior.

#### `parameter.WithAllowReserved()`

Allows reserved characters in parameter value.

#### `parameter.WithExamples(examples map[string]parameter.Example)`

Sets multiple parameter examples.

## 6. Response Functions

### `response.New(schema interface{}, code, description string) Response`

Creates a new response definition.

**Parameters:**

- `schema`: Response schema (Go struct or type)
- `code`: HTTP status code
- `description`: Response description

**Example:**

```go
resp := response.New(User{}, "200", "User found successfully")
```

### `response.NewWithLinks(schema interface{}, code, description string, links map[string]endpoint.Link) Response`

Creates a response with links.

**Example:**

```go
resp := response.NewWithLinks(User{}, "201", "User created",
    map[string]endpoint.Link{
        "GetUser": endpoint.NewEnhancedLink().
            SetOperationId("getUserById").
            AddParameter("id", "$response.body#/id"),
    },
)
```

### `response.NewWithContent(code, description string, content map[string]response.MediaType) Response`

Creates a response with custom content types.

**Example:**

```go
resp := response.NewWithContent("200", "User found", map[string]response.MediaType{
    "application/json": {
        Schema: definition.SchemaFromStruct(User{}),
        Examples: map[string]response.Example{
            "user": {
                Summary: "Example user",
                Value: User{ID: 1, Name: "John"},
            },
        },
    },
    "application/xml": {
        Schema: definition.SchemaFromStruct(User{}),
    },
})
```

## 7. Tag Functions

### `tag.New(name, description string, options ...TagOption) Tag`

Creates a new tag.

**Parameters:**

- `name`: Tag name
- `description`: Tag description
- `options`: Tag options

**Example:**

```go
userTag := tag.New("users", "User management operations",
    tag.WithExternalDocs("https://docs.example.com/users", "User docs"),
)
```

### Tag Options

#### `tag.WithExternalDocs(url, description string)`

Adds external documentation to tag.

## 8. External Documentation

### `swagno3.NewExternalDocs(url, description string) *ExternalDocs`

Creates external documentation reference.

**Parameters:**

- `url`: Documentation URL (required)
- `description`: Documentation description

**Example:**

```go
extDocs := swagno3.NewExternalDocs(
    "https://docs.example.com/api",
    "Complete API Documentation",
)
```

## 9. Callback Functions (OpenAPI 3.0)

### `endpoint.NewEnhancedCallback() *EnhancedCallback`

Creates a new webhook callback.

**Example:**

```go
callback := endpoint.NewEnhancedCallback()
callback.AddExpression("{$request.body#/webhookUrl}", &endpoint.PathItem{
    Post: &endpoint.Operation{
        Description: "Webhook endpoint",
        RequestBody: &endpoint.RequestBody{
            Content: map[string]endpoint.MediaType{
                "application/json": {
                    Schema: webhookSchema,
                },
            },
        },
    },
})
```

## 10. Link Functions (OpenAPI 3.0)

### `endpoint.NewEnhancedLink() *EnhancedLink`

Creates a new response link.

**Example:**

```go
link := endpoint.NewEnhancedLink().
    SetOperationId("getUserById").
    AddParameter("id", "$response.body#/id").
    SetDescription("Get the created user")
```

### Link Methods

#### `(l *EnhancedLink) SetOperationRef(operationRef string) *EnhancedLink`

Sets operation reference.

#### `(l *EnhancedLink) SetOperationId(operationId string) *EnhancedLink`

Sets operation ID.

#### `(l *EnhancedLink) AddParameter(name string, expression interface{}) *EnhancedLink`

Adds parameter mapping.

#### `(l *EnhancedLink) SetRequestBody(body interface{}) *EnhancedLink`

Sets request body for link.

#### `(l *EnhancedLink) SetDescription(description string) *EnhancedLink`

Sets link description.

## 11. Schema Functions

### `definition.SchemaFromStruct(model interface{}) definition.Schema`

Creates schema from Go struct.

**Example:**

```go
schema := definition.SchemaFromStruct(User{})
```

### `definition.NewEnhancedSchema(schemaType string) *EnhancedSchema`

Creates enhanced schema with OpenAPI 3.0 features.

**Example:**

```go
schema := definition.NewEnhancedSchema("object").
    SetSummary("User information").
    SetContentMediaType("application/json")
```

## 12. Validation Functions

### `(ed *ExternalDocs) Validate() error`

Validates external documentation.

### `(ec *EnhancedCallback) Validate() error`

Validates callback definition.

### `(el *EnhancedLink) Validate() error`

Validates link definition.

### `ValidateRuntimeExpression(expr string) error`

Validates OpenAPI 3.0 runtime expressions.

**Example:**

```go
err := endpoint.ValidateRuntimeExpression("$request.body#/webhookUrl")
if err != nil {
    log.Fatal("Invalid runtime expression:", err)
}
```

## 13. Constants

### HTTP Methods

```go
const (
    GET     MethodType = "GET"
    POST    MethodType = "POST"
    PUT     MethodType = "PUT"
    DELETE  MethodType = "DELETE"
    PATCH   MethodType = "PATCH"
    HEAD    MethodType = "HEAD"
    OPTIONS MethodType = "OPTIONS"
    TRACE   MethodType = "TRACE"
)
```

### Parameter Locations

```go
const (
    Path   Location = "path"
    Query  Location = "query"
    Header Location = "header"
    Cookie Location = "cookie"
)
```

### MIME Types

```go
const (
    JSON                = "application/json"
    XML                 = "application/xml"
    FormURLEncoded      = "application/x-www-form-urlencoded"
    MultipartFormData   = "multipart/form-data"
    TextPlain           = "text/plain"
    TextHTML            = "text/html"
    OctetStream         = "application/octet-stream"
)
```

This comprehensive API reference covers all the major functions and types available in Swagno v3 for creating OpenAPI 3.0.3 documentation.
