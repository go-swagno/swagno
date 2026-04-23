package swagno3

import (
	"errors"

	"github.com/go-swagno/swagno/v3/components/extensions"
)

// ErrExternalDocsURLRequired is returned when URL field is empty in ExternalDocs
var ErrExternalDocsURLRequired = errors.New("externalDocs URL field is required")

// ExternalDocs represents external documentation object
// https://spec.openapis.org/oas/v3.0.3#external-documentation-object
type ExternalDocs struct {
	Description string                `json:"description,omitempty"`
	URL         string                `json:"url"` // REQUIRED
	Extensions  extensions.Extensions `json:"-"`
}

func (ed ExternalDocs) MarshalJSON() ([]byte, error) {
	type alias ExternalDocs
	return extensions.Merge(alias(ed), ed.Extensions)
}

// NewExternalDocs creates a new ExternalDocs instance
func NewExternalDocs(url string, description string) *ExternalDocs {
	return &ExternalDocs{
		URL:         url,
		Description: description,
	}
}

// Validate validates the ExternalDocs according to OpenAPI 3.0.3 rules
func (ed *ExternalDocs) Validate() error {
	if ed.URL == "" {
		return ErrExternalDocsURLRequired
	}
	return nil
}
