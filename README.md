# swagno: _no annotations, no files, no command_

<img src="https://user-images.githubusercontent.com/1047345/188009539-ea9d0106-979d-4f98-83a3-0d7df6969c9f.png" alt="Swagno" align="right" width="200"/>

Swagno is an approach to create Swagger Documentation 2.0 without any **annotation**, **exported file** or any **command** to run.
You can declare your documentation details in code and get a json string to serve with a handler.

## About the Project

This project inspired by [Swaggo](https://github.com/swaggo/swag). Swaggo, uses annotations, exports files and needs to run by command. If you don't like this way, [Swag**no**](https://github.com/go-swagno/swagno) appears as a good alternative.

## Contents

- [Getting started](#getting-started)
- [Supported Web Frameworks](#supported-web-frameworks)
- [How to use with Fiber](#how-to-use-with-gin)
- [How to use with Gin](#how-to-use-with-gin)
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
go get github.com/go-swagno/swagno
```

2. Import swagno (We suggest "." import)

```go
import (
  . "github.com/go-swagno/swagno"
)
```

You can import without explicit period (.) like this: `import "github.com/go-swagno/swagno"` but then you have to add `swagno.` to functions, structs etc. ( `[]swagno.Endpoint` , `swagno.EndPoint` , `swagno.Params()` etc.)

3. Create your endpoints (check [Endpoints](#endpoints-api)). Example:

```go
endpoints := []Endpoint{
  EndPoint(GET, "/product", "product", Params(), nil, []models.Product{}, models.ErrorResponse{}, "Get all products", nil),
  EndPoint(GET, "/product", "product", Params(IntParam("id", true, "")), nil, models.Product{}, models.ErrorResponse{}, "", nil),
  EndPoint(POST, "/product", "product", Params(), models.ProductPost{}, models.Product{}, models.ErrorResponse{}, "", nil),
}
```

4. Create Swagger(swagno) instance

```go
sw := CreateSwagger("Swagger API", "1.0")
```

5. Use AddEndpoints _(or swagno.AddEndpoints)_ function to add endpoints arrays to Swagno

```go
AddEndpoints(endpoints)
// you can add more arrays
// AddEndpoints(productEndpoints)
// AddEndpoints(merchantEndpoints)
```

6. Generate json as string and give it to your handler to serve. You can create your own handler or use our [Supported Web Frameworks](#supported-web-frameworks)

`sw.GenerateDocs()` -> to generate swagger json from endpoints

**For Gin:** [swagno-gin](https://github.com/go-swagno/swagno-gin)

```go
// gin example -> https://github.com/go-swagno/swagno-gin
a.GET("/swagger/*any", swagger.SwaggerHandler(sw.GenerateDocs()))
```

**For Fiber:** [swagno-fiber](https://github.com/go-swagno/swagno-fiber)

```go
// fiber example -> https://github.com/go-swagno/swagno-fiber
swagger.SwaggerHandler(a, sw.GenerateDocs(), swagger.Config{Prefix: "/swagger"})
```

## Supported Web Frameworks

- [fiber](https://github.com/go-swagno/swagno-fiber)
- [gin](https://github.com/go-swagno/swagno-gin)
- ... more on the way

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
swagger.SwaggerHandler(a, sw.GenerateDocs(), swagger.Config{Prefix: "/swagger"})
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

```go
sw := CreateSwagger("Swagger API", "1.0") -> (title, version)
sw := CreateSwagger("Swagger API", "1.0", "/v2", "localhost") -> (title, version, basePath, host)
```

### Adding Contact and License info (optional)

```
sw.Info.Contact.Email = "anilsenay3@gmail.com"
sw.Info.Contact.Name = "anilsenay"
sw.Info.Contact.Url = "https://anilsenay.com"
sw.Info.License.Name = "Apache 2.0"
sw.Info.License.Url = "http://www.apache.org/licenses/LICENSE-2.0.html"
sw.Info.TermsOfService = "http://swagger.io/terms/"
```

### Adding Tags (optional)

Allows adding meta data to a single tag. If you don't need meta data for your tags, you can skip this.

There is 3 alternative way for describing tags with descriptions.

```go
sw.AddTags(Tag("product", "Product operations"), Tag("merchant", "Merchant operations"))
```

```go
sw.AddTags(SwaggerTag{Name: "WithStruct", Description: "WithStruct operations"})
```

```go
sw.Tags = append(sw.Tags, SwaggerTag{Name: "headerparams", Description: "headerparams operations"})
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

Defination:

```go
EndPoint(method MethodType, path string, tags string, params []Parameter, body interface{}, ret interface{}, err interface{}, des string, secuirty []map[string][]string, args ...string)
```

You need to create an Endpoint array []Endpoint and add your endpoints in this array. Example:

```go
endpoints := []Endpoint{
  EndPoint(GET, "/product", "product", Params(), nil, []models.Product{}, models.ErrorResponse{}, "Get all products", nil),
  EndPoint(GET, "/product", "product", Params(IntParam("id", true, "")), nil, models.Product{}, models.ErrorResponse{}, "", nil),
  EndPoint(POST, "/product", "product", Params(), models.ProductPost{}, models.Product{}, models.ErrorResponse{}, "", nil),
}
// add endpoints array to Swagno
AddEndpoints(endpoints)
```

- Arguments: (Method, Path, Tag, Params, Body, Response, Error Response, Description, Security)

**NOTE: If you not imported with explicit period (.), you need to get from swagno package:**

```go
endpoints := []swagno.Endpoint{
  swagno.EndPoint(swagno.GET, "/product", "product", swagno.Params(), nil, []models.Product{}, models.ErrorResponse{}, "Get all products", nil),
  swagno.EndPoint(swagno.GET, "/product", "product", swagno.Params(swagno.IntParam("id", true, "")), nil, models.Product{}, models.ErrorResponse{}, "", nil),
  swagno.EndPoint(swagno.POST, "/product", "product", swagno.Params(), models.ProductPost{}, models.Product{}, models.ErrorResponse{}, "", nil),
}
// add endpoints array to Swagno
swagno.AddEndpoints(endpoints)
```

If you don't like this functional approach, you can use directly struct:

```go
endpoints := []Endpoint{
  {Method: "GET", Path: "/product4", Description: "product", Params: Params(IntParam("id", true, "")), Return: models.Product{}, Error: models.ErrorResponse{}, Tags: []string{"WithStruct"}},
}
```

❗ **Don't forget to add your endpoints array to Swagno** ❗

```go
AddEndpoints(endpoints)
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

Path of your endpoint without adding parameter options
For example, you have endpoint as `/product/{id}?someParam=true` you need to add path as "/product" only, without params.

### Tags

Tags as string seperated by comma -> "tag1,tag2"

### Params

You can use Params() function to generate params array:

```go
Params(StrParam("merchant", true, ""), IntParam("id", true, "")) // -> /product/{merchant}/{id}
```

Or you can use []Parameter array:

```go
[]Parameter{{Name: "id", Type: "integer", In: "path", Required: true}}
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

| CollenctionFormat | Description	|
| ----------------- | ----------------- |
| csv               | comma separated values foo,bar.|
| ssv               | space separated values foo bar.|
| tsv               | tab separated values foo\tbar.|
| pipes             | pipe separated values foo \| bar. |
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
BasicAuth()
```

```go
ApiKeyAuth("api_key")
```

```go
OAuth("oauth2_name", "read:pets")
// you can add more scope name as argument
OAuth("oauth2_name", "read:pets", "write:pets", "...")
```

And use in `EndPoint` function:

```go
EndPoint(GET, "/product", "product", Params(), nil, []models.Product{}, models.Error{}, "description", ApiKeyAuth("api_key", "header"))
```

You can add more than one security to your endpoint with `Security()` function:

```go
EndPoint(GET, "/product", "product", Params(), nil, []models.Product{}, models.ErrorResponse{}, "Get all products", Security(ApiKeyAuth("api_key", "header"), BasicAuth()))
```

#### BasicAuth

If you want to use basic auth to an endpoint, you can use `BasicAuth()` function.

```go
BasicAuth("Basic Auth Description")
```

#### ApiKeyAuth

If you want to use api key auth to an endpoint, you can use `ApiKeyAuth()` function. It needs name as argument. This name must match with one of your Swagno security definations.

```go
ApiKeyAuth("api_key")
```

#### OAuth2Auth

If you want to use oauth2 auth to an endpoint, you can use `OAuth2Auth()` function. It needs name as argument. This name must match with one of your Swagno security definations. Then you can add scopes as arguments:

```go
OAuth2Auth("oauth2_name", "read:pets", "write:pets")
```

### Consumes (optional)

For default there is only one consumes type: "application/json", you don't need to add it. If you want to add more consumes types, you can add them as string as seperated by commas to EndPoint function's extra option:

```go
EndPoint(GET, "/product", "product", Params(), nil, []models.Product{}, models.ErrorResponse{}, "Get all products", nil, "application/xml,text/plain"),
```

**NOTE: If you used FileParam() in your endpoint, you don't need to add "multipart/form-data" to consumes. It will add automatically.**

### Produces (optional)

For default there are two produces types: "application/json" and "application/xml", you don't need to add them. If you want to add more produces types, you can add them as string as seperated by commas to EndPoint function's extra option:

```go
// without extra consumes -> nil as consumes
EndPoint(GET, "/product", "product", Params(), nil, []models.Product{}, models.ErrorResponse{}, "Get all products", nil, nil, "application/xml,text/plain"),
// with extra consumes
EndPoint(GET, "/product", "product", Params(), nil, []models.Product{}, models.ErrorResponse{}, "Get all products", nil, "application/xml,text/plain", "text/plain,text/html"),
```

# Contribution

We are welcome to any contribution. Swagno still has some missing features. Also we want to enrich handler implementations for other web frameworks.
