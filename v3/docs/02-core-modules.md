# Core Modules and Basic OpenAPI 3.0 Structures

In this section, we will examine Swagno v3's core modules and basic OpenAPI 3.0 structures in detail.

## 1. Main OpenAPI Module (`openapi.go`)

### OpenAPI Struct

```go
type OpenAPI struct {
    OpenAPI      string                       `json:"openapi" default:"3.0.3"`
    Info         Info                         `json:"info"`
    Servers      []Server                     `json:"servers,omitempty"`
    Paths        map[string]endpoint.PathItem `json:"paths"`
    Components   *Components                  `json:"components,omitempty"`
    Tags         []tag.Tag                    `json:"tags,omitempty"`
    Security     []map[string][]string        `json:"security,omitempty"`
    ExternalDocs *ExternalDocs                `json:"externalDocs,omitempty"`
    endpoints    []*endpoint.EndPoint
}
```

**OpenAPI struct** holds all OpenAPI 3.0.3 documentation and generates JSON compliant with OpenAPI 3.0.3 specification.

#### Core Functions:

##### `New(c Config) *OpenAPI`

Creates a new OpenAPI instance.

```go
openapi := v3.New(v3.Config{
    Title: "My API",
    Version: "v1.0.0",
    Description: "API description",
    Contact: &v3.Contact{
        Name:  "API Support",
        Email: "support@example.com",
    },
})
```

##### `AddEndpoints(e []*endpoint.EndPoint)`

Adds multiple endpoints.

```go
openapi.AddEndpoints([]*endpoint.EndPoint{
    endpoint1,
    endpoint2,
})
```

##### `AddEndpoint(e *endpoint.EndPoint)`

Adds a single endpoint.

```go
openapi.AddEndpoint(userEndpoint)
```

##### `AddTags(tags ...tag.Tag)`

Adds OpenAPI tags.

```go
openapi.AddTags(
    tag.New("users", "User operations"),
    tag.New("products", "Product operations"),
)
```

##### `AddServer(url string, description string)`

Adds a server to the OpenAPI specification.

```go
openapi.AddServer("https://api.example.com/v1", "Production server")
openapi.AddServer("https://staging.example.com/v1", "Staging server")
```

### Config Struct

```go
type Config struct {
    Title          string        // Title of the OpenAPI documentation
    Version        string        // API version
    Summary        string        // API summary (OpenAPI 3.0 feature)
    Description    string        // API description
    Servers        []Server      // Server configurations
    License        *License      // License information
    Contact        *Contact      // Contact information
    TermsOfService string        // Terms of service
    ExternalDocs   *ExternalDocs // External documentation
}
```

### Info Struct

```go
type Info struct {
    Title          string   `json:"title"`
    Summary        string   `json:"summary,omitempty"`        // New in OpenAPI 3.0
    Description    string   `json:"description,omitempty"`
    Version        string   `json:"version"`
    TermsOfService string   `json:"termsOfService,omitempty"`
    Contact        *Contact `json:"contact,omitempty"`
    License        *License `json:"license,omitempty"`
}
```

### Server Struct (New in OpenAPI 3.0)

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

**Usage:**

```go
// Simple server
openapi.AddServer("https://api.example.com", "Production server")

// Server with variables (advanced usage)
server := Server{
    URL: "https://{environment}.example.com/{version}",
    Description: "Configurable server",
    Variables: map[string]ServerVariable{
        "environment": {
            Default: "api",
            Enum: []string{"api", "staging", "dev"},
            Description: "Environment name",
        },
        "version": {
            Default: "v1",
            Description: "API version",
        },
    },
}
```

### Components Struct (New in OpenAPI 3.0)

```go
type Components struct {
    Schemas         map[string]definition.Schema       `json:"schemas,omitempty"`
    Responses       map[string]endpoint.JsonResponse   `json:"responses,omitempty"`
    Parameters      map[string]parameter.JsonParameter `json:"parameters,omitempty"`
    Examples        map[string]ComponentExample        `json:"examples,omitempty"`
    RequestBodies   map[string]endpoint.RequestBody    `json:"requestBodies,omitempty"`
    Headers         map[string]ComponentHeader         `json:"headers,omitempty"`
    SecuritySchemes map[string]SecurityScheme          `json:"securitySchemes,omitempty"`
    Links           map[string]endpoint.Link           `json:"links,omitempty"`
    Callbacks       map[string]endpoint.Callback       `json:"callbacks,omitempty"`
}
```

### Contact and License Structs

```go
type Contact struct {
    Name  string `json:"name,omitempty"`
    URL   string `json:"url,omitempty"`   // Changed from Url to URL
    Email string `json:"email,omitempty"`
}

type License struct {
    Name string `json:"name"`            // Required in OpenAPI 3.0
    URL  string `json:"url,omitempty"`   // Changed from URL to url
}
```

## 2. JSON Generation Module (`generate.go`)

