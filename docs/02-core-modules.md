# Core Modules and Basic Structures

In this section, we will examine Swagno's core modules and basic structures in detail.

## 1. Main Swagger Module (`swagger.go`)

### Swagger Struct

```go
type Swagger struct {
    Swagger             string                                      `json:"swagger" default:"2.0"`
    Info                Info                                        `json:"info"`
    Paths               map[string]map[string]endpoint.JsonEndPoint `json:"paths"`
    BasePath            string                                      `json:"basePath" default:"/"`
    Host                string                                      `json:"host" default:""`
    Definitions         map[string]definition.Definition            `json:"definitions"`
    Schemes             []string                                    `json:"schemes,omitempty"`
    Tags                []tag.Tag                                   `json:"tags,omitempty"`
    SecurityDefinitions map[string]securityDefinition               `json:"securityDefinitions,omitempty"`
    endpoints           []*endpoint.EndPoint
}
```

**Swagger struct** holds all swagger documentation and generates JSON compliant with Swagger 2.0 specification.

#### Core Functions:

##### `New(c Config) *Swagger`

Creates a new Swagger instance.

```go
sw := swagno.New(swagno.Config{
    Title: "My API",
    Version: "v1.0.0",
    Host: "localhost:8080",
    Description: "API description",
})
```

##### `AddEndpoints(e []*endpoint.EndPoint)`

Adds multiple endpoints.

```go
sw.AddEndpoints([]*endpoint.EndPoint{
    endpoint1,
    endpoint2,
})
```

##### `AddEndpoint(e *endpoint.EndPoint)`

Adds a single endpoint.

```go
sw.AddEndpoint(userEndpoint)
```

##### `AddTags(tags ...tag.Tag)`

Adds Swagger tags.

```go
sw.AddTags(
    tag.New("users", "User operations"),
    tag.New("products", "Product operations"),
)
```

### Config Struct

```go
type Config struct {
    Title          string   // Title of the Swagger documentation
    Version        string   // API version
    Description    string   // API description
    Host           string   // Host URL
    Path           string   // Base path
    License        *License // License information
    Contact        *Contact // Contact information
    TermsOfService string   // Terms of service
}
```

### Info Struct

```go
type Info struct {
    Title          string   `json:"title"`
    Description    string   `json:"description"`
    Version        string   `json:"version"`
    TermsOfService string   `json:"termsOfService,omitempty"`
    Contact        *Contact `json:"contact,omitempty"`
    License        *License `json:"license,omitempty"`
}
```

### Contact and License Structs

```go
type Contact struct {
    Name  string `json:"name,omitempty"`
    Url   string `json:"url,omitempty"`
    Email string `json:"email,omitempty"`
}

type License struct {
    Name string `json:"name,omitempty"`
    URL  string `json:"url,omitempty"`
}
```

## 2. JSON Generation Module (`generate.go`)

This module converts swagger documentation to JSON format.

### Main Functions:

#### `ToJson() ([]byte, error)`

Converts the Swagger object to JSON format. Returns error if any issue occurs.

```go
jsonData, err := sw.ToJson()
if err != nil {
    log.Fatal(err)
}
```

#### `MustToJson() []byte`

Converts the Swagger object to JSON format. Panics if any error occurs.

```go
jsonData := sw.MustToJson()
```

#### `generateSwaggerJson()`

Internal function. Manages the JSON generation process for Swagger object:

1. Checks endpoints
2. Generates definitions
3. Creates paths
4. Links responses

#### `generateSwaggerDefinition()`

Automatically generates model definitions from endpoints:

```go
func (s *Swagger) generateSwaggerDefinition() {
    for _, endpoint := range s.endpoints {
        if endpoint.Body != nil {
            s.createDefinition(endpoint.Body)
        }
        s.createDefinitions(endpoint.SuccessfulReturns())
        s.createDefinitions(endpoint.Errors())
    }
}
```

#### `createDefinition(t interface{})`

Creates swagger definition from given interface.

#### `appendResponses()`

Converts response array to JsonResponse map.

## 3. Security Module (`auth.go`)

Manages API security definitions.

### Supported Security Types:

#### Basic Authentication

```go
func (s Swagger) SetBasicAuth(description ...string) {
    desc := "Basic Authentication"
    if len(description) > 0 {
        desc = description[0]
    }
    s.SecurityDefinitions["basicAuth"] = securityDefinition{
        Type:        "basic",
        Description: desc,
    }
}
```

**Usage:**

```go
sw.SetBasicAuth("Basic auth for API access")
```

#### API Key Authentication

```go
func (s Swagger) SetApiKeyAuth(name string, in string, description ...string)
```

**Parameters:**

- `name`: Name of the API key
- `in`: Location of the API key ("query", "header")
- `description`: Description (optional)

**Usage:**

```go
sw.SetApiKeyAuth("X-API-Key", "header", "API key authentication")
```

#### OAuth2 Authentication

```go
func (s Swagger) SetOAuth2Auth(name string, flow string, authorizationUrl string, tokenUrl string, scopes map[string]string, description ...string)
```

**Parameters:**

- `name`: OAuth2 name
- `flow`: OAuth2 flow type ("implicit", "password", "application", "accessCode")
- `authorizationUrl`: Authorization URL
- `tokenUrl`: Token URL
- `scopes`: OAuth2 scopes
- `description`: Description (optional)

**Usage:**

```go
scopes := security.Scopes(
    security.Scope("read:users", "Read user data"),
    security.Scope("write:users", "Write user data"),
)
sw.SetOAuth2Auth("oauth2", "password", "", "http://localhost:8080/oauth/token", scopes)
```

### securityDefinition Struct

```go
type securityDefinition struct {
    Type             string            `json:"type"`
    Description      string            `json:"description,omitempty"`
    Name             string            `json:"name,omitempty"`
    In               string            `json:"in,omitempty"`
    Flow             string            `json:"flow,omitempty"`
    AuthorizationUrl string            `json:"authorizationUrl,omitempty"`
    TokenUrl         string            `json:"tokenUrl,omitempty"`
    Scopes           map[string]string `json:"scopes,omitempty"`
}
```

## 4. Utility Functions

### `buildSwagger(c Config) *Swagger`

Creates new Swagger instance from Config. Sets default values:

- Title: "Swagger API" (if empty)
- Version: "1.0" (if empty)
- Path: "/" (if empty)
- Schemes: ["http", "https"]

### `ExportSwaggerDocs(out_file string) string`

Saves swagger documentation to file and returns JSON string:

```go
jsonStr := sw.ExportSwaggerDocs("swagger.json")
```

## Test Structure

The project has comprehensive test coverage (`generate_test.go`):

### Test Scenarios:

1. **Basic Functionality Test**: Basic endpoint, parameter and response tests
2. **Deeply Nested Model Test**: Complex, nested model tests

### Test Data Structure:

```go
var desc = "Lorem ipsum dolor sit amet..."

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
}
```

Thanks to this modular structure, Swagno offers both an extensible and testable architecture.
