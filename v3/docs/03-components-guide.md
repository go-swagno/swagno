# Components Guide for OpenAPI 3.0

This section provides detailed explanations of all Swagno v3 components and their OpenAPI 3.0.3 features.

## 1. Definition Component (`components/definition/`)

### Base Schema Definition (`definition.go`)

The definition component generates OpenAPI 3.0 schemas from Go structs.

#### Schema Struct

```go
type Schema struct {
    Type                 string                    `json:"type,omitempty"`
    Format               string                    `json:"format,omitempty"`
    Title                string                    `json:"title,omitempty"`
    Description          string                    `json:"description,omitempty"`
    Default              interface{}               `json:"default,omitempty"`
    Example              interface{}               `json:"example,omitempty"`
    Required             []string                  `json:"required,omitempty"`
    Properties           map[string]Schema         `json:"properties,omitempty"`
    Items                *Schema                   `json:"items,omitempty"`
    AdditionalProperties *Schema                   `json:"additionalProperties,omitempty"`
    Enum                 []interface{}             `json:"enum,omitempty"`

    // OpenAPI 3.0 specific fields
    Nullable             bool                      `json:"nullable,omitempty"`
    ReadOnly             bool                      `json:"readOnly,omitempty"`
    WriteOnly            bool                      `json:"writeOnly,omitempty"`
    Deprecated           bool                      `json:"deprecated,omitempty"`
    Discriminator        *Discriminator            `json:"discriminator,omitempty"`
    XML                  *XML                      `json:"xml,omitempty"`
    ExternalDocs         *ExternalDocs             `json:"externalDocs,omitempty"`

    // Validation
    MultipleOf           *float64                  `json:"multipleOf,omitempty"`
    Maximum              *float64                  `json:"maximum,omitempty"`
    ExclusiveMaximum     *bool                     `json:"exclusiveMaximum,omitempty"`
    Minimum              *float64                  `json:"minimum,omitempty"`
    ExclusiveMinimum     *bool                     `json:"exclusiveMinimum,omitempty"`
    MaxLength            *int64                    `json:"maxLength,omitempty"`
    MinLength            *int64                    `json:"minLength,omitempty"`
    Pattern              string                    `json:"pattern,omitempty"`
    MaxItems             *int64                    `json:"maxItems,omitempty"`
    MinItems             *int64                    `json:"minItems,omitempty"`
    UniqueItems          bool                      `json:"uniqueItems,omitempty"`
    MaxProperties        *int64                    `json:"maxProperties,omitempty"`
    MinProperties        *int64                    `json:"minProperties,omitempty"`

    // Composition
    AllOf                []Schema                  `json:"allOf,omitempty"`
    OneOf                []Schema                  `json:"oneOf,omitempty"`
    AnyOf                []Schema                  `json:"anyOf,omitempty"`
    Not                  *Schema                   `json:"not,omitempty"`
}
```

#### Key OpenAPI 3.0 Features:

##### Nullable Types

```go
type User struct {
    ID    int     `json:"id"`
    Name  string  `json:"name"`
    Email *string `json:"email,omitempty"` // Automatically becomes nullable: true
}
```

##### ReadOnly and WriteOnly

```go
type User struct {
    ID        int       `json:"id" readonly:"true"`         // ReadOnly: true
    Password  string    `json:"password" writeonly:"true"`  // WriteOnly: true
    Name      string    `json:"name"`
}
```

##### Discriminator for Polymorphism

```go
type Animal struct {
    Type string `json:"type" discriminator:"type"`
    Name string `json:"name"`
}

type Dog struct {
    Animal
    Breed string `json:"breed"`
}

type Cat struct {
    Animal
    Lives int `json:"lives"`
}
```

### Enhanced Schema Definition (`enhanced_schema.go`)

Provides additional OpenAPI 3.0.3 features beyond the base schema.

#### EnhancedSchema Features:

