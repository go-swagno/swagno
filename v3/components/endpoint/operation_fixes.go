package endpoint

// OperationExternalDocs represents external docs for an operation
type OperationExternalDocs struct {
	Description string `json:"description,omitempty"`
	URL         string `json:"url"` // REQUIRED
}

// NewOperationExternalDocs creates external docs for operation
func NewOperationExternalDocs(url string, description string) *OperationExternalDocs {
	return &OperationExternalDocs{
		URL:         url,
		Description: description,
	}
}

// OperationServer represents a server object for operations
type OperationServer struct {
	URL         string                             `json:"url"`
	Description string                             `json:"description,omitempty"`
	Variables   map[string]OperationServerVariable `json:"variables,omitempty"`
}

// OperationServerVariable represents a server variable object for operations
type OperationServerVariable struct {
	Enum        []string `json:"enum,omitempty"`
	Default     string   `json:"default"`
	Description string   `json:"description,omitempty"`
}

// NewOperationServer creates a new OperationServer instance
func NewOperationServer(url string, description string) *OperationServer {
	return &OperationServer{
		URL:         url,
		Description: description,
	}
}
