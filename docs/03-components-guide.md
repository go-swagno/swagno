# Detailed Components Guide

In this section, we will examine all components in Swagno's `components/` folder in detail.

## 1. Definition Component (`components/definition/`)

This component automatically generates Swagger model definitions from Go structs.

### Definition Struct

```go
type Definition struct {
    Type       string                          `json:"type"`
    Properties map[string]DefinitionProperties `json:"properties,omitempty"`
    Required   []string                        `json:"required,omitempty"`
}
```

### DefinitionProperties Struct

```go
type DefinitionProperties struct {
    Type    string                     `json:"type,omitempty"`
    Format  string                     `json:"format,omitempty"`
    Ref     string                     `json:"$ref,omitempty"`
    Items   *DefinitionPropertiesItems `json:"items,omitempty"`
    Example interface{}                `json:"example,omitempty"`
    IsRequired bool                    `json:"-"`  // Internal usage
}
```

### DefinitionGenerator

Main struct responsible for generating model definitions:

```go
type DefinitionGenerator struct {
    Definitions map[string]Definition
}
```

#### Main Functions:

##### `NewDefinitionGenerator(definitionMap map[string]Definition) *DefinitionGenerator`

Creates a new generator instance.

##### `CreateDefinition(t interface{})`

Creates definition from given interface. Supported types:

- **Slice**: For array types
- **Struct**: For object types
- **CustomResponse**: For custom response types

**Working Process:**

1. Type determination with reflection
2. Finding element type if slice
3. Parsing struct fields
4. Merging embedded structs
5. Determining required fields

##### `createStructDefinitions(structType reflect.Type) map[string]DefinitionProperties`

Generates properties from struct types. Supported field types:

**Array Types:**

```go
// []string, []int etc.
Tags []string `json:"tags"`

// []Struct types
Sizes []Size `json:"sizes"`

// []*Struct types
SizePtrs []*Size `json:"size_ptrs"`
```

**Struct Types:**

```go
// Special time handling
SaleDate time.Time `json:"sale_date"`
EndDate *time.Time `json:"end_date"`

// Other struct types
Complex ComplexModel `json:"complex"`
```

**Pointer Types:**

```go
// Pointer to struct
CategoryId *uint64 `json:"category_id,omitempty"`
```

**Map Types:**

```go
// Map types generate special definitions
Metadata map[string]interface{} `json:"metadata"`
```

**Interface Types:**

```go
// interface{} types are marked as "Ambiguous Type"
Data interface{} `json:"data"`
```

##### `mergeEmbeddedStructFields(properties map[string]DefinitionProperties)`

Merges embedded struct fields into main struct:

```go
type EmbeddedStruct struct {
    Sizes                    // Embedded
    OtherField int `json:"other_field"`
}
```

Embedded struct fields are moved to upper level.

##### `findRequiredFields(properties map[string]DefinitionProperties) []string`

Determines required fields. Rules:

- Required if `required:"true"` tag exists
- Required if no `omitempty` tag
- Pointer fields are optional by default

##### `isRequired(field reflect.StructField) bool`

Checks if field is required.

### Helper Functions:

#### `timeProperty(field reflect.StructField, required bool) DefinitionProperties`

Generates special property for time.Time fields:

```go
// "string" type, "date-time" format
```

#### `durationProperty(field reflect.StructField, required bool) DefinitionProperties`

Generates special property for time.Duration fields:

```go
// "integer" type
```

#### `refProperty(field reflect.StructField, required bool) DefinitionProperties`

Generates reference property.

#### `defaultProperty(field reflect.StructField) DefinitionProperties`

Generates default property.

## 2. Endpoint Component (`components/endpoint/`)

Defines API endpoints and converts them to JSON format.

### MethodType Enum

