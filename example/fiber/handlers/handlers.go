package handlers

import (
	. "github.com/anilsenay/swagno"
	"github.com/anilsenay/swagno/example/fiber/models"
	swagger "github.com/anilsenay/swagno/fiber/handler"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
}

func NewHandler() Handler {
	return Handler{}
}

func (h Handler) SetRoutes(a *fiber.App) {
	a.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).SendString("Hello World!")
	}).Name("index")
}

func (h *Handler) SetSwagger(a *fiber.App) {
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
	}

	sw := CreateSwagger("Swagger API", "1.0")
	// sw.ExportSwaggerDocs("api/swagger/doc.json") // optional
	swagger.SwaggerHandler(a, sw.GenerateDocs(endpoints), swagger.Config{Prefix: "/swagger"})
}
