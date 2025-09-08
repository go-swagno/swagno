# API Reference and Function Documentation

This section provides detailed documentation of all public APIs, functions, and methods in Swagno.

## 1. Main Swagger API (`swagno` package)

### 1.1. Swagger Struct

```go
type Swagger struct {
    Swagger             string                                      `json:"swagger"`
    Info                Info                                        `json:"info"`
    Paths               map[string]map[string]endpoint.JsonEndPoint `json:"paths"`
    BasePath            string                                      `json:"basePath"`
    Host                string                                      `json:"host"`
    Definitions         map[string]definition.Definition            `json:"definitions"`
    Schemes             []string                                    `json:"schemes,omitempty"`
    Tags                []tag.Tag                                   `json:"tags,omitempty"`
    SecurityDefinitions map[string]securityDefinition               `json:"securityDefinitions,omitempty"`
    // private field
    endpoints           []*endpoint.EndPoint
}
```

### 1.2. Constructor Functions

#### `New(c Config) *Swagger`

Creates a new Swagger instance.

**Parameters:**

- `c Config`: Swagger configuration object

**Return:**

- `*Swagger`: New swagger instance

**Example:**

```go
sw := swagno.New(swagno.Config{
    Title: "My API",
    Version: "v1.0.0",
    Host: "localhost:8080",
    Description: "My awesome API",
    Contact: &swagno.Contact{
        Name: "Support Team",
        Email: "support@example.com",
    },
    License: &swagno.License{
        Name: "MIT",
        URL: "https://opensource.org/licenses/MIT",
    },
    TermsOfService: "https://example.com/terms",
})
```

### 1.3. Endpoint Management Methods

#### `AddEndpoints(e []*endpoint.EndPoint)`

Adds multiple endpoints at once.

**Parameters:**

- `e []*endpoint.EndPoint`: Array of endpoints

**Example:**

```go
endpoints := []*endpoint.EndPoint{
    endpoint.New(endpoint.GET, "/users", ...),
    endpoint.New(endpoint.POST, "/users", ...),
}
sw.AddEndpoints(endpoints)
```

#### `AddEndpoint(e *endpoint.EndPoint)`

Adds a single endpoint.

**Parameters:**

- `e *endpoint.EndPoint`: Endpoint to add

**Example:**

```go
userEndpoint := endpoint.New(endpoint.GET, "/users/{id}", ...)
sw.AddEndpoint(userEndpoint)
```

#### `AddTags(tags ...tag.Tag)`

Adds Swagger tags.

**Parameters:**

- `tags ...tag.Tag`: Variadic tag parameters

**Example:**

```go
sw.AddTags(
    tag.New("users", "User operations"),
    tag.New("products", "Product operations"),
)
```

### 1.4. JSON Generation Methods

#### `ToJson() ([]byte, error)`

Converts Swagger object to JSON byte array.

**Return:**

- `[]byte`: JSON byte array
- `error`: Error if any

**Example:**

```go
jsonData, err := sw.ToJson()
if err != nil {
    log.Fatal("JSON generation failed:", err)
}
```

#### `MustToJson() []byte`

Converts Swagger object to JSON byte array, panics on error.

**Return:**

- `[]byte`: JSON byte array

**Example:**

```go
jsonData := sw.MustToJson() // Panics on error
```

#### `ExportSwaggerDocs(out_file string) string`

Saves Swagger documentation to file and returns JSON string.

**Parameters:**

- `out_file string`: File path

**Return:**

- `string`: JSON string

**Example:**

```go
jsonStr := sw.ExportSwaggerDocs("swagger.json")
```

### 1.5. Security Methods

#### `SetBasicAuth(description ...string)`

Adds Basic authentication definition.

**Parameters:**

- `description ...string`: Optional description

**Example:**

```go
sw.SetBasicAuth("Basic authentication required")
```

#### `SetApiKeyAuth(name string, in string, description ...string)`

Adds API Key authentication definition.

**Parameters:**

- `name string`: API key name
- `in string`: Location ("query", "header")
- `description ...string`: Optional description