```go
type MethodType string

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

### EndPoint Struct

```go
type EndPoint struct {
    method            MethodType
    path              string
    params            []*parameter.Parameter
    tags              []string
    Body              interface{}
    successfulReturns []response.Response
    errors            []response.Response
    description       string
    summary           string
    consume           []mime.MIME
    produce           []mime.MIME
    security          []map[string][]string
}
```

### JsonEndPoint Struct

Endpoint representation for Swagger JSON:

```go
type JsonEndPoint struct {
    Description string                    `json:"description"`
    Consumes    []mime.MIME               `json:"consumes"`
    Produces    []mime.MIME               `json:"produces"`
    Tags        []string                  `json:"tags"`
    Summary     string                    `json:"summary"`
    OperationId string                    `json:"operationId,omitempty"`
    Parameters  []parameter.JsonParameter `json:"parameters"`
    Responses   map[string]JsonResponse   `json:"responses"`
    Security    []map[string][]string     `json:"security,omitempty"`
}
```

### Endpoint Creation:

#### `New(m MethodType, path string, opts ...EndPointOption) *EndPoint`

Creates new endpoint:

```go
endpoint.New(
    endpoint.GET,
    "/users/{id}",
    endpoint.WithTags("users"),
    endpoint.WithParams(parameter.IntParam("id", parameter.Path, parameter.WithRequired())),
    endpoint.WithSuccessfulReturns([]response.Response{response.New(User{}, "200", "OK")}),
    endpoint.WithErrors([]response.Response{response.New(ErrorResponse{}, "404", "Not Found")}),
)
```

### EndPoint Options:

#### `WithTags(tag ...string)`

Adds tags to endpoint:

```go
endpoint.WithTags("users", "public")
```

#### `WithParams(params ...*parameter.Parameter)`

Adds parameters to endpoint:

```go
endpoint.WithParams(
    parameter.IntParam("id", parameter.Path, parameter.WithRequired()),
    parameter.StrParam("filter", parameter.Query),
)
```

#### `WithBody(body interface{})`

Defines request body:

```go
endpoint.WithBody(CreateUserRequest{})
```

#### `WithSuccessfulReturns(ret []response.Response)`

Defines successful responses:

```go
endpoint.WithSuccessfulReturns([]response.Response{
    response.New(User{}, "200", "User found"),
    response.New([]User{}, "200", "Users list"),
})
```

#### `WithErrors(err []response.Response)`

Defines error responses:

```go
endpoint.WithErrors([]response.Response{
    response.New(ErrorResponse{}, "400", "Bad Request"),
    response.New(ErrorResponse{}, "404", "Not Found"),
})
```

#### `WithDescription(des string)` and `WithSummary(s string)`

Adds description and summary.

#### `WithConsume(consume []mime.MIME)` and `WithProduce(produce []mime.MIME)`

Defines Content-Types:

```go
endpoint.WithConsume([]mime.MIME{mime.JSON}),
endpoint.WithProduce([]mime.MIME{mime.JSON, mime.XML}),
```

#### `WithSecurity(security []map[string][]string)`

Defines security requirements.

### Main Functions:

#### `AsJson() JsonEndPoint`

Converts EndPoint to JSON representation.

#### `BodyJsonParameter() *parameter.JsonParameter`

Creates JSON parameter for request body:

```go
// For single object
bodySchema := JsonResponseSchema{
    Ref: "#/definitions/UserCreateRequest",
}

// For array
bodySchema := JsonResponseSchema{
    Type: "array",
    Items: &JsonResponseSchemeItems{
        Ref: "#/definitions/UserCreateRequest",
    },
}
```

## 3. Fields Component (`components/fields/`)

Utility functions for parsing struct fields.

### `ExampleTag(field reflect.StructField) interface{}`

Parses `example` tag:

```go
Name string `json:"name" example:"John Doe"`
Age  int    `json:"age" example:"25"`
```

Tries to convert numeric values to uint64, returns as string if unsuccessful.

### `JsonTag(field reflect.StructField) string`

Parses `json` tag:

```go
Name string `json:"name,omitempty"`  // returns "name"
```

### `IsOmitempty(field reflect.StructField) bool`

Checks `omitempty` tag:

```go
Name string `json:"name,omitempty"`  // returns true
```

### `IsRequired(field reflect.StructField) bool`

Checks `required` tag:

```go
Name string `json:"name" required:"true"`  // returns true
```

### `Type(t string) string`

Converts Go types to Swagger types:

```go
"int", "int64", "uint64" → "integer"
"float32", "float64"     → "number"
"bool"                   → "boolean"
"array", "slice"         → "array"
"interface"              → "interface"
// Others return as is
```

## 4. HTTP Response Component (`components/http/response/`)

Manages HTTP response structures.

### Response Interface

```go
type Response interface {
    Description() string
    ReturnCode() string
}
```

### CustomResponse Struct

Implements Response interface:

```go
type CustomResponse struct {
    Model             any
    returnCodeString  string
    descriptionString string
}
```

### Response Creation:

#### `New(model any, returnCode string, description string) CustomResponse`

Creates new response:

```go
response.New(User{}, "200", "User found successfully")
response.New([]User{}, "200", "Users retrieved")
response.New(ErrorResponse{}, "404", "User not found")
```

### ResponseGenerator

Generates response schemas:

#### `Generate(model any) *parameter.JsonResponseSchema`

Generates JSON response schema from model:

**For Slice:**

```go
// []User →
{
    "type": "array",
    "items": {
        "$ref": "#/definitions/User"
    }
}
```

**For Struct:**

```go
// User →
{
    "$ref": "#/definitions/User"
}
```

**For Map:**

```go
// map[string]interface{} →
{
    "$ref": "#/definitions/map[string]interface {}"
}
```

## 5. MIME Component (`components/mime/`)

Defines MIME types:

```go
type MIME string

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