This module converts OpenAPI documentation to JSON format with OpenAPI 3.0.3 features.

### Main Functions:

#### `ToJson() ([]byte, error)`

Converts the OpenAPI object to JSON format. Returns error if any issue occurs.

```go
jsonData, err := openapi.ToJson()
if err != nil {
    log.Fatal(err)
}
```

#### `MustToJson() []byte`

Converts the OpenAPI object to JSON format. Panics if any error occurs.

```go
jsonData := openapi.MustToJson()
```

#### `generateOpenAPIJson()`

Internal function. Manages the JSON generation process for OpenAPI object:

1. Validates endpoints and security
2. Generates component schemas
3. Creates path items with OpenAPI 3.0 features
4. Links request bodies and responses
5. Processes callbacks and links

#### `generateOpenAPIDefinition()`

Automatically generates schema definitions in components/schemas from endpoints:

```go
func (o *OpenAPI) generateOpenAPIDefinition() {
    if o.Components == nil {
        o.Components = &Components{
            Schemas: make(map[string]definition.Schema),
        }
    }

    for _, endpoint := range o.endpoints {
        if endpoint.Body != nil {
            o.createSchema(endpoint.Body)
        }
        o.createSchemas(endpoint.SuccessfulReturns())
        o.createSchemas(endpoint.Errors())
    }
}
```

#### `createSchema(t interface{})`

Creates OpenAPI 3.0 schema from given interface and stores it in components/schemas.

#### `processRequestBodies()`

Processes request bodies with proper content type handling (new in OpenAPI 3.0).

## 3. Security Module (`auth.go`)

Manages OpenAPI 3.0 security definitions with enhanced features.

### Security Scheme Types:

#### Basic Authentication

```go
func (o *OpenAPI) SetBasicAuth(description ...string) {
    desc := "Basic Authentication"
    if len(description) > 0 {
        desc = description[0]
    }
    o.ensureComponents()
    o.Components.SecuritySchemes["basicAuth"] = SecurityScheme{
        Type:        "http",
        Scheme:      "basic",
        Description: desc,
    }
}
```

**Usage:**

```go
openapi.SetBasicAuth("Basic auth for API access")
```

#### Bearer Authentication (New in OpenAPI 3.0)

```go
func (o *OpenAPI) SetBearerAuth(bearerFormat string, description ...string) {
    desc := "Bearer Authentication"
    if len(description) > 0 {
        desc = description[0]
    }
    o.ensureComponents()
    o.Components.SecuritySchemes["bearerAuth"] = SecurityScheme{
        Type:         "http",
        Scheme:       "bearer",
        BearerFormat: bearerFormat,
        Description:  desc,
    }
}
```

**Usage:**

```go
openapi.SetBearerAuth("JWT", "Bearer authentication using JWT tokens")
```

#### API Key Authentication

```go
func (o *OpenAPI) SetApiKeyAuth(name string, in string, description ...string)
```

**Parameters:**

- `name`: Name of the API key
- `in`: Location of the API key ("query", "header", "cookie")
- `description`: Description (optional)

**Usage:**

```go
openapi.SetApiKeyAuth("X-API-Key", "header", "API key authentication")
```

#### OAuth2 Authentication (Enhanced in OpenAPI 3.0)

```go
func (o *OpenAPI) SetOAuth2Auth(flows *OAuthFlows, description ...string)
```

**OpenAPI 3.0 OAuth2 Flows:**

```go
type OAuthFlows struct {
    Implicit          *OAuthFlow `json:"implicit,omitempty"`
    Password          *OAuthFlow `json:"password,omitempty"`
    ClientCredentials *OAuthFlow `json:"clientCredentials,omitempty"`
    AuthorizationCode *OAuthFlow `json:"authorizationCode,omitempty"`
}

type OAuthFlow struct {
    AuthorizationUrl string            `json:"authorizationUrl,omitempty"`
    TokenUrl         string            `json:"tokenUrl,omitempty"`
    RefreshUrl       string            `json:"refreshUrl,omitempty"`
    Scopes           map[string]string `json:"scopes"`
}
```

**Usage:**

```go
flows := &v3.OAuthFlows{
    AuthorizationCode: &v3.OAuthFlow{
        AuthorizationUrl: "https://example.com/oauth/authorize",
        TokenUrl:        "https://example.com/oauth/token",
        RefreshUrl:      "https://example.com/oauth/refresh",
        Scopes: map[string]string{
            "read":  "Read access to resources",
            "write": "Write access to resources",
            "admin": "Administrative access",
        },
    },
    ClientCredentials: &v3.OAuthFlow{
        TokenUrl: "https://example.com/oauth/token",
        Scopes: map[string]string{
            "api": "API access",
        },
    },
}
openapi.SetOAuth2Auth(flows, "OAuth2 authentication with multiple flows")
```

#### OpenID Connect Authentication (New in OpenAPI 3.0)

