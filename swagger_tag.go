package swagno

// https://swagger.io/specification/v2/#tagObject
type SwaggerTag struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func Tag(name string, description string) SwaggerTag {
	return SwaggerTag{
		Name:        name,
		Description: description,
	}
}

func (swagger *Swagger) AddTags(tags ...SwaggerTag) {
	swagger.Tags = append(swagger.Tags, tags...)
}