**Usage:**

```go
endpoint.WithConsume([]mime.MIME{mime.JSON, mime.URLFORM})
endpoint.WithProduce([]mime.MIME{mime.JSON, mime.XML})
```

## 6. Parameter Component (`components/parameter/`)

Defines and manages API parameters.

### Location Enum

```go
type Location string

const (
    Query  Location = "query"
    Header Location = "header"
    Path   Location = "path"
    Form   Location = "formData"
)
```

### ParamType Enum

```go
type ParamType string

const (
    String  ParamType = "string"
    Number  ParamType = "number"
    Integer ParamType = "integer"
    Boolean ParamType = "boolean"
    Array   ParamType = "array"
    File    ParamType = "file"
)
```

### Parameter Struct

```go
type Parameter struct {
    name             string
    typeValue        ParamType
    in               Location
    required         bool
    description      string
    enum             []interface{}
    defaultValue     interface{}
    format           string
    min              int64
    max              int64
    minLen           int64
    maxLen           int64
    pattern          string
    maxItems         int64
    minItems         int64
    uniqueItems      bool
    multipleOf       int64
    collectionFormat CollectionFormat
}
```

### Parameter Creation Functions:

#### `IntParam(name string, l Location, opts ...Option) *Parameter`

Creates integer parameter:

```go
parameter.IntParam("id", parameter.Path,
    parameter.WithRequired(),
    parameter.WithMin(1),
    parameter.WithMax(1000),
)
```

#### `StrParam(name string, l Location, opts ...Option) *Parameter`

Creates string parameter:

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

Creates file upload parameter.

#### `IntEnumParam(name string, l Location, arr []int64, opts ...Option) *Parameter`

Creates integer enum parameter:

```go
parameter.IntEnumParam("status", parameter.Query, []int64{1, 2, 3})
```

#### `StrEnumParam(name string, l Location, arr []string, opts ...Option) *Parameter`

Creates string enum parameter:

```go
parameter.StrEnumParam("type", parameter.Query, []string{"user", "admin", "guest"})
```

### Parameter Options:

- `WithRequired()`: Required parameter
- `WithDescription(string)`: Description
- `WithDefault(interface{})`: Default value
- `WithFormat(string)`: Format (e.g., "date-time")
- `WithMin(int64)`: Minimum value
- `WithMax(int64)`: Maximum value
- `WithMinLen(int64)`: Minimum length
- `WithMaxLen(int64)`: Maximum length
- `WithPattern(string)`: Regex pattern
- `WithMaxItems(int64)`: Max item count for arrays
- `WithMinItems(int64)`: Min item count for arrays
- `WithUniqueItems(bool)`: Should array items be unique
- `WithMultipleOf(int64)`: Number must be multiple of
- `WithCollectionFormat(CollectionFormat)`: Array format ("csv", "ssv", etc.)

## 7. Security Component (`components/security/`)

Manages security scopes and schemas.

### Security Scope Functions:

#### `Scope(name string, description string) swaggerSecurityScope`

Creates a single security scope.

#### `Scopes(scopes ...swaggerSecurityScope) map[string]string`

Converts multiple scopes to map:

```go
scopes := security.Scopes(
    security.Scope("read:users", "Read user data"),
    security.Scope("write:users", "Write user data"),
    security.Scope("admin", "Administrative access"),
)
```

### Security Structs:

#### `BasicAuth`, `ApiKeyAuth`, `OAuth`

Structs for different authentication types.

## 8. Tag Component (`components/tag/`)

Manages Swagger tags.

### Tag Struct

```go
type Tag struct {
    Name         string        `json:"name"`
    Description  string        `json:"description"`
    ExternalDocs *ExternalDocs `json:"externalDocs,omitempty"`
}
```

### Tag Creation:

#### `New(name string, description string, opts ...TagOpts) Tag`

Creates new tag:

```go
tag.New("users", "User management operations",
    tag.WithExternalDocs("User Guide", "https://docs.example.com/users"),
)
```

These components work together to create Swagno's powerful and flexible swagger documentation generation system.
