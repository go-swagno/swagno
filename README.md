# swagno: _no annotations, no files, no command_

<img src="https://user-images.githubusercontent.com/1047345/188009539-ea9d0106-979d-4f98-83a3-0d7df6969c9f.png" alt="Swagno" align="right" width="200"/>

Swagno is an approach to create Swagger Documentation 2.0 without any **annotation**, **exported file** or any **command** to run.
You can declare your documentation details in code and get a json string to serve with a handler.

## About the Project

This project inspired by [Swaggo](https://github.com/swaggo/swag). Swaggo, uses annotations, exports files and needs to run by command. If you don't like this way, [Swag**no**](https://github.com/go-swagno/swagno) appears as a good alternative.

This project was then forked from [go-swagno/swagno](https://github.com/go-swagno/swagno) with the goal of being more idiomatic and  user friendly. While trying to achieve these goals a lot of breaking changes were made and hence why this repo is separate as opposed to being merged into the original.

## Improvements

Some improvements from the orignal [Swagno](https://github.com/go-swagno/swagno) are shown below:

- More idiomatic and easier to read API calls
- Added more type safety
- Uses functional option paremeters to allow for more flexible and robust Endpoint creation
- Allows for more than one response and error model types to be used when making endpoints
- Fixed some core structural and semantic bug fixes when rendering swagger/openapi page

Before:

```go
endpoints := []Endpoint{
  // constructor function doesn't allow for options, only allows for one error response with not configuring of things like response code and description, and no type safety
  EndPoint(GET, "/product", "product", Params(), nil, []models.Product{}, models.ErrorResponse{}, "Get all products", nil), 
  EndPoint(GET, "/product/{id}", "product", Params(IntParam("id", true, "")), nil, models.Product{}, models.ErrorResponse{}, "", nil),
  EndPoint(POST, "/product", "product", Params(), models.ProductPost{}, models.Product{}, models.ErrorResponse{}, "", nil),
  // no return
  EndPoint(POST, "/product-no-return", "product", Params(), nil, nil, models.ErrorResponse{}, "", nil),
  // no error
  EndPoint(POST, "/product-no-error", "product", Params(), nil, nil, nil, "", nil),
}
```

After:

```go
// New constructor is idiomatic, allows for the client to use as many or as little options as they like, and allows for multiple responses and errors to be modelled for endpoints instead of restricting to just one.
endpoints := []*endpoint.EndPoint{
  endpoint.New(
   endpoint.WithMethod(endpoint.GET),
   endpoint.WithPath("/product/page"),
   endpoint.WithTags("product"),
   endpoint.WithSuccessfulReturns([]response.Info{models.UnsuccessfulResponse{}}),
   endpoint.WithErrors([]response.Info{models.EmptySuccessfulResponse{}}),
   endpoint.WithDescription(desc),
   endpoint.WithProduce([]string{"application/json", "application/xml"}),
   endpoint.WithConsume([]string{"application/json"}),
  ),
  endpoint.New(
   endpoint.WithMethod(endpoint.GET),
   endpoint.WithPath("/product"),
   endpoint.WithTags("product"),
   endpoint.WithParams(parameter.IntParam("id", parameter.WithRequired())),
   endpoint.WithSuccessfulReturns([]response.Info{models.SuccessfulResponse{}}),
   endpoint.WithErrors([]response.Info{models.UnsuccessfulResponse{}}),
  ),
  endpoint.New(
   endpoint.WithMethod(endpoint.GET),
   endpoint.WithPath("/product/{id}/detail"),
   endpoint.WithTags("product"),
   endpoint.WithParams(parameter.IntParam("id", parameter.WithRequired())),
   endpoint.WithSuccessfulReturns([]response.Info{models.SuccessfulResponse{}}),
   endpoint.WithErrors([]response.Info{models.UnsuccessfulResponse{}}),
  )
}
```

## Contents

- [Getting started](#getting-started)
- [Implementation Status](#implementation-status)
- [Create Your Swagger](#create-your-swagger)
  - [General Swagger Info](#general-swagger-info)
  - [Adding Contact and License info (optional)](#adding-contact-and-license-info-optional)
  - [Adding Tags (optional)](#adding-tags-optional)
  - [Security](#security)
    - [Basic Auth](#basic-auth)
    - [API Key Auth](#api-key-auth)
    - [OAuth2](#oauth2-auth)
  - [Endpoints (API)](#endpoints-api)
    - [Arguments](#arguments)
- [Contribution](#contribution)

## Getting started

1. Get swagno package in your project

```sh
go get github.com/domhoward14/swagno
```

2. Import swagno

```go
import "github.com/domhoward14/swagno"
import "github.com/go-swagno/swagno-http/swagger" // recommended if you want to use their http handler for serving swagger docs
```

3. Create your endpoints (check [Endpoints](#endpoints-api)). Example:

```go
 endpoints := []*endpoint.EndPoint{
  endpoint.New(
   endpoint.WithMethod(endpoint.GET),
   endpoint.WithPath("/product/page"),
   endpoint.WithTags("product"),
   endpoint.WithSuccessfulReturns([]response.Info{models.UnsuccessfulResponse{}}),
   endpoint.WithErrors([]response.Info{models.EmptySuccessfulResponse{}}),
   endpoint.WithDescription(desc),
   endpoint.WithProduce([]string{"application/json", "application/xml"}),
   endpoint.WithConsume([]string{"application/json"}),
  ),
  endpoint.New(
   endpoint.WithMethod(endpoint.GET),
   endpoint.WithPath("/product"),
   endpoint.WithTags("product"),
   endpoint.WithParams(parameter.IntParam("id", parameter.WithRequired())),
   endpoint.WithSuccessfulReturns([]response.Info{models.SuccessfulResponse{}}),
   endpoint.WithErrors([]response.Info{models.UnsuccessfulResponse{}}),
  ),
  endpoint.New(
   endpoint.WithMethod(endpoint.GET),
   endpoint.WithPath("/product/{id}/detail"),
   endpoint.WithTags("product"),
   endpoint.WithParams(parameter.IntParam("id", parameter.WithRequired())),
   endpoint.WithSuccessfulReturns([]response.Info{models.SuccessfulResponse{}}),
   endpoint.WithErrors([]response.Info{models.UnsuccessfulResponse{}}),
  ),
  endpoint.New(
    endpoint.WithMethod(endpoint.POST),
    endpoint.WithPath("/product"),
    endpoint.WithTags("product"),
    endpoint.WithBody(models.ProductPost{}),
    endpoint.WithSuccessfulReturns([]response.Info{models.SuccessfulResponse{}}),
    endpoint.WithErrors([]response.Info{models.UnsuccessfulResponse{}}),
    endpoint.WithProduce([]mime.MIME{mime.JSON, mime.XML}),
  ),
}
```

4. Create Swagger(swagno) instance

```go
sw := swagno.New(swagno.Config{Title: "Testing API", Version: "v1.0.0"})
```

5. Use sw.AddEndpoints function to add endpoints arrays to Swagno

```go
sw.AddEndpoints(endpoints)

// you can also add more arrays
sw.AddEndpoints(productEndpoints)
sw.AddEndpoints(merchantEndpoints)
```

6. Generate json as string and give it to your handler to serve. You can create your own handler or use the swagno http handler

```go
 http.HandleFunc("/swagger/", swagger.SwaggerHandler(sw.GenerateDocs()))
 fmt.Println("Server is running on http://localhost:8080")
 http.ListenAndServe(":8080", nil)
```

## Implementation Status

As purpose of this section, you can compare **swagno** status with **swaggo**

[Swagger 2.0 document](https://swagger.io/docs/specification/2-0/basic-structure/)

- [x] Basic Structure
- [x] API Host and Base Path
- [x] Paths and Operations
- [x] Describing Parameters
- [x] Describing Request Body
- [x] Describing Responses
- [x] MIME Types -> need to improve
- [x] Authentication
  - [x] Basic Authentication
  - [x] API Keys
  - [x] OAuth2
- [ ] Adding Examples
- [x] File Upload -> need to improve
- [x] Enums
- [x] Grouping Operations With Tags
- [ ] Swagger Extensions

# Create Your Swagger

## General Swagger Info

You can use the swagger config when creating new swagger object

```go
type Config struct {
 Title   string   // title of the Swagger documentation
 Version string   // version of the Swagger documentation
 Host    string   // host URL for the API
 Path    string   // path to the Swagger JSON file
 License *License // license information for the Swagger documentation
 Contact *Contact // contact information for the Swagger documentation
}
```

```go
sw := swagno.New(swagno.Config{Title: "Testing API", Version: "v1.0.0"}) // optionally you can also use the License and Info properties as well
```

### Adding Tags (optional)

Allows adding meta data to a single tag. If you don't need meta data for your tags, you can skip this.

There is 3 alternative way for describing tags with descriptions.

```go
sw.AddTags(tag.Tag("product", "Product operations"), tag.NewTag("merchant", "Merchant operations"))
```

```go
sw.AddTags(tag.Tag{Name: "WithStruct", Description: "WithStruct operations"})
```

```go
sw.Tags = append(sw.Tags, tag.Tag{Name: "headerparams", Description: "headerparams operations"})
```

## Security

If you want to add security to your swagger, you can use `SetBasicAuth`, `SetApiKeyAuth`, `SetOAuth2Auth` functions.

```go
sw.SetBasicAuth()
sw.SetApiKeyAuth("api_key", "header")
sw.SetOAuth2Auth("oauth2_name", "password", "http://localhost:8080/oauth2/token", "http://localhost:8080/oauth2/authorize", Scopes(Scope("read:pets", "read your pets"), Scope("write:pets", "modify pets in your account")))
```

#### Basic Auth

If you have a basic auth with username and password, you can use `SetBasicAuth` function. It has default name as "basicAuth". You can add description as argument:

```go
sw.SetBasicAuth()
// with description
sw.SetBasicAuth("Description")
```

#### Api Key Auth

If you have an api key auth, you can use `SetApiKeyAuth` function.

Parameters:

- `name` -> name of the api key
- `in` -> location of the api key. It can be `header` or `query`
- `description` (optional) -> you can also add description as argument

```go
sw.SetApiKeyAuth("api_key", "header")
// with description
sw.SetApiKeyAuth("api_key", "header", "Description")
```

#### OAuth2 Auth

If you have an oauth2 auth, you can use `SetOAuth2Auth` function. You can also add description as argument:

Parameters:

- `name` -> name of the oauth2
- `flow` -> flow type of the oauth2. It can be `implicit`, `password`, `application`, `accessCode`
- `authorizationUrl` -> authorization url of the oauth2 (set this if flow is `impilicit` or `accessCode`, else you can set empty string)
- `tokenUrl` -> token url of the oauth2 (set this if flow is `password`, `application` or `accessCode`, else you can set empty string)
- `scopes` -> scopes of the oauth2
- `description` (optional) -> you can also add description as argument

```go
sw.SetOAuth2Auth("oauth2_name", "password", "", "http://localhost:8080/oauth2/token", Scopes(Scope("read:pets", "read your pets"), Scope("write:pets", "modify pets in your account")))
```

For scopes, you can use `Scopes` function. It takes `Scope` as variadic parameter.
Parameters of `Scope`:

- `name` -> name of the scope
- `description` -> description of the scope

## Endpoints (API)

Definition:

```go
type EndPoint struct {
 method            string
 path              string
 params            []*parameter.Parameter
 tags              []string
 Body              interface{}
 successfulReturns []response.Info
 errors            []response.Info
 description       string
 summary           string
 consume           []mime.MIME
 produce           []mime.MIME
 security          []map[string][]string
}
```

You need to create an Endpoint array []Endpoint and add your endpoints in this array. Example:

```go
import "github.com/domhoward14/swagno/components/endpoint"

endpoints := []endpoint.Endpoint{
  endpoint.New(
    endpoint.WithMethod(endpoint.POST),
    endpoint.WithPath("/product"),
    endpoint.WithTags("product"),
    endpoint.WithBody(models.ProductPost{}),
    endpoint.WithSuccessfulReturns([]response.Info{models.SuccessfulResponse{}}),
    endpoint.WithErrors([]response.Info{models.UnsuccessfulResponse{}}),
    endpoint.WithProduce([]mime.MIME{mime.JSON, mime.XML}),
    ),
}
// add endpoints array to Swagno
sw.AddEndpoints(endpoints)
```

**Note:** You can simply add only one endpoint by using `AddEndpoint(endpoint)`

### Configuring Endpoints

Arguments: The `Endpoint` object is configured via the `With<property>` commands provided in the `github.com/domhoward14/swagno/components/endpoint package`

| Function             | Description                                           |
|----------------------|-------------------------------------------------------|
| `WithMethod(method string)` | Configures the HTTP method of the `EndPoint`.          |
| `WithPath(path string)` | Sets the path for the `EndPoint`.                       |
| `WithParams(params []*parameter.Parameter)` | Adds parameters to the `EndPoint`.                            |
| `WithTags(tags ...string)` | Assigns tags to the `EndPoint` for grouping and categorization.   |
| `WithBody(body interface{})` | Sets the request body structure expected by the `EndPoint`.      |
| `WithSuccessfulReturns(successfulReturns ...response.Info)` | Defines the information for successful responses from the `EndPoint`. |
| `WithErrors(errors ...response.Info)` | Specifies the error responses the `EndPoint` could return.       |
| `WithDescription(description string)` | Provides a detailed description of what the `EndPoint` does.     |
| `WithSummary(summary string)` | Gives a brief summary of the `EndPoint` purpose.                  |
| `WithConsume(consume ...mime.MIME)` | Lists the MIME types the `EndPoint` can consume (input formats).  |
| `WithProduce(produce ...mime.MIME)` | Lists the MIME types the `EndPoint` can produce (output formats). |
| `WithSecurity(security ...map[string][]string)` | Sets security requirements for the `EndPoint`, such as required scopes or auth methods. |

❗ **Don't forget to add your endpoints array to Swagno prior to serving requests** ❗

```go
sw.AddEndpoints(endpoints)
```

### Method

Options: GET, POST, PUT, DELETE, OPTION, PATCH, HEAD

```go
// MethodType represents HTTP request methods.
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

### Path

Path of your endpoint without adding `query` parameter options
For example, you have endpoint as `/product/{id}?someParam=true` you need to add path as `/product/{id}` only, without query params.

```go
endpoint.WithPath("/product/{id}")
```

### Tags

Tag name as string to categorize a endpoint

### Parameters

You can use `endpoint.WithParams()` function to generate params array for an `Endpoint` object:

```go
// path should be -> /product/{merchant}/{id}
endpoint.WithParams(
    parameter.IntParam("id", parameter.WithIn(parameter.Path), parameter.WithRequired()), parameter.StrParam("merchant", parameter.WithIn(parameter.Path), parameter.WithRequired()),
    parameter.StrEnumParam("type", enumFileTypes, parameter.WithIn(parameter.Path), parameter.WithRequired(), parameter.WithDescription("The type of upload")),
    ),
    swagno.Params(StrParam("merchant", true, ""), IntParam("id", true, ""))
```

#### Parameter Types

Below are all the parameter types that the `EndPoint object can take as input`

Here is the markdown table for all functions listed in the code:

| Function                 | Description                                                           |
|--------------------------|-----------------------------------------------------------------------|
| `IntParam(name string, ...)` | Creates an integer parameter.                                         |
| `StrParam(name string, ...)` | Creates a string parameter.                                           |
| `BoolParam(name string, ...)` | Creates a boolean parameter.                                          |
| `FileParam(name string, ...)` | Creates a file upload parameter.                                      |
| `IntQuery(name string, ...)` | Creates an integer query parameter.                                   |
| `StrQuery(name string, ...)` | Creates a string query parameter.                                     |
| `BoolQuery(name string, ...)` | Creates a boolean query parameter.                                    |
| `IntHeader(name string, ...)` | Creates an integer header parameter.                                  |
| `StrHeader(name string, ...)` | Creates a string header parameter.                                    |
| `BoolHeader(name string, ...)` | Creates a boolean header parameter.                                   |
| `IntEnumParam(name string, arr []int64, ...)` | Creates an integer enum parameter.                                 |
| `StrEnumParam(name string, arr []string, ...)` | Creates a string enum parameter.                                   |
| `IntEnumQuery(name string, arr []int64, ...)` | Creates an integer enum query parameter.                           |
| `StrEnumQuery(name string, arr []string, ...)` | Creates a string enum query parameter.                             |
| `IntEnumHeader(name string, arr []int64, ...)` | Creates an integer enum header parameter.                          |
| `StrEnumHeader(name string, arr []string, ...)` | Creates a string enum header parameter.                           |
| `IntArrParam(name string, arr []int64, ...)` | Creates an integer array parameter.                                 |
| `StrArrParam(name string, arr []string, ...)` | Creates a string array parameter.                                   |
| `IntArrQuery(name string, arr []int64, ...)` | Creates an integer array query parameter.                           |
| `StrArrQuery(name string, arr []string, ...)` | Creates a string array query parameter.                             |
| `IntArrHeader(name string, arr []int64, ...)` | Creates an integer array header parameter.                          |
| `StrArrHeader(name string, arr []string, ...)` | Creates a string array header parameter.                           |

### Parameter Options

Just like the `endpoint` package, the `parameter` package also comes with a suite of functional `With<Option>` options to configure a parameter.

| Modifier Function              | Description                                              |
|--------------------------------|----------------------------------------------------------|
| `WithType(t ParamType)`        | Sets the type of a parameter (integer, string, boolean). |
| `WithIn(in Location)`          | Defines where the parameter is expected (query, header). |
| `WithRequired()`               | Makes the parameter required.                            |
| `WithDescription(description string)` | Provides a description for the parameter.                |
| `WithDefault(defaultValue interface{})`    | Sets a default value for the parameter.                  |
| `WithFormat(format string)`           | Sets the format field for the parameter.                       |
| `WithMin(min int)`                 | sets the Min field of a Parameter.             |
| `WithMax(max int)`                 | sets the Max field of a Parameter.             |
| `WithMinLen(minLen int)`           | sets the MinLen field of a Parameter.             |
| `WithMaxLen(maxLen int)`           | sets the MaxLen field of a Parameter.             |
| `WithPattern(pattern)`         | sets the Pattern field of a Parameter.     |
| `WithMaxItems(maxItems)`       | sets the WithMaxItems field of a Parameter.    |
| `WithMinItems(minItems)`       | sets the WithMinItems field of a Parameter.    |
| `WithUniqueItems(uniqueItems)` | Sets the WithUniqueItems filed of a Parameter                     |
| `WithMultipleOf(multipleOf)`   | Sets the WithMultipleOf filed of a Parameter         |
| `WithCollectionFormat(c)`      | Sets the WithCollectionFormat filed of a Parameter |

## Defining Models

You can define models for the following type of objects below. [Look here](example/models/generic.go) for specific model examples.

❗ When modeling error/successful responses they need to implement the interface below to properly generate their Swagger/OpenAPI output.❗

```go
package response

// Info is an interface for response information.
type Info interface {
 GetDescription() string
 GetReturnCode() string
}
```

### Successful Response/Return

use a struct model instance like `models.EmptySuccessfulResponse{}` or nil

### Error Response

use a struct model instance like `models.UnsuccessfulResponse{}` or nil

### Request Body

use a struct model instance like `models.PostBody{}` or nil

### Security

Before using this function, you need to define your security definitions in Swagno struct. For example:

```go
sw.SetBasicAuth()
sw.SetApiKeyAuth("api_key", "query")
sw.SetOAuth2Auth("oauth2_name", "password", "http://localhost:8080/oauth2/token", "http://localhost:8080/oauth2/authorize", Scopes(Scope("read:pets", "read your pets"), Scope("write:pets", "modify pets in your account")))
```

### Consumes/Produces (optional)

Supported MIME types:

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

Example:

```
endpoint.WithProduce([]mime.MIME{mime.JSON, mime.XML}),
endpoint.WithConsume([]mime.MIME{mime.JSON}),
```

**NOTE: If you used FileParam() in your endpoint, you don't need to add "multipart/form-data" to consumes. It will add automatically.**

# Contribution

We are welcome to any contribution. Swagno still has some missing features. Also we want to enrich handler implementations for other web frameworks.
