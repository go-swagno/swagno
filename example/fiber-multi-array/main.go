package main

import (
	"fmt"

	swagno "github.com/go-swagno/swagno"
	"github.com/go-swagno/swagno-fiber/swagger"
	"github.com/go-swagno/swagno/example/fiber-multi-array/handlers"
	"github.com/gofiber/fiber/v2"
)

func main() {

	productHandler := handlers.NewProductHandler()
	merchantHandler := handlers.NewMerchantHandler()

	app := fiber.New()

	// set mock routes
	productHandler.SetProductRoutes(app)
	merchantHandler.SetMerchantRoutes(app)

	// set swagger routes
	sw := swagno.CreateNewSwagger("Swagger API", "1.0")
	swagno.AddEndpoints(handlers.ProductEndpoints)
	swagno.AddEndpoints(handlers.MerchantEndpoints)

	swagger.SwaggerHandler(app, sw.GenerateDocs(), swagger.Config{Prefix: "/swagger"})

	// Listen app
	app.Listen(fmt.Sprintf(
		"%s:%s",
		"localhost",
		"8080"),
	)

}
