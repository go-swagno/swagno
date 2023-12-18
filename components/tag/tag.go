package tag

// https://swagger.io/specification/v2/#tagObject
type Tag struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// New returns a new Tag.
func New(name string, description string) Tag {
	return Tag{
		Name:        name,
		Description: description,
	}
}
