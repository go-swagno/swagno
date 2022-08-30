package main

import (
	"github.com/anilsenay/swagno/example/gin/handlers"
	"github.com/gin-gonic/gin"
)

func main() {

	handler := handlers.NewHandler()

	app := gin.Default()

	// set mock routes
	handler.SetRoutes(app)

	// set swagger routes
	handler.SetSwagger(app)

	app.Run()

}
