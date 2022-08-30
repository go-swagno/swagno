package swagger

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
)

type Config struct {
	Prefix string
}

var swaggerDoc string

var defaultConfig = Config{
	Prefix: "/swagger",
}

func SwaggerHandler(a *fiber.App, doc []byte, config ...Config) {
	if len(config) != 0 {
		defaultConfig = config[0]
	}
	if swaggerDoc == "" {
		swaggerDoc = string(doc)
	}
	a.Use(defaultConfig.Prefix+"/doc.json", func(c *fiber.Ctx) error {
		return c.SendString(swaggerDoc)
	})
	a.Use(defaultConfig.Prefix, filesystem.New(filesystem.Config{
		Root: http.Dir("./swaggerFiles"),
	}))
}