**Example:**

```go
sw.SetApiKeyAuth("X-API-Key", "header", "API key authentication")
```

#### `SetOAuth2Auth(name string, flow string, authorizationUrl string, tokenUrl string, scopes map[string]string, description ...string)`

Adds OAuth2 authentication definition.

**Parameters:**

- `name string`: OAuth2 name
- `flow string`: Flow type ("implicit", "password", "application", "accessCode")
- `authorizationUrl string`: Authorization URL
- `tokenUrl string`: Token URL
- `scopes map[string]string`: Scope map
- `description ...string`: Optional description

**Example:**

```go
scopes := security.Scopes(
    security.Scope("read", "Read access"),
    security.Scope("write", "Write access"),
)
sw.SetOAuth2Auth("oauth2", "password", "", "http://localhost:8080/oauth/token", scopes)
```

### 1.6. Configuration Structs

#### `Config`

```go
type Config struct {
    Title          string
    Version        string
    Description    string
    Host           string
    Path           string
    License        *License
    Contact        *Contact
    TermsOfService string
}
```

#### `Info`

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

#### `Contact`

```go
type Contact struct {
    Name  string `json:"name,omitempty"`
    Url   string `json:"url,omitempty"`
    Email string `json:"email,omitempty"`
}
```

#### `License`

```go
type License struct {
    Name string `json:"name,omitempty"`
    URL  string `json:"url,omitempty"`
}
```

## 2. Endpoint API (`components/endpoint`)

### 2.1. Endpoint Constructor

#### `New(m MethodType, path string, opts ...EndPointOption) *EndPoint`

Creates a new endpoint.

**Parameters:**

- `m MethodType`: HTTP method (GET, POST, PUT, DELETE, PATCH, OPTIONS, HEAD)
- `path string`: Endpoint path
- `opts ...EndPointOption`: Variadic option functions

**Return:**

- `*EndPoint`: New endpoint instance

### 2.2. HTTP Methods

```go
const (
    GET     MethodType = "GET"
    POST    MethodType = "POST"
    PUT     MethodType = "PUT"
    DELETE  MethodType = "DELETE"
    PATCH   MethodType = "PATCH"
    OPTIONS MethodType = "OPTIONS"
    HEAD    MethodType = "HEAD"
)
```

### 2.3. Endpoint Options

#### `WithTags(tag ...string) EndPointOption`

Adds tags to the endpoint.

**Example:**

```go
endpoint.WithTags("users", "public", "v1")
```

#### `WithParams(params ...*parameter.Parameter) EndPointOption`

Adds parameters to the endpoint.

**Example:**

```go
endpoint.WithParams(
    parameter.IntParam("id", parameter.Path, parameter.WithRequired()),
    parameter.StrParam("filter", parameter.Query),
)
```

#### `WithBody(body interface{}) EndPointOption`

Defines request body.

**Example:**

```go
endpoint.WithBody(CreateUserRequest{})
```

#### `WithSuccessfulReturns(ret []response.Response) EndPointOption`

Defines successful responses.

**Example:**

```go
endpoint.WithSuccessfulReturns([]response.Response{
    response.New(User{}, "200", "User found"),
    response.New(UserList{}, "200", "Users retrieved"),
})
```

#### `WithErrors(err []response.Response) EndPointOption`

Defines error responses.

**Example:**

```go
endpoint.WithErrors([]response.Response{
    response.New(ErrorResponse{}, "400", "Bad Request"),
    response.New(ErrorResponse{}, "404", "Not Found"),
})
```

#### `WithDescription(des string) EndPointOption`

Adds endpoint description.

#### `WithSummary(s string) EndPointOption`

Adds endpoint summary.

#### `WithConsume(consume []mime.MIME) EndPointOption`

Defines accepted content-types.

**Example:**

```go
endpoint.WithConsume([]mime.MIME{mime.JSON, mime.URLFORM})
```

#### `WithProduce(produce []mime.MIME) EndPointOption`

Defines produced content-types.

**Example:**

