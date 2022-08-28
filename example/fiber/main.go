package main

import (
	"fmt"

	"github.com/anilsenay/swagno/example/fiber/handlers"
	"github.com/gofiber/fiber/v2"
)

func main() {

	handler := handlers.NewHandler()

	app := fiber.New()

	// set mock routes
	handler.SetRoutes(app)

	// set swagger routes
	handler.SetSwagger(app)

	app.Listen(fmt.Sprintf(
		"%s:%s",
		"localhost",
		"8080"),
	)

}
