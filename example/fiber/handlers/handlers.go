package handlers

import (
	. "github.com/anilsenay/swagno"
	"github.com/anilsenay/swagno/example/fiber/models"
	"github.com/gofiber/fiber/v2"
	fiber_swagger "github.com/gofiber/swagger"
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
		EndPoint(GET, "/product", "product", Params(), nil, models.Product{}, models.ErrorResponse{}, "Get all products"),
		EndPoint(GET, "/product", "product", Params(IntParam("id", true, "")), nil, models.Product{}, models.ErrorResponse{}, ""),
		EndPoint(POST, "/product", "product", Params(), models.ProductPost{}, models.Product{}, models.ErrorResponse{}, ""),

		// /merchant/{merchantId}?id={id} -> get product of a merchant
		EndPoint(GET, "/merchant", "merchant", Params(StrParam("merchant", true, ""), IntQuery("id", true, "product id")), nil, models.Product{}, models.ErrorResponse{}, ""),
	}

	sw := CreateSwagger("Service Brochure Go", "1.0")
	sw.Register(endpoints)
	// sw.ExportSwaggerDocs("api/swagger/doc.json") // optional

	a.Get("/swagger/*", func(ctx *fiber.Ctx) error {
		err := fiber_swagger.HandlerDefault(ctx)
		return err
	})
}
