# Swagno: _`no` annotations, `no` files, `no` commands_

![Swagno Logo](https://user-images.githubusercontent.com/1047345/188009539-ea9d0106-979d-4f98-83a3-0d7df6969c9f.png "Swagno")

Swagno redefines the way Swagger Documentation 2.0 is created, embedding documentation seamlessly into your codebase for a clutter-free, streamlined experience. This tool does away with the hassles of annotations, exported files, and command executions. Simplify your documentation process with Swagno. Embrace the ease: Swagno - no annotations, no exports, no commands!

## About the Project

This project inspired by [Swaggo](https://github.com/swaggo/swag). Swaggo, uses annotations, exports files and needs to run by command. If you don't like this way, [Swag**no**](https://github.com/go-swagno/swagno) appears as a good alternative.

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
- [Examples](example/)

## Getting started

0. Server Example [here](example/server.go)

1. Get swagno package in your project

```sh
go get github.com/go-swagno/swagno
```

2. Import swagno

```go
import "github.com/go-swagno/swagno"
import "github.com/go-swagno/swagno-http/swagger" // recommended if you want to use their http handler for serving swagger docs
```

3. Create your endpoints (check [Endpoints](#endpoints-api)) with it's corresponding parameters. Example:

```go
 endpoints := []*endpoint.EndPoint{
  endpoint.New(
   endpoint.WithMethod(endpoint.GET),
   endpoint.WithPath("/product/page"),
   endpoint.WithTags("product"),
   endpoint.WithSuccessfulReturns([]response.Info{models.UnsuccessfulResponse{}}),
   endpoint.WithErrors([]response.Info{models.EmptySuccessfulResponse{}}),
   endpoint.WithDescription(desc),
   endpoint.WithProduce([]mime.MIME{mime.JSON, mime.XML}),
   endpoint.WithConsume([]mime.MIME{mime.JSON}),
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

## Supported Web Frameworks
- [fiber](https://github.com/go-swagno/swagno-fiber)
- [gin](https://github.com/go-swagno/swagno-gin)
- [gorilla/mux](https://github.com/go-swagno/swagno-http)
- [net/http](https://github.com/go-swagno/swagno-http)

## How to use with Fiber

You can read detailed document and find better examples in [swagno-fiber](https://github.com/go-swagno/swagno-fiber)

Example:

1. Get swagno-fiber

```sh
go get github.com/go-swagno/swagno-fiber
```

2. Import swagno-fiber

```go
import "github.com/go-swagno/swagno-fiber/swagger"
```

3.

```go
...
// assume you declare your endpoints and "sw"(swagno) instance
swagger.SwaggerHandler(a, sw.MustToJson(), swagger.Config{Prefix: "/swagger"})
...
```

You can find a detailed example in [https://github.com/go-swagno/swagno/example/fiber](https://github.com/go-swagno/swagno/tree/master/example/fiber)

## How to use with Gin

You can read detailed document and find better examples in [swagno-gin](https://github.com/go-swagno/swagno-gin)

Example:

1. Get swagno-gin

```sh
go get github.com/go-swagno/swagno-gin
```

2. Import swagno-gin

```go
import "github.com/go-swagno/swagno-gin/swagger"
```

3.

```go
...
// assume you declare your endpoints and "sw"(swagno) instance
a.GET("/swagger/*any", swagger.SwaggerHandler(sw.GenerateDocs()))
...
```

You can find a detailed example in [https://github.com/go-swagno/swagno/example/gin](https://github.com/go-swagno/swagno/tree/master/example/gin)

## Implementation Status

As purpose of this section, you can compare **swagno** status with **swaggo**

[Swagger 2.0 document](https://swagger.io/docs/specification/2-0/basic-structure/)

See how Swagno compares to Swaggo in terms of Swagger 2.0 features:

- Basic Structure: ✅
- API Host and Base Path: ✅
- Paths and Operations: ✅
- Describing Parameters: ✅
- Describing Request Body: ✅
- Describing Responses: ✅
- MIME Types: 🔄 (Improvement needed)
- Authentication: ✅
- File Upload: 🔄 (Improvement needed)
- Enums: ✅
- Grouping Operations With Tags: ✅
- Swagger Extensions: 🔜 (Coming soon)
- Swagger Validation: 🔜 (Coming soon)


# Create Your Swagger

## General Swagger Info

Swagger v2.0 specifications can be found [here](https://swagger.io/specification/v2/)

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
import "github.com/go-swagno/swagno/components/endpoint"

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

### Endpoint Options

Arguments: The `Endpoint` object is configured via the `With<property>` functional options provided in the `github.com/go-swagno/swagno/components/endpoint package`

| Function             | Description                                           |
|----------------------|-------------------------------------------------------|
| `WithMethod(method string)` | Sets the HTTP method of the `EndPoint`.          |
| `WithPath(path string)` | Sets the path for the `EndPoint`.                       |
| `WithParams(params []*parameter.Parameter)` | Adds parameters to the `EndPoint`.                            |
| `WithTags(tags ...string)` | Assigns tags to the `EndPoint` for grouping and categorization.   |
| `WithBody(body interface{})` | Sets the request body structure expected by the `EndPoint`.      |
| `WithSuccessfulReturns(successfulReturns ...response.Info)` | Sets the successful responses from the `EndPoint`. Needs to implement the `response.Info` interface |
| `WithErrors(errors ...response.Info)` | Sets the error responses the `EndPoint` could return. Needs to implement the `response.Info` interface       |
| `WithDescription(description string)` | Provides a detailed description of what the `EndPoint` does.     |
| `WithSummary(summary string)` | Gives a brief summary of the `EndPoint` purpose.                  |
| `WithConsume(consume ...mime.MIME)` | Sets the MIME types the `EndPoint` can consume (input formats).  |
| `WithProduce(produce ...mime.MIME)` | Sets the MIME types the `EndPoint` can produce (output formats). |
| `WithSecurity(security ...map[string][]string)` | Sets security requirements for the `EndPoint`, such as required scopes or auth methods. |

❗ **Don't forget to add your endpoints array to Swagno prior to serving requests** ❗

```go
sw.AddEndpoints(endpoints)
```

### Parameters

You can use `endpoint.WithParams()` function to generate params array for an `Endpoint` object:

```go
// path should be -> /product/{merchant}/{id}
endpoint.WithParams(
    parameter.StrParam("id", parameter.WithIn(parameter.Path), parameter.WithRequired()), 
    parameter.StrParam("merchant", parameter.WithIn(parameter.Path), parameter.WithRequired()),
    ),
```

### Parameter Location

Each parameter value can be assigned to a different location for the api request (i.e. [query, header, path, form]) using `WithIn`

```go
parameter.WithIn(parameter.Query)
```

| Location Type | Description                 |
|---------------|-----------------------------|
| `Query`       | Used for parameters in the URL query string. |
| `Header`      | Used for parameters in the HTTP header.      |
| `Path`        | Used for parameters within the path of the URL. |
| `Form`        | Used for parameters submitted through form data in POST requests. |

#### Parameter Types

Below are all the parameter types that the `EndPoint object can take as input`

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

Just like the `endpoint` package, the `parameter` package also comes with a set of functional `With<Option>` options to configure a parameter.

| Modifier Function              | Description                                              |
|--------------------------------|----------------------------------------------------------|
| `WithType(t ParamType)`        | Sets the type of a parameter (integer, string, boolean, and etc.). |
| `WithIn(in Location)`          | Defines where the parameter is expected (query, header). |
| `WithRequired()`               | Makes the parameter required.                            |
| `WithDescription(description string)` | Provides a description for the parameter.                |
| `WithDefault(defaultValue interface{})`    | Sets a default value for the parameter.                  |
| `WithFormat(format string)`           | Sets the format field for the parameter.                       |
| `WithMin(min int)`                 | sets the Min field of a Parameter.             |
| `WithMax(max int)`                 | sets the Max field of a Parameter.             |
| `WithMinLen(minLen int)`           | sets the MinLen field of a Parameter.             |
| `WithMaxLen(maxLen int)`           | sets the MaxLen field of a Parameter.             |
| `WithPattern(pattern string)`         | sets the Pattern field of a Parameter.     |
| `WithMaxItems(maxItems int)`       | sets the WithMaxItems field of a Parameter.    |
| `WithMinItems(minItems int)`       | sets the WithMinItems field of a Parameter.    |
| `WithUniqueItems(uniqueItems bool)` | Sets the WithUniqueItems filed of a Parameter                     |
| `WithMultipleOf(multipleOf int64)`   | Sets the WithMultipleOf filed of a Parameter         |
| `WithCollectionFormat(c CollectionFormat)`      | Sets the WithCollectionFormat filed of a Parameter |

## Defining Models

You can define models for the following type of objects below. [Look here](example/models/generic.go) for specific model examples.

❗ When modeling error/successful responses they need to implement the interface below to properly generate their Swagger/OpenAPI output.❗

```go
package response

// Info is an interface for response information.
type Info interface {
 Description() string
 ReturnCode() string
}
```

### Successful Response/Return

use a struct model instance like `models.EmptySuccessfulResponse{}` or nil

### Error Response

use a struct model instance like `models.UnsuccessfulResponse{}` or nil

### Request Body

use a struct model instance like `models.PostBody{}` or nil

### Security (optional)

Also provides functions to set different security configurations for swagger doc

```go
sw.SetBasicAuth()
sw.SetApiKeyAuth("api_key", "query")
sw.SetOAuth2Auth("oauth2_name", "password", "http://localhost:8080/oauth2/token", "http://localhost:8080/oauth2/authorize", security.Scopes(security.Scope("read:pets", "read your pets"), security.Scope("write:pets", "modify pets in your account")))
```

# Contribution

We are welcome to any contribution. Swagno still has some missing features. Also we want to enrich handler implementations for other web frameworks.
