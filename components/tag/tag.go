package tag

// https://swagger.io/specification/v2/#tagObject
type Tag struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func NewTag(name string, description string) Tag {
	return Tag{
		Name:        name,
		Description: description,
	}
}
