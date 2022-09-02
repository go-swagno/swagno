package handlers

import (
	. "github.com/go-swagno/swagno"
	"github.com/go-swagno/swagno/example/fiber/models"
	"github.com/gofiber/fiber/v2"
)

type MerchantHandler struct {
}

func NewMerchantHandler() MerchantHandler {
	return MerchantHandler{}
}

func (h MerchantHandler) SetMerchantRoutes(a *fiber.App) {
	a.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).SendString("Hello World!")
	}).Name("index")
}

var MerchantEndpoints = []Endpoint{
	// /merchant/{merchantId}?id={id} -> get product of a merchant
	EndPoint(GET, "/merchant", "merchant", Params(StrParam("merchant", true, ""), IntQuery("id", true, "product id")), nil, models.Product{}, models.ErrorResponse{}, ""),
	EndPoint(GET, "/merchant", "merchant", Params(StrQuery("merchant", true, "")), nil, models.Product{}, models.ErrorResponse{}, ""),
}