```go
type EnhancedSchema struct {
    Schema // Embed existing schema

    // Additional OpenAPI 3.0.3 fields
    Summary string `json:"summary,omitempty"`

    // JSON Schema 2020-12 compatibility
    If   *EnhancedSchema `json:"if,omitempty"`
    Then *EnhancedSchema `json:"then,omitempty"`
    Else *EnhancedSchema `json:"else,omitempty"`

    // Content validation
    ContentMediaType string `json:"contentMediaType,omitempty"`
    ContentEncoding  string `json:"contentEncoding,omitempty"`

    // Extended validation
    Contains              *EnhancedSchema `json:"contains,omitempty"`
    MinContains           *int64          `json:"minContains,omitempty"`
    MaxContains           *int64          `json:"maxContains,omitempty"`
    UnevaluatedItems      *EnhancedSchema `json:"unevaluatedItems,omitempty"`
    UnevaluatedProperties *EnhancedSchema `json:"unevaluatedProperties,omitempty"`
}
```

#### Enhanced Schema Usage:

```go
// Create enhanced schema
schema := definition.NewEnhancedSchema("object").
    SetSummary("User information").
    SetContentMediaType("application/json").
    SetIfThenElse(ifSchema, thenSchema, elseSchema)
```

## 2. Endpoint Component (`components/endpoint/`)

### Core Endpoint Functionality (`endpoints.go`)

#### EndPoint Struct

```go
type EndPoint struct {
    Method              MethodType
    Path                string
    Tags                []string
    Summary             string
    Description         string
    OperationId         string
    Deprecated          bool
    Parameters          []*parameter.Parameter
    Body                interface{}
    SuccessfulReturns   []response.Response
    Errors              []response.Response
    Security            []map[string][]string
    ExternalDocs        *ExternalDocs
    Callbacks           map[string]Callback
    Servers             []OperationServer
}
```

#### Endpoint Creation:

```go
endpoint.New(
    endpoint.GET,
    "/users/{id}",
    endpoint.WithTags("users"),
    endpoint.WithSummary("Get user by ID"),
    endpoint.WithDescription("Retrieve a specific user by their ID"),
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
    endpoint.WithErrors([]response.Response{
        response.New(ErrorResponse{}, "404", "User not found"),
    }),
    endpoint.WithSecurity([]map[string][]string{
        {"bearerAuth": {}},
    }),
)
```

#### Endpoint Options:

##### `WithTags(tags ...string)`

Groups endpoints by tags.

##### `WithSummary(summary string)`

Sets endpoint summary.

##### `WithDescription(description string)`

Sets detailed endpoint description.

##### `WithOperationId(operationId string)`

Sets unique operation ID.

##### `WithDeprecated()`

Marks endpoint as deprecated.

##### `WithParams(params ...*parameter.Parameter)`

Adds parameters to endpoint.

##### `WithBody(body interface{})`

Sets request body schema.

##### `WithSuccessfulReturns(responses []response.Response)`

Sets successful response definitions.

##### `WithErrors(responses []response.Response)`

Sets error response definitions.

##### `WithSecurity(security []map[string][]string)`

Sets endpoint-level security requirements.

##### `WithExternalDocs(url, description string)`

Adds external documentation.

##### `WithCallbacks(callbacks map[string]Callback)`

Adds webhook callbacks (OpenAPI 3.0 feature).

##### `WithServers(servers []OperationServer)`

Sets operation-specific servers.

### Callbacks and Links (`callbacks_links.go`)

OpenAPI 3.0 introduces callbacks and links for describing webhooks and response relationships.

#### Callbacks

```go
type EnhancedCallback struct {
    Expressions map[string]*PathItem `json:"-"`
}
```

**Usage:**

