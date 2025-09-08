package tag

// https://spec.openapis.org/oas/v3.0.3#tag-object
type Tag struct {
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	ExternalDocs *ExternalDocs `json:"externalDocs,omitempty"`
}

type ExternalDocs struct {
	URL         string `json:"url"`
	Description string `json:"description,omitempty"`
}

type TagOpts func(*Tag)

func WithExternalDocs(url string, description string) TagOpts {
	return func(t *Tag) {
		t.ExternalDocs = &ExternalDocs{
			URL:         url,
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
