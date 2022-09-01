# swagno: *no annotations, no files, no command*

Swagno is an approach to create Swagger Documentation 2.0 without any **annotation**, **exported file** or any **command** to run.
You can declare your documentation details in code and get a json string to serve with a handler.

## About the Project

This project inspired by [Swaggo](https://github.com/swaggo/swag). Swaggo, uses annotations, export files and need to run by command. If you don't like this way, [Swag**no**](https://github.com/go-swagno/swagno) appears as a good alternative. 

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
	- [Endpoints (API)](#endpoints-api)
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
You can import without explicit period (.) like this: `import "github.com/go-swagno/swagno"` but then you have to add **swagno.** to examples in this Readme file.

3. Create your endpoints (explained later in the document). Example:
```go
endpoints := []Endpoint{
	EndPoint(GET, "/product", "product", Params(), nil, []models.Product{}, models.ErrorResponse{}, "Get all products"),
	EndPoint(GET, "/product", "product", Params(IntParam("id", true, "")), nil, models.Product{}, models.ErrorResponse{}, ""),
	EndPoint(POST, "/product", "product", Params(), models.ProductPost{}, models.Product{}, models.ErrorResponse{}, ""),
}
```
3. Create Swagger(swagno) instance
```go
sw := CreateSwagger("Swagger API", "1.0")
```
4. Generate json as string and give it to your handler to serve. You can create your own handler or use our [Supported Web Frameworks](#supported-web-frameworks)

```go
// sw.GenerateDocs(endpoints) -> to generate swagger json from endpoints
---
// gin example -> https://github.com/go-swagno/swagno-gin
a.GET("/swagger/*any", swagger.SwaggerHandler(sw.GenerateDocs(endpoints)))
---
// fiber example -> https://github.com/go-swagno/swagno-fiber
swagger.SwaggerHandler(a, sw.GenerateDocs(endpoints), swagger.Config{Prefix: "/swagger"})
```

## Supported Web Frameworks
- [fiber](https://github.com/go-swagno/swagno-fiber)
- [gin](https://github.com/go-swagno/swagno-gin)
- ... more on the way

## How to use with Fiber
You can read detailed document and find better examples in [swagno-fiber](https://github.com/go-swagno/swagno-fiber)

Example:
1. Get swagno-fiber
```go
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
swagger.SwaggerHandler(a, sw.GenerateDocs(endpoints), swagger.Config{Prefix: "/swagger"})
...
```
You can find a detailed example in [https://github.com/go-swagno/swagno/example/fiber](https://github.com/go-swagno/swagno/tree/master/example/fiber)

## How to use with Gin
You can read detailed document and find better examples in [swagno-gin](https://github.com/go-swagno/swagno-gin)

Example:
1. Get swagno-gin
```go
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
a.GET("/swagger/*any", swagger.SwaggerHandler(sw.GenerateDocs(endpoints)))
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
- [ ] Authentication
  - [ ] Basic Authentication
  - [ ] API Keys
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

## Endpoints (API)

You need to create an Endpoint array []Endpoint and add your endpoints in this array. Example:
```go
endpoints := []Endpoint{
  EndPoint(GET, "/product", "product", Params(), nil, []models.Product{}, models.ErrorResponse{}, "Get all products"),
  EndPoint(GET, "/product", "product", Params(IntParam("id", true, "")), nil, models.Product{}, models.ErrorResponse{}, ""),
  EndPoint(POST, "/product", "product", Params(), models.ProductPost{}, models.Product{}, models.ErrorResponse{}, ""),
}
```
- Arguments: (Method, Path, Tag, Params, Body, Response, Error Response, Description) 

**NOTE: If you not imported with explicit period (.), you need to get from swagno package:**
```go
endpoints := []swagno.Endpoint{
  swagno.EndPoint(swagno.GET, "/product", "product", swagno.Params(), nil, []models.Product{}, models.ErrorResponse{}, "Get all products"),
  swagno.EndPoint(swagno.GET, "/product", "product", swagno.Params(swagno.IntParam("id", true, "")), nil, models.Product{}, models.ErrorResponse{}, ""),
  swagno.EndPoint(swagno.POST, "/product", "product", swagno.Params(), models.ProductPost{}, models.Product{}, models.ErrorResponse{}, ""),
}
```


If you don't like this functional approach, you can use directly struct:
```go
endpoints := []Endpoint{
  {Method: "GET", Path: "/product4", Description: "product", Params: Params(IntParam("id", true, "")), Return: models.Product{}, Error: models.ErrorResponse{}, Tags: []string{"WithStruct"}},
}
```

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
- **IntParam** *(name string, required bool, description string, args ...Fields)*
- **StrParam** *(name string, required bool, description string, args ...Fields)*
- **BoolParam** *(name string, required bool, description string, args ...Fields)*
- **FileParam** *(name string, required bool, description string, args ...Fields)*
- **IntQuery** *(name string, required bool, description string, args ...Fields)*
- **StrQuery** *(name string, required bool, description string, args ...Fields)*
- **BoolQuery** *(name string, required bool, description string, args ...Fields)*
- **IntHeader** *(name string, required bool, description string, args ...Fields)*
- **StrHeader** *(name string, required bool, description string, args ...Fields)*
- **BoolHeader** *(name string, required bool, description string, args ...Fields)*
- **IntEnumParam** *(name string, arr []int64, required bool, description string, args ...Fields)*
- **StrEnumParam** *(name string, arr []string, required bool, description string, args ...Fields)*
- **IntEnumQuery** *(name string, arr []int64, required bool, description string, args ...Fields)*
- **StrEnumQuery** *(name string, arr []string, required bool, description string, args ...Fields)*
- **IntEnumHeader** *(name string, arr []int64, required bool, description string, args ...Fields)*
- **StrEnumHeader** *(name string, arr []string, required bool, description string, args ...Fields)*
- **IntArrParam** *(name string, arr []int64, required bool, description string, args ...Fields)*
- **StrArrParam** *(name string, arr []string, required bool, description string, args ...Fields)*
- **IntArrQuery** *(name string, arr []int64, required bool, description string, args ...Fields)*
- **StrArrQuery** *(name string, arr []string, required bool, description string, args ...Fields)*
- **IntArrHeader** *(name string, arr []int64, required bool, description string, args ...Fields)*
- **StrArrHeader** *(name string, arr []string, required bool, description string, args ...Fields)*

#### Parameter Options
| Parameter Option  | Description   |
| ----------------- | ------------- |
| Name              | name of parameter |
| Type              | type of parameter: integer, number(for float/double), string, array, boolean, file      |
| In                | options: path, query, formData, header, array      |
| Required          | true or false |
| Description       | parameter description as string  |
| Enum              | int64 array or string array |
| Items             | |
| Default           | default value of parameter |
| Format            | format of parameter: https://swagger.io/specification/v2/#dataTypeFormat |
| Min               | min value of parameter value |
| Max               | max value of parameter value |
| MinLen            | min length of parameter value |
| MaxLen            | max length of parameter value |
| Pattern           | see: https://datatracker.ietf.org/doc/html/draft-fge-json-schema-validation-00#section-5.2.3 |
| MaxItems          | max items if type is array |
| MinItems          | min items if type is array |
| UniqueItems       | true or false |
| MultipleOf        | see: https://datatracker.ietf.org/doc/html/draft-fge-json-schema-validation-00#section-5.1.1 |
| CollenctionFormat | if type is "array", checkout the table above: |

| CollenctionFormat | Description   |
| ----------------- | ------------- |
| csv               | comma separated values foo,bar. |
| ssv               | space separated values foo bar. |
| tsv               | tab separated values foo\tbar.  |
| pipes             | pipe separated values foo|bar.  |
| multi             | corresponds to multiple parameter instances instead of multiple values for a single instance foo=bar&foo=baz. This is valid only for parameters in "query" or "formData". |

### Body
use a struct model instance like `models.ProductPost{}`

### Response/Return
use a struct model instance like `models.Product{}`

### Error Response
use a struct model instance like `models.ErrorResponse`

### Description
Endpoint description as string

# Contribution
We are welcome to any contribution. Swagno still has some missing features. Also we want to enrich handler implementations for other web frameworks.
