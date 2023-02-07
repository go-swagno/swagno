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
		EndPoint(GET, "/product", "product", Params(), nil, []models.Product{}, models.ErrorResponse{}, "Get all products", nil).AddResponses(NewResponse("500", "Internal Server Error", models.ErrorResponse{})),
		EndPoint(GET, "/product/{id}", "product", Params(IntParam("id", true, "")), nil, models.Product{}, models.ErrorResponse{}, "", nil),
		EndPoint(POST, "/product", "product", Params(), models.ProductPost{}, models.Product{}, models.ErrorResponse{}, "", nil),

		// no return
		EndPoint(POST, "/product-no-return", "product", Params(), nil, nil, models.ErrorResponse{}, "", nil),
		// no error
		EndPoint(POST, "/product-no-error", "product", Params(), nil, nil, nil, "", nil),

		// ids query enum
		EndPoint(GET, "/products", "product", Params(IntEnumQuery("ids", []int64{1, 2, 3}, true, "")), nil, models.Product{}, models.ErrorResponse{}, "", nil),
		// ids path enum
		EndPoint(GET, "/products2/{ids}", "product", Params(IntEnumParam("ids", []int64{1, 2, 3}, true, "")), nil, models.Product{}, models.ErrorResponse{}, "", nil),
		// with fields
		EndPoint(GET, "/productsMinMax", "product", Params(IntArrQuery("ids", nil, true, "test", Fields{Min: 0, Max: 10, Default: 5})), nil, models.Product{}, models.ErrorResponse{}, "", nil),
		// string array query
		EndPoint(GET, "/productsArr", "product", Params(StrArrQuery("strs", nil, true, "")), nil, models.Product{}, models.ErrorResponse{}, "", nil),
		EndPoint(GET, "/productsArrWithEnums", "product", Params(StrArrQuery("strs", []string{"test1", "test2"}, true, "")), nil, models.Product{}, models.ErrorResponse{}, "", nil),
		EndPoint(GET, "/productsArrWithEnumsInPath/{strs}", "product", Params(StrArrParam("strs", []string{"test1", "test2"}, true, "")), nil, models.Product{}, models.ErrorResponse{}, "", nil),

		// /merchant/{merchantId}?id={id} -> get product of a merchant
		EndPoint(GET, "/merchant/{merchant}", "merchant", Params(StrParam("merchant", true, ""), IntQuery("id", true, "product id")), nil, models.Product{}, models.ErrorResponse{}, "", nil),

		// with headers
		EndPoint(POST, "/product-header", "header params", Params(IntHeader("header1", false, "")), models.ProductPost{}, models.Product{}, models.ErrorResponse{}, "", nil),
		EndPoint(POST, "/product2-header", "header params", Params(IntEnumHeader("header1", []int64{1, 2, 3}, false, ""), StrEnumHeader("header2", []string{"a", "b", "c"}, false, "")), models.ProductPost{}, models.Product{}, models.ErrorResponse{}, "", nil),
		EndPoint(POST, "/product3-header", "header params", Params(IntArrHeader("header1", []int64{1, 2, 3}, false, "")), models.ProductPost{}, models.Product{}, models.ErrorResponse{}, "", nil),

		// with file
		EndPoint(POST, "/productUpload", "upload", Params(FileParam("file", true, "File to upload")), nil, models.Product{}, models.ErrorResponse{}, "", nil),

		// without EndPoint function
		{Method: "GET", Path: "/product4/{id}", Description: "product", Params: Params(IntParam("id", true, "")), Responses: []Response{{Code: "200", Description: "OK", Body: models.Product{}}, {Code: "400", Description: "Not Found", Body: models.ErrorResponse{}}}, Tags: []string{"WithStruct"}},
		// without EndPoint function and without Params
		{Method: "GET", Path: "/product5/{id}", Description: "product", Params: []Parameter{{Name: "id", Type: "integer", In: "path", Required: true}}, Responses: []Response{{Code: "200", Description: "OK", Body: models.Product{}}, {Code: "400", Description: "Not Found", Body: models.ErrorResponse{}}}, Tags: []string{"WithStruct"}},

		// with security
		EndPoint(POST, "/secure-product", "Secure", Params(), models.ProductPost{}, models.Product{}, models.ErrorResponse{}, "Only Basic Auth", BasicAuth()),
		EndPoint(POST, "/multi-secure-product", "Secure", Params(), models.ProductPost{}, models.Product{}, models.ErrorResponse{}, "Basic Auth + Api Key Auth", Security(ApiKeyAuth("api_key"), BasicAuth())),
		EndPoint(POST, "/secure-product-oauth", "Secure", Params(), models.ProductPost{}, models.Product{}, models.ErrorResponse{}, "OAuth", OAuth("oauth2_name", "read:pets")),
	}

	sw := CreateNewSwagger("Swagger API", "1.0")
	AddEndpoints(endpoints)

	// set auth
	sw.SetBasicAuth()
	sw.SetApiKeyAuth("api_key", "query")
	sw.SetOAuth2Auth("oauth2_name", "password", "http://localhost:8080/oauth2/token", "http://localhost:8080/oauth2/authorize", Scopes(Scope("read:pets", "read your pets"), Scope("write:pets", "modify pets in your account")))

	// 3 alternative way for describing tags with descriptions
	sw.AddTags(Tag("product", "Product operations"), Tag("merchant", "Merchant operations"))
	sw.AddTags(SwaggerTag{Name: "WithStruct", Description: "WithStruct operations"})
	sw.Tags = append(sw.Tags, SwaggerTag{Name: "headerparams", Description: "headerparams operations"})

	// if you want to export your swagger definition to a file
	// sw.ExportSwaggerDocs("api/swagger/doc.json") // optional

	a.GET("/swagger/*any", swagger.SwaggerHandler(sw.GenerateDocs()))
}
