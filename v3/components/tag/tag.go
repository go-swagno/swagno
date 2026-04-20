package tag

import "github.com/go-swagno/swagno/v3/components/extensions"

// https://spec.openapis.org/oas/v3.0.3#tag-object
type Tag struct {
	Name         string                `json:"name"`
	Description  string                `json:"description"`
	ExternalDocs *ExternalDocs         `json:"externalDocs,omitempty"`
	Extensions   extensions.Extensions `json:"-"`
}

func (t Tag) MarshalJSON() ([]byte, error) {
	type alias Tag
	return extensions.Merge(alias(t), t.Extensions)
}

type ExternalDocs struct {
	URL         string                `json:"url"`
	Description string                `json:"description,omitempty"`
	Extensions  extensions.Extensions `json:"-"`
}

func (ed ExternalDocs) MarshalJSON() ([]byte, error) {
	type alias ExternalDocs
	return extensions.Merge(alias(ed), ed.Extensions)
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