```go
// Create callback for webhook
callback := endpoint.NewEnhancedCallback()
callback.AddExpression("{$request.body#/webhookUrl}", &endpoint.PathItem{
    Post: &endpoint.Operation{
        RequestBody: &endpoint.RequestBody{
            Content: map[string]endpoint.MediaType{
                "application/json": {
                    Schema: &definition.Schema{
                        Type: "object",
                        Properties: map[string]definition.Schema{
                            "event": {Type: "string"},
                            "data":  {Type: "object"},
                        },
                    },
                },
            },
        },
    },
})

// Add to endpoint
endpoint.New(
    endpoint.POST,
    "/subscribe",
    endpoint.WithCallbacks(map[string]endpoint.Callback{
        "webhook": callback,
    }),
)
```

#### Links

```go
type EnhancedLink struct {
    OperationRef string                 `json:"operationRef,omitempty"`
    OperationId  string                 `json:"operationId,omitempty"`
    Parameters   map[string]interface{} `json:"parameters,omitempty"`
    RequestBody  interface{}            `json:"requestBody,omitempty"`
    Description  string                 `json:"description,omitempty"`
    Server       *OperationServer       `json:"server,omitempty"`
}
```

**Usage:**

```go
// Create link to related operation
link := endpoint.NewEnhancedLink().
    SetOperationId("getUserById").
    AddParameter("id", "$response.body#/id").
    SetDescription("Get the created user")

// Add to response
response.New(User{}, "201", "User created").
    WithLinks(map[string]endpoint.Link{
        "GetUser": link,
    })
```

### Operation Fixes (`operation_fixes.go`)

Contains OpenAPI 3.0 specific operation enhancements and fixes.

## 3. Parameter Component (`components/parameter/`)

### Base Parameter Functionality (`parameter.go`)

#### Parameter Interface

```go
type Parameter interface {
    Name() string
    In() Location
    Description() string
    Required() bool
    Type() string
    Format() string
    Example() interface{}
    Default() interface{}
}
```

#### Parameter Creation Functions:

##### String Parameters

```go
// Basic string parameter
parameter.StringParam("name", parameter.Query)

// With options
parameter.StringParam("name", parameter.Query,
    parameter.WithRequired(),
    parameter.WithDescription("User name"),
    parameter.WithExample("John Doe"),
    parameter.WithPattern("^[a-zA-Z ]+$"),
    parameter.WithMinLength(1),
    parameter.WithMaxLength(100),
)
```

##### Integer Parameters

```go
parameter.IntParam("id", parameter.Path,
    parameter.WithRequired(),
    parameter.WithDescription("User ID"),
    parameter.WithExample(123),
    parameter.WithMin(1),
    parameter.WithMax(999999),
)
```

##### Float Parameters

```go
parameter.FloatParam("price", parameter.Query,
    parameter.WithDescription("Price filter"),
    parameter.WithExample(29.99),
    parameter.WithMin(0.0),
    parameter.WithMax(1000.0),
)
```

##### Boolean Parameters

```go
parameter.BoolParam("active", parameter.Query,
    parameter.WithDescription("Filter by active status"),
    parameter.WithExample(true),
)
```

##### Array Parameters

```go
parameter.ArrayParam("tags", parameter.Query,
    parameter.WithDescription("Filter by tags"),
    parameter.WithItems(parameter.StringItems()),
    parameter.WithExample([]string{"tag1", "tag2"}),
)
```

### OpenAPI 3.0 Parameter Features (`openapi3_parameter.go`)

Enhanced parameter features specific to OpenAPI 3.0.

#### Parameter Styles and Explode

```go
type OpenAPI3Parameter struct {
    Parameter // Embed base parameter

    // OpenAPI 3.0 specific fields
    Style         string                     `json:"style,omitempty"`
    Explode       *bool                      `json:"explode,omitempty"`
    AllowReserved bool                       `json:"allowReserved,omitempty"`
    Schema        *definition.Schema         `json:"schema,omitempty"`
    Examples      map[string]ComponentExample `json:"examples,omitempty"`
    Content       map[string]MediaType       `json:"content,omitempty"`
}
```

#### Style Examples:

