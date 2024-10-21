package tag

// https://swagger.io/specification/v2/#tagObject
type Tag struct {
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	ExternalDocs *ExternalDocs `json:"externalDocs,omitempty"`
}

type ExternalDocs struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type TagOpts func(*Tag)

func WithExternalDocs(name string, description string) TagOpts {
	return func(t *Tag) {
		t.ExternalDocs = &ExternalDocs{
			Name:        name,
			Description: description,
		}
	}
}

// New returns a new Tag.
func New(name string, description string, opts ...TagOpts) Tag {
	t := Tag{
		Name:        name,
		Description: description,
	}

	for _, opt := range opts {
		opt(&t)
	}

	return t
}
