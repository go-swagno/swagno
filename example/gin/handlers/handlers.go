package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	. "github.com/go-swagno/swagno"
	"github.com/go-swagno/swagno-gin/swagger"
	"github.com/go-swagno/swagno/example/gin/models"
)

type Handler struct {
}

func NewHandler() Handler {
	return Handler{}
}

func (h Handler) SetRoutes(a *gin.Engine) {
	a.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
}

func (h *Handler) SetSwagger(a *gin.Engine) {
	endpoints := []Endpoint{
		EndPoint(GET, "/product", "product", Params(), nil, []models.Product{}, models.ErrorResponse{}, "Get all products"),
		EndPoint(GET, "/product", "product", Params(IntParam("id", true, "")), nil, models.Product{}, models.ErrorResponse{}, ""),
		EndPoint(POST, "/product", "product", Params(), models.ProductPost{}, models.Product{}, models.ErrorResponse{}, ""),

		// ids query enum
		EndPoint(GET, "/products", "product", Params(IntEnumQuery("ids", []int64{1, 2, 3}, true, "")), nil, models.Product{}, models.ErrorResponse{}, ""),
		// ids path enum
		EndPoint(GET, "/products2", "product", Params(IntEnumParam("ids", []int64{1, 2, 3}, true, "")), nil, models.Product{}, models.ErrorResponse{}, ""),
		// with fields
		EndPoint(GET, "/productsMinMax", "product", Params(IntArrQuery("ids", nil, true, "test", Fields{Min: 0, Max: 10, Default: 5})), nil, models.Product{}, models.ErrorResponse{}, ""),
		// string array query
		EndPoint(GET, "/productsArr", "product", Params(StrArrQuery("strs", nil, true, "")), nil, models.Product{}, models.ErrorResponse{}, ""),
		EndPoint(GET, "/productsArrWithEnums", "product", Params(StrArrQuery("strs", []string{"test1", "test2"}, true, "")), nil, models.Product{}, models.ErrorResponse{}, ""),
		EndPoint(GET, "/productsArrWithEnumsInPath", "product", Params(StrArrParam("strs", []string{"test1", "test2"}, true, "")), nil, models.Product{}, models.ErrorResponse{}, ""),

		// /merchant/{merchantId}?id={id} -> get product of a merchant
		EndPoint(GET, "/merchant", "merchant", Params(StrParam("merchant", true, ""), IntQuery("id", true, "product id")), nil, models.Product{}, models.ErrorResponse{}, ""),

		// with headers
		EndPoint(POST, "/product", "header params", Params(IntHeader("header1", false, "")), models.ProductPost{}, models.Product{}, models.ErrorResponse{}, ""),
		EndPoint(POST, "/product2", "header params", Params(IntEnumHeader("header1", []int64{1, 2, 3}, false, ""), StrEnumHeader("header2", []string{"a", "b", "c"}, false, "")), models.ProductPost{}, models.Product{}, models.ErrorResponse{}, ""),
		EndPoint(POST, "/product3", "header params", Params(IntArrHeader("header1", []int64{1, 2, 3}, false, "")), models.ProductPost{}, models.Product{}, models.ErrorResponse{}, ""),

		// with file
		EndPoint(POST, "/productUpload", "upload", Params(FileParam("file", true, "File to upload")), nil, models.Product{}, models.ErrorResponse{}, ""),

		// without EndPoint function
		{Method: "GET", Path: "/product4", Description: "product", Params: Params(IntParam("id", true, "")), Return: models.Product{}, Error: models.ErrorResponse{}, Tags: []string{"WithStruct"}},
		// without EndPoint function and without Params
		{Method: "GET", Path: "/product5", Description: "product", Params: []Parameter{{Name: "id", Type: "integer", In: "path", Required: true}}, Return: models.Product{}, Error: models.ErrorResponse{}, Tags: []string{"WithStruct"}},
	}

	sw := CreateSwagger("Swagger API", "1.0")

	// 3 alternative way for describing tags with descriptions
	sw.AddTags(Tag("product", "Product operations"), Tag("merchant", "Merchant operations"))
	sw.AddTags(SwaggerTag{Name: "WithStruct", Description: "WithStruct operations"})
	sw.Tags = append(sw.Tags, SwaggerTag{Name: "headerparams", Description: "headerparams operations"})

	// if you want to export your swagger definition to a file
	// sw.ExportSwaggerDocs("api/swagger/doc.json") // optional

	a.GET("/swagger/*any", swagger.SwaggerHandler(sw.GenerateDocs(endpoints)))
}
