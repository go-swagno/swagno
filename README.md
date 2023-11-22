# swagno: _no annotations, no files, no command_

<img src="https://user-images.githubusercontent.com/1047345/188009539-ea9d0106-979d-4f98-83a3-0d7df6969c9f.png" alt="Swagno" align="right" width="200"/>

Swagno is an approach to create Swagger Documentation 2.0 without any **annotation**, **exported file** or any **command** to run.
You can declare your documentation details in code and get a json string to serve with a handler.

## About the Project

This project inspired by [Swaggo](https://github.com/swaggo/swag). Swaggo, uses annotations, exports files and needs to run by command. If you don't like this way, [Swag**no**](https://github.com/go-swagno/swagno) appears as a good alternative.

This project was then forked from [Swagno](https://github.com/go-swagno/swagno) with the goals of being more idiomatic and robust. In order to achieve these goals a lot of breaking changes were made and hence why this repo is separate as opposed to being merged into the original.

## Improvements
Some improvements from the orignal [Swagno](https://github.com/go-swagno/swagno) are shown below:

- More idiomatic and easier to read API calls
- Added more type safety
- Uses functional option paremeters to allow for more flexible and robust Endpoint creation
- Allows for more than one response and error model types to be used when making endpoints
- Bug fixes with OpenAPI output

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
			endpoint.WithParams([]parameter.Parameter{parameter.IntParam("id", parameter.WithRequired(true))}),
			endpoint.WithSuccessfulReturns([]response.Info{models.SuccessfulResponse{}}),
			endpoint.WithErrors([]response.Info{models.UnsuccessfulResponse{}}),
		),
		endpoint.New(
			endpoint.WithMethod(endpoint.GET),
			endpoint.WithPath("/product/{id}/detail"),
			endpoint.WithTags("product"),
			endpoint.WithParams([]parameter.Parameter{parameter.IntParam("id", parameter.WithRequired(true))}),
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
			endpoint.WithParams([]parameter.Parameter{parameter.IntParam("id", parameter.WithRequired(true))}),
			endpoint.WithSuccessfulReturns([]response.Info{models.SuccessfulResponse{}}),
			endpoint.WithErrors([]response.Info{models.UnsuccessfulResponse{}}),
		),
		endpoint.New(
			endpoint.WithMethod(endpoint.GET),
			endpoint.WithPath("/product/{id}/detail"),
			endpoint.WithTags("product"),
			endpoint.WithParams([]parameter.Parameter{parameter.IntParam("id", parameter.WithRequired(true))}),
			endpoint.WithSuccessfulReturns([]response.Info{models.SuccessfulResponse{}}),
			endpoint.WithErrors([]response.Info{models.UnsuccessfulResponse{}}),
		),
		endpoint.New(
			endpoint.WithMethod(endpoint.POST),
			endpoint.WithPath("/product"),
			endpoint.WithTags("product"),
			endpoint.WithParams([]parameter.Parameter{}),
			endpoint.WithBody(models.ProductPost{}),
			endpoint.WithSuccessfulReturns([]response.Info{models.SuccessfulResponse{}}),
			endpoint.WithErrors([]response.Info{models.UnsuccessfulResponse{}}),
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
- you can use the swagger config when creating new swagger object
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
EndPoint(method MethodType, path string, tags string, params []Parameter, body interface{}, ret interface{}, err interface{}, des string, secuirty []map[string][]string, args ...string)
```

You need to create an Endpoint array []Endpoint and add your endpoints in this array. Example:

```go
endpoints := []swagno.Endpoint{
  swagno.EndPoint(GET, "/product", "product", Params(), nil, []models.Product{}, models.ErrorResponse{}, "Get all products", nil),
  swagno.EndPoint(GET, "/product/{id}", "product", Params(IntParam("id", true, "")), nil, models.Product{}, models.ErrorResponse{}, "", nil),
  swagno.EndPoint(POST, "/product", "product", Params(), models.ProductPost{}, models.Product{}, models.ErrorResponse{}, "", nil),
}
// add endpoints array to Swagno
sw.AddEndpoints(endpoints)
```

**Note:** You can simply add only one endpoint by using `AddEndpoint(endpoint)`

- Arguments: (Method, Path, Tag, Params, Body, Response, Error Response, Description, Security)

**NOTE: If you not imported with explicit period (.), you need to get from swagno package:**

```go
endpoints := []swagno.Endpoint{
  swagno.EndPoint(swagno.GET, "/product", "product", swagno.Params(), nil, []models.Product{}, models.ErrorResponse{}, "Get all products", nil),
  swagno.EndPoint(swagno.GET, "/product/{id}", "product", swagno.Params(swagno.IntParam("id", true, "")), nil, models.Product{}, models.ErrorResponse{}, "", nil),
  swagno.EndPoint(swagno.POST, "/product", "product", swagno.Params(), models.ProductPost{}, models.Product{}, models.ErrorResponse{}, "", nil),
}
// add endpoints array to Swagno
sw.AddEndpoints(endpoints)
```

If you don't like this functional approach, you can use directly struct:

```go
endpoints := []Endpoint{
  {Method: "GET", Path: "/product/{id}", Description: "product", Params: Params(IntParam("id", true, "")), Return: models.Product{}, Error: models.ErrorResponse{}, Tags: []string{"WithStruct"}},
}
```

❗ **Don't forget to add your endpoints array to Swagno** ❗

```go
sw.AddEndpoints(endpoints)
```

### Arguments:

- [Method](#method)
- [Path](#path)
- [Tags](#tags)
- [Params](#params)
- [Body](#body)
- [Return](#responsereturn)
- [Error](#error-response)
- [Description](#description)
- [Security](#security)
- [Consumes](#consumes-optional) (optional / extra argument)
- [Produces](#produces-optional) (optional / extra argument)

### Method

Options: GET, POST, PUT, DELETE, OPTION, PATCH, HEAD

### Path

Path of your endpoint without adding `query` parameter options
For example, you have endpoint as `/product/{id}?someParam=true` you need to add path as `/product/{id}` only, without query params.

### Tags

Tags as string seperated by comma -> "tag1,tag2"

### Params

You can use Params() function to generate params array:

```go
// path should be -> /product/{merchant}/{id}
swagno.Params(StrParam("merchant", true, ""), IntParam("id", true, ""))
```

Or you can use []Parameter array:

```go
[]swagno.swagno.Parameter{{Name: "id", Type: "integer", In: "path", Required: true}}
```

#### Parameter Functions

- **IntParam** _(name string, required bool, description string, args ...Fields)_
- **StrParam** _(name string, required bool, description string, args ...Fields)_
- **BoolParam** _(name string, required bool, description string, args ...Fields)_
- **FileParam** _(name string, required bool, description string, args ...Fields)_
- **IntQuery** _(name string, required bool, description string, args ...Fields)_
- **StrQuery** _(name string, required bool, description string, args ...Fields)_
- **BoolQuery** _(name string, required bool, description string, args ...Fields)_
- **IntHeader** _(name string, required bool, description string, args ...Fields)_
- **StrHeader** _(name string, required bool, description string, args ...Fields)_
- **BoolHeader** _(name string, required bool, description string, args ...Fields)_
- **IntEnumParam** _(name string, arr []int64, required bool, description string, args ...Fields)_
- **StrEnumParam** _(name string, arr []string, required bool, description string, args ...Fields)_
- **IntEnumQuery** _(name string, arr []int64, required bool, description string, args ...Fields)_
- **StrEnumQuery** _(name string, arr []string, required bool, description string, args ...Fields)_
- **IntEnumHeader** _(name string, arr []int64, required bool, description string, args ...Fields)_
- **StrEnumHeader** _(name string, arr []string, required bool, description string, args ...Fields)_
- **IntArrParam** _(name string, arr []int64, required bool, description string, args ...Fields)_
- **StrArrParam** _(name string, arr []string, required bool, description string, args ...Fields)_
- **IntArrQuery** _(name string, arr []int64, required bool, description string, args ...Fields)_
- **StrArrQuery** _(name string, arr []string, required bool, description string, args ...Fields)_
- **IntArrHeader** _(name string, arr []int64, required bool, description string, args ...Fields)_
- **StrArrHeader** _(name string, arr []string, required bool, description string, args ...Fields)_

#### Parameter Options

| Parameter Option  | Description                                                                                  |
| ----------------- | -------------------------------------------------------------------------------------------- |
| Name              | name of parameter                                                                            |
| Type              | type of parameter: integer, number(for float/double), string, array, boolean, file           |
| In                | options: path, query, formData, header, array                                                |
| Required          | true or false                                                                                |
| Description       | parameter description as string                                                              |
| Enum              | int64 array or string array                                                                  |
| Items             |                                                                                              |
| Default           | default value of parameter                                                                   |
| Format            | format of parameter: https://swagger.io/specification/v2/#dataTypeFormat                     |
| Min               | min value of parameter value                                                                 |
| Max               | max value of parameter value                                                                 |
| MinLen            | min length of parameter value                                                                |
| MaxLen            | max length of parameter value                                                                |
| Pattern           | see: https://datatracker.ietf.org/doc/html/draft-fge-json-schema-validation-00#section-5.2.3 |
| MaxItems          | max items if type is array                                                                   |
| MinItems          | min items if type is array                                                                   |
| UniqueItems       | true or false                                                                                |
| MultipleOf        | see: https://datatracker.ietf.org/doc/html/draft-fge-json-schema-validation-00#section-5.1.1 |
| CollenctionFormat | if type is "array", checkout the table above:                                                |

| CollenctionFormat | Description                                                                                                                                                                |
| ----------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| csv               | comma separated values foo,bar.                                                                                                                                            |
| ssv               | space separated values foo bar.                                                                                                                                            |
| tsv               | tab separated values foo\tbar.                                                                                                                                             |
| pipes             | pipe separated values foo \| bar.                                                                                                                                          |
| multi             | corresponds to multiple parameter instances instead of multiple values for a single instance foo=bar&foo=baz. This is valid only for parameters in "query" or "formData".  |

### Body

use a struct model instance like `models.ProductPost{}` or nil

### Response/Return

use a struct model instance like `models.Product{}` or nil

### Error Response

use a struct model instance like `models.ErrorResponse` or nil

### Description

Endpoint description as string

### Security

Before using this function, you need to define your security definitions in Swagno struct. For example:

```go
sw.SetBasicAuth()
sw.SetApiKeyAuth("api_key", "query")
sw.SetOAuth2Auth("oauth2_name", "password", "http://localhost:8080/oauth2/token", "http://localhost:8080/oauth2/authorize", Scopes(Scope("read:pets", "read your pets"), Scope("write:pets", "modify pets in your account")))
```

If you want to add security to your endpoint, you can use one of `BasicAuth()`, `ApiKeyAuth()`, `OAuth()` functions:

```go
swagno.BasicAuth()
```

```go
swagno.ApiKeyAuth("api_key")
```

```go
swagno.OAuth("oauth2_name", "read:pets")
// you can add more scope name as argument
swagno.OAuth("oauth2_name", "read:pets", "write:pets", "...")
```

And use in `EndPoint` function:

```go
swagno.EndPoint(GET, "/product", "product", swagno.(), nil, []models.Product{}, models.Error{}, "description", swagno.ApiKeyAuth("api_key", "header"))
```

You can add more than one security to your endpoint with `Security()` function:

```go
swagno.EndPoint(GET, "/product", "product", swagno.(), nil, []models.Product{}, models.ErrorResponse{}, "Get all products", swagno.Security(ApiKeyAuth("api_key", "header"), swagno.()))
```

#### BasicAuth

If you want to use basic auth to an endpoint, you can use `BasicAuth()` function.

```go
swagno.BasicAuth("Basic Auth Description")
```

#### ApiKeyAuth

If you want to use api key auth to an endpoint, you can use `ApiKeyAuth()` function. It needs name as argument. This name must match with one of your Swagno security definations.

```go
swagno.ApiKeyAuth("api_key")
```

#### OAuth2Auth

If you want to use oauth2 auth to an endpoint, you can use `OAuth2Auth()` function. It needs name as argument. This name must match with one of your Swagno security definations. Then you can add scopes as arguments:

```go
swagno.OAuth2Auth("oauth2_name", "read:pets", "write:pets")
```

### Consumes (optional)

For default there is only one consumes type: "application/json", you don't need to add it. If you want to add more consumes types, you can add them as string as seperated by commas to EndPoint function's extra option:

```go
swagno.EndPoint(GET, "/product", "product", swagno.Params(), nil, []models.Product{}, models.ErrorResponse{}, "Get all products", nil, "application/xml,text/plain"),
```

**NOTE: If you used FileParam() in your endpoint, you don't need to add "multipart/form-data" to consumes. It will add automatically.**

### Produces (optional)

For default there are two produces types: "application/json" and "application/xml", you don't need to add them. If you want to add more produces types, you can add them as string as seperated by commas to EndPoint function's extra option:

```go
// without extra consumes -> nil as consumes
swagno.EndPoint(GET, "/product", "product", swagno.Params(), nil, []models.Product{}, models.ErrorResponse{}, "Get all products", nil, nil, "application/xml,text/plain"),
// with extra consumes
swagno.EndPoint(GET, "/product", "product", swagno.Params(), nil, []models.Product{}, models.ErrorResponse{}, "Get all products", nil, "application/xml,text/plain", "text/plain,text/html"),
```

# Contribution

We are welcome to any contribution. Swagno still has some missing features. Also we want to enrich handler implementations for other web frameworks.
