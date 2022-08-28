package swagger

import (
	"github.com/swaggo/swag"
)

func (s Swagger) Register(endpoints []Endpoint) {
	var SwaggerInfo = &swag.Spec{
		Version:          s.Info.Version,
		Host:             s.Host,
		BasePath:         s.BasePath,
		Title:            s.Info.Title,
		InfoInstanceName: "swagger",
		SwaggerTemplate:  string(s.generateDocs(endpoints)),
	}

	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