```go
// Matrix style parameter
parameter.StringParam("coordinates", parameter.Path,
    parameter.WithStyle("matrix"),
    parameter.WithExplode(true),
)

// Label style parameter
parameter.StringParam("tags", parameter.Path,
    parameter.WithStyle("label"),
    parameter.WithExplode(false),
)

// Form style parameter (default for query)
parameter.ArrayParam("filters", parameter.Query,
    parameter.WithStyle("form"),
    parameter.WithExplode(true),
)

// Simple style parameter (default for path and header)
parameter.StringParam("id", parameter.Path,
    parameter.WithStyle("simple"),
)

// SpaceDelimited style for arrays
parameter.ArrayParam("items", parameter.Query,
    parameter.WithStyle("spaceDelimited"),
)

// PipeDelimited style for arrays
parameter.ArrayParam("items", parameter.Query,
    parameter.WithStyle("pipeDelimited"),
)

// DeepObject style for objects
parameter.ObjectParam("filter", parameter.Query,
    parameter.WithStyle("deepObject"),
    parameter.WithExplode(true),
)
```

#### Multiple Examples:

```go
parameter.StringParam("format", parameter.Query,
    parameter.WithExamples(map[string]parameter.Example{
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

## 4. Response Component (`components/http/response/`)

### Response Interface

```go
type Response interface {
    Description() string
    ReturnCode() string
    Schema() interface{}
    Headers() map[string]parameter.Parameter
    Examples() map[string]interface{}
    Links() map[string]endpoint.Link
}
```

### Enhanced Request/Response (`enhanced_request_response.go`)

OpenAPI 3.0 enhanced request and response handling with content types.

#### Response with Content Types

```go
// Simple response
response.New(User{}, "200", "User found")

// Response with multiple content types
response.NewWithContent("200", "User found", map[string]response.MediaType{
    "application/json": {
        Schema: definition.SchemaFromStruct(User{}),
        Examples: map[string]response.Example{
            "user": {
                Summary: "Example user",
                Value: User{
                    ID:   1,
                    Name: "John Doe",
                    Email: "john@example.com",
                },
            },
        },
    },
    "application/xml": {
        Schema: definition.SchemaFromStruct(User{}),
    },
})

// Basic response
response.New(User{}, "200", "User found")

// Response with links
response.New(User{}, "201", "User created").
    WithLinks(map[string]endpoint.Link{
        "GetUser": endpoint.NewEnhancedLink().
            SetOperationId("getUserById").
            AddParameter("id", "$response.body#/id"),
    })
```

#### Request Body Handling

```go
// Simple request body
endpoint.WithBody(CreateUserRequest{})

// Request body with multiple content types
endpoint.WithRequestBody(endpoint.RequestBody{
    Description: "User data",
    Required:    true,
    Content: map[string]endpoint.MediaType{
        "application/json": {
            Schema: definition.SchemaFromStruct(CreateUserRequest{}),
            Examples: map[string]endpoint.Example{
                "user": {
                    Summary: "Example user creation",
                    Value: CreateUserRequest{
                        Name:  "John Doe",
                        Email: "john@example.com",
                    },
                },
            },
        },
        "application/xml": {
            Schema: definition.SchemaFromStruct(CreateUserRequest{}),
        },
        "application/x-www-form-urlencoded": {
            Schema: definition.SchemaFromStruct(CreateUserRequest{}),
        },
    },
})
```

## 5. Security Component (`components/security/`)

### Base Security (`security.go`)

#### Security Requirements

```go
// Single security scheme
endpoint.WithSecurity([]map[string][]string{
    {"bearerAuth": {}},
})

// Multiple security schemes (OR)
endpoint.WithSecurity([]map[string][]string{
    {"bearerAuth": {}},
    {"apiKeyAuth": {}},
})