```go
endpoint.WithProduce([]mime.MIME{mime.JSON, mime.XML})
```

#### `WithSecurity(security []map[string][]string) EndPointOption`

Defines security requirements.

**Example:**

```go
endpoint.WithSecurity([]map[string][]string{
    {"oauth2": {"read", "write"}},
    {"apiKey": {}},
})
```

### 2.4. Endpoint Methods

#### `Method() MethodType`

Returns the endpoint's HTTP method.

#### `Path() string`

Returns the endpoint's path.

#### `Params() []*parameter.Parameter`

Returns the endpoint's parameters.

#### `SuccessfulReturns() []response.Response`

Returns successful responses.

#### `Errors() []response.Response`

Returns error responses.

#### `AsJson() JsonEndPoint`

Converts endpoint to JSON representation.

#### `BodyJsonParameter() *parameter.JsonParameter`

Creates JSON parameter for request body.

## 3. Parameter API (`components/parameter`)

### 3.1. Parameter Constructors

#### `IntParam(name string, l Location, opts ...Option) *Parameter`

Creates integer parameter.

**Example:**

```go
parameter.IntParam("id", parameter.Path,
    parameter.WithRequired(),
    parameter.WithMin(1),
    parameter.WithMax(1000),
)
```

#### `StrParam(name string, l Location, opts ...Option) *Parameter`

Creates string parameter.

**Example:**

```go
parameter.StrParam("name", parameter.Query,
    parameter.WithMinLen(2),
    parameter.WithMaxLen(50),
    parameter.WithPattern("^[a-zA-Z]+$"),
)
```

#### `BoolParam(name string, l Location, opts ...Option) *Parameter`

Creates boolean parameter.

#### `FileParam(name string, opts ...Option) *Parameter`

Creates file parameter.

#### `IntEnumParam(name string, l Location, arr []int64, opts ...Option) *Parameter`

Creates integer enum parameter.

**Example:**

```go
parameter.IntEnumParam("status", parameter.Query, []int64{1, 2, 3})
```

#### `StrEnumParam(name string, l Location, arr []string, opts ...Option) *Parameter`

Creates string enum parameter.

**Example:**

```go
parameter.StrEnumParam("type", parameter.Query, []string{"user", "admin"})
```

#### `IntArrParam(name string, l Location, arr []int64, opts ...Option) *Parameter`

Creates integer array parameter.

#### `StrArrParam(name string, l Location, arr []string, opts ...Option) *Parameter`

Creates string array parameter.

### 3.2. Parameter Locations

```go
const (
    Query  Location = "query"
    Header Location = "header"
    Path   Location = "path"
    Form   Location = "formData"
)
```

### 3.3. Parameter Types

```go
const (
    String  ParamType = "string"
    Number  ParamType = "number"
    Integer ParamType = "integer"
    Boolean ParamType = "boolean"
    Array   ParamType = "array"
    File    ParamType = "file"
)
```

### 3.4. Parameter Options

#### `WithRequired() Option`

Makes parameter required.

#### `WithDescription(description string) Option`

Adds parameter description.

#### `WithDefault(defaultValue interface{}) Option`

Sets default value.

#### `WithFormat(format string) Option`

Sets format (e.g., "date-time").

#### `WithMin(min int64) Option`

Sets minimum value.

#### `WithMax(max int64) Option`

Sets maximum value.

#### `WithMinLen(minLen int64) Option`

Sets minimum length.

#### `WithMaxLen(maxLen int64) Option`

Sets maximum length.

#### `WithPattern(pattern string) Option`

Sets regex pattern.

#### `WithMaxItems(maxItems int64) Option`

Sets array maximum item count.

#### `WithMinItems(minItems int64) Option`

Sets array minimum item count.

#### `WithUniqueItems(uniqueItems bool) Option`

Sets whether array items should be unique.

#### `WithMultipleOf(multipleOf int64) Option`

Forces number to be multiple of specified value.

#### `WithCollectionFormat(c CollectionFormat) Option`

Sets array serialization format.

