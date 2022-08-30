package swagger

import (
	"github.com/gin-contrib/static"

	"github.com/gin-gonic/gin"
)

type Config struct {
	Prefix string
}

var swaggerDoc string

var defaultConfig = Config{
	Prefix: "/swagger",
}

func SwaggerHandler(a *gin.Engine, doc []byte, config ...Config) {
	if len(config) != 0 {
		defaultConfig = config[0]
	}
	if swaggerDoc == "" {
		swaggerDoc = string(doc)
	}

	a.GET(defaultConfig.Prefix+"/doc.json", func(c *gin.Context) {
		c.String(200, swaggerDoc)
	})

	a.Use(static.Serve(defaultConfig.Prefix, static.LocalFile("./swaggerFiles", true)))
}