// Multiple security schemes with scopes (AND)
endpoint.WithSecurity([]map[string][]string{
    {
        "oauth2": {"read", "write"},
        "bearerAuth": {},
    },
})
```

### Security Scopes (`security_scope.go`)

```go
// OAuth2 scopes
scopes := security.Scopes(
    security.Scope("read", "Read access to resources"),
    security.Scope("write", "Write access to resources"),
    security.Scope("admin", "Administrative access"),
)
```

### OpenAPI 3.0 Security Features (`openapi3_security.go`)

Enhanced security features for OpenAPI 3.0.

#### Security Scheme Examples:

```go
// Bearer with JWT format
openapi.SetBearerAuth("JWT", "Bearer authentication using JWT tokens")

// API Key in cookie
openapi.SetApiKeyAuth("sessionId", "cookie", "Session ID authentication")

// OAuth2 with multiple flows
flows := &v3.OAuthFlows{
    AuthorizationCode: &v3.OAuthFlow{
        AuthorizationUrl: "https://example.com/oauth/authorize",
        TokenUrl:        "https://example.com/oauth/token",
        RefreshUrl:      "https://example.com/oauth/refresh",
        Scopes: map[string]string{
            "read":   "Read access",
            "write":  "Write access",
            "delete": "Delete access",
        },
    },
    ClientCredentials: &v3.OAuthFlow{
        TokenUrl: "https://example.com/oauth/token",
        Scopes: map[string]string{
            "api": "API access",
        },
    },
}
openapi.SetOAuth2Auth(flows, "OAuth2 with authorization code and client credentials")

// OpenID Connect
openapi.SetOpenIdConnectAuth(
    "https://example.com/.well-known/openid_configuration",
    "OpenID Connect authentication",
)
```

## 6. Tag Component (`components/tag/`)

### Tag Structure

```go
type Tag struct {
    Name         string                `json:"name"`
    Description  string                `json:"description,omitempty"`
    ExternalDocs *ExternalDocumentation `json:"externalDocs,omitempty"`
}
```

### Tag Creation:

```go
// Simple tag
tag.New("users", "User operations")

// Tag with external documentation
tag.New("users", "User operations",
    tag.WithExternalDocs(
        "https://docs.example.com/users",
        "User API Documentation",
    ),
)
```

## 7. MIME Component (`components/mime/`)

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

### Usage in Endpoints:

```go
endpoint.New(
    endpoint.POST,
    "/users",
    endpoint.WithConsume([]mime.MIME{mime.JSON, mime.XML}),
    endpoint.WithProduce([]mime.MIME{mime.JSON, mime.XML}),
)
```

## 8. Fields Component (`components/fields/`)

### Field Parsing

The fields component handles Go struct field parsing with enhanced support for OpenAPI 3.0 features.

#### Struct Tags Support:

```go
type User struct {
    ID        int     `json:"id" description:"User unique identifier" example:"123"`
    Name      string  `json:"name" description:"User full name" example:"John Doe" minLength:"1" maxLength:"100"`
    Email     *string `json:"email,omitempty" description:"User email address" example:"john@example.com" format:"email"`
    Age       int     `json:"age" description:"User age" example:"30" minimum:"0" maximum:"150"`
    IsActive  bool    `json:"is_active" description:"Whether user is active" example:"true"`
    CreatedAt string  `json:"created_at" description:"Creation timestamp" example:"2023-01-01T00:00:00Z" format:"date-time" readonly:"true"`
    Password  string  `json:"password" description:"User password" writeonly:"true"`
    Metadata  map[string]interface{} `json:"metadata,omitempty" description:"Additional user metadata"`
}
```

#### Supported Tags:

- `json`: JSON field name and omitempty
- `description`: Field description
- `example`: Example value
- `format`: String format (email, date-time, uuid, etc.)
- `minLength`, `maxLength`: String length validation
- `minimum`, `maximum`: Number validation
- `pattern`: Regex pattern for strings
- `readonly`: ReadOnly property
- `writeonly`: WriteOnly property
- `deprecated`: Mark field as deprecated

This comprehensive component system allows Swagno v3 to generate complete OpenAPI 3.0.3 documentation with all modern features while maintaining ease of use.