### 3.5. Collection Formats

```go
const (
    CSV   CollectionFormat = "csv"    // comma separated: a,b,c
    SSV   CollectionFormat = "ssv"    // space separated: a b c
    TSV   CollectionFormat = "tsv"    // tab separated: a\tb\tc
    Pipes CollectionFormat = "pipes"  // pipe separated: a|b|c
    Multi CollectionFormat = "multi"  // multiple parameters: a=1&a=2
)
```

## 4. Response API (`components/http/response`)

### 4.1. Response Interface

```go
type Response interface {
    Description() string
    ReturnCode() string
}
```

### 4.2. Response Constructor

#### `New(model any, returnCode string, description string) CustomResponse`

Creates a new response.

**Parameters:**

- `model any`: Response model
- `returnCode string`: HTTP status code
- `description string`: Response description

**Return:**

- `CustomResponse`: Response instance

**Example:**

```go
response.New(User{}, "200", "User retrieved successfully")
response.New([]User{}, "200", "Users list")
response.New(ErrorResponse{}, "404", "User not found")
```

### 4.3. ResponseGenerator

#### `NewResponseGenerator() *ResponseGenerator`

Creates a new response generator.

#### `Generate(model any) *parameter.JsonResponseSchema`

Generates JSON response schema from model.

## 5. MIME Types API (`components/mime`)

### 5.1. MIME Constants

```go
const (
    JSON       MIME = "application/json"
    XML        MIME = "application/xml"
    URLFORM    MIME = "application/x-www-form-urlencoded"
    MULTIFORM  MIME = "multipart/form-data"
    PLAINTEXT  MIME = "text/plain"
    HTML       MIME = "text/html"
    JAVASCRIPT MIME = "application/javascript"
)
```

## 6. Security API (`components/security`)

### 6.1. Scope Functions

#### `Scope(name string, description string) swaggerSecurityScope`

Creates a single security scope.

#### `Scopes(scopes ...swaggerSecurityScope) map[string]string`

Converts multiple scopes to map.

**Example:**

```go
scopes := security.Scopes(
    security.Scope("read:users", "Read user data"),
    security.Scope("write:users", "Write user data"),
    security.Scope("admin", "Administrative access"),
)
```

## 7. Tag API (`components/tag`)

### 7.1. Tag Constructor

#### `New(name string, description string, opts ...TagOpts) Tag`

Creates a new tag.

**Parameters:**

- `name string`: Tag name
- `description string`: Tag description
- `opts ...TagOpts`: Variadic option functions

**Return:**

- `Tag`: Tag instance

**Example:**

```go
tag.New("users", "User management operations",
    tag.WithExternalDocs("User Guide", "https://docs.example.com/users"),
)
```

### 7.2. Tag Options

#### `WithExternalDocs(name string, description string) TagOpts`

Adds external documentation.

## 8. Fields Utility API (`components/fields`)

### 8.1. Struct Tag Parsers

#### `ExampleTag(field reflect.StructField) interface{}`

Parses `example` tag.

#### `JsonTag(field reflect.StructField) string`

Parses `json` tag.

#### `IsOmitempty(field reflect.StructField) bool`

Checks for `omitempty` tag.

#### `IsRequired(field reflect.StructField) bool`

Checks for `required` tag.

#### `Type(t string) string`

Converts Go types to Swagger types.

## 9. Error Handling

Error cases in Swagno:

### 9.1. JSON Generation Errors

- `ToJson()` function returns error
- `MustToJson()` function panics

### 9.2. Parameter Validation

- Invalid parameter combinations are logged
- Default values used on reflection errors

### 9.3. Definition Generation

- Interface{} types marked as "Ambiguous Type"
- Default handling for unsupported types

## 10. Performance Considerations

### 10.1. Memory Usage

- Definitions are cached
- Recursive struct detection prevents infinite loops

### 10.2. Reflection Overhead

- Struct parsing done at runtime, not build time
- No field caching (parsed every time)

This API reference covers all public interfaces of Swagno with detailed usage examples.
