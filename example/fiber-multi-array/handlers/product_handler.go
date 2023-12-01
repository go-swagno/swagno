package handlers

import (
	. "github.com/go-swagno/swagno"
	"github.com/go-swagno/swagno/example/fiber/models"
	"github.com/gofiber/fiber/v2"
)

type ProductHandler struct {
}

func NewProductHandler() ProductHandler {
	return ProductHandler{}
}

func (h ProductHandler) SetProductRoutes(a *fiber.App) {
	a.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).SendString("Hello World!")
	}).Name("index")
}

var ProductEndpoints = []Endpoint{
	EndPoint(GET, "/product", "product", Params(), nil, []models.Product{}, models.ErrorResponse{}, "Get all products", nil),
	EndPoint(GET, "/product/{id}", "product", Params(IntParam("id", true, "")), nil, models.Product{}, models.ErrorResponse{}, "", nil),
	EndPoint(POST, "/product", "product", Params(), models.ProductPost{}, models.Product{}, models.ErrorResponse{}, "", nil),

	// /merchant/{merchantName}?id={id} -> get product of a merchant
	EndPoint(GET, "/merchant/{merchant}", "merchant", Params(StrParam("merchant", true, ""), IntQuery("id", true, "product id")), nil, models.Product{}, models.ErrorResponse{}, "", nil),
}