```go
func (o *OpenAPI) SetOpenIdConnectAuth(openIdConnectUrl string, description ...string) {
    desc := "OpenID Connect Authentication"
    if len(description) > 0 {
        desc = description[0]
    }
    o.ensureComponents()
    o.Components.SecuritySchemes["openIdConnect"] = SecurityScheme{
        Type:             "openIdConnect",
        OpenIdConnectUrl: openIdConnectUrl,
        Description:      desc,
    }
}
```

**Usage:**

```go
openapi.SetOpenIdConnectAuth(
    "https://example.com/.well-known/openid_configuration",
    "OpenID Connect authentication",
)
```

### SecurityScheme Struct

```go
type SecurityScheme struct {
    Type             string      `json:"type"`
    Description      string      `json:"description,omitempty"`
    Name             string      `json:"name,omitempty"`
    In               string      `json:"in,omitempty"`
    Scheme           string      `json:"scheme,omitempty"`           // New in OpenAPI 3.0
    BearerFormat     string      `json:"bearerFormat,omitempty"`     // New in OpenAPI 3.0
    Flows            *OAuthFlows `json:"flows,omitempty"`            // Enhanced in OpenAPI 3.0
    OpenIdConnectUrl string      `json:"openIdConnectUrl,omitempty"` // New in OpenAPI 3.0
}
```

## 4. External Documentation Module (`external_docs.go`)

### ExternalDocs Struct

```go
type ExternalDocs struct {
    Description string `json:"description,omitempty"`
    URL         string `json:"url"` // REQUIRED
}
```

### Functions:

#### `NewExternalDocs(url string, description string) *ExternalDocs`

Creates a new ExternalDocs instance.

```go
extDocs := v3.NewExternalDocs(
    "https://docs.example.com/api",
    "Complete API Documentation",
)
```

#### `Validate() error`

Validates the ExternalDocs according to OpenAPI 3.0.3 rules.

**Usage in Config:**

```go
openapi := v3.New(v3.Config{
    Title:       "My API",
    Version:     "v1.0.0",
    ExternalDocs: v3.NewExternalDocs(
        "https://docs.example.com/api",
        "Complete API Documentation",
    ),
})
```

## 5. Utility Functions

### `buildOpenAPI(c Config) *OpenAPI`

Creates new OpenAPI instance from Config. Sets default values:

- Title: "OpenAPI API" (if empty)
- Version: "1.0.0" (if empty)
- OpenAPI: "3.0.3"
- Default server: "/" (if no servers provided)

### `ExportOpenAPIDocs(out_file string) string`

Saves OpenAPI documentation to file and returns JSON string:

```go
jsonStr := openapi.ExportOpenAPIDocs("openapi.json")
```

### `ensureComponents()`

Internal utility to ensure Components struct is initialized:

```go
func (o *OpenAPI) ensureComponents() {
    if o.Components == nil {
        o.Components = &Components{
            Schemas:         make(map[string]definition.Schema),
            SecuritySchemes: make(map[string]SecurityScheme),
        }
    }
}
```

## 6. Enhanced Component Types

### ComponentExample

```go
type ComponentExample struct {
    Summary       string      `json:"summary,omitempty"`
    Description   string      `json:"description,omitempty"`
    Value         interface{} `json:"value,omitempty"`
    ExternalValue string      `json:"externalValue,omitempty"`
}
```

### ComponentHeader

```go
type ComponentHeader struct {
    Description     string                      `json:"description,omitempty"`
    Required        bool                        `json:"required,omitempty"`
    Deprecated      bool                        `json:"deprecated,omitempty"`
    AllowEmptyValue bool                        `json:"allowEmptyValue,omitempty"`
    Style           string                      `json:"style,omitempty"`
    Explode         bool                        `json:"explode,omitempty"`
    AllowReserved   bool                        `json:"allowReserved,omitempty"`
    Schema          *definition.Schema          `json:"schema,omitempty"`
    Example         interface{}                 `json:"example,omitempty"`
    Examples        map[string]ComponentExample `json:"examples,omitempty"`
}
```

## Test Structure

The project has comprehensive test coverage (`generate_test.go`, `openapi_test.go`):

### Test Scenarios:

1. **Basic Functionality Test**: Basic endpoint, parameter and response tests
2. **Deeply Nested Model Test**: Complex, nested model tests
3. **OpenAPI 3.0 Features Test**: Security, servers, request bodies

### Test Data Structure:

```go
testCases := []struct {
    name      string
    endpoints []*endpoint.EndPoint
    file      string
}{
    {
        name: "Basic Functionality Test",
        endpoints: [...],
        file: "testdata/expected_output/bft.json",
    },
    {
        name: "OpenAPI 3.0 Features Test",
        endpoints: [...],
        file: "testdata/expected_output/openapi3.json",
    },
}
```

Thanks to this enhanced modular structure, Swagno v3 offers full OpenAPI 3.0.3 compliance while maintaining extensibility and testability.
