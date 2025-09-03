package endpoint

import (
	"fmt"
	"regexp"
	"strings"
)

// EnhancedCallback with runtime expression support
type EnhancedCallback struct {
	Expressions map[string]*PathItem `json:"-"` // Runtime expressions as keys
}

// EnhancedLink with full OpenAPI 3.0 support
type EnhancedLink struct {
	OperationRef string                 `json:"operationRef,omitempty"`
	OperationId  string                 `json:"operationId,omitempty"`
	Parameters   map[string]interface{} `json:"parameters,omitempty"`
	RequestBody  interface{}            `json:"requestBody,omitempty"`
	Description  string                 `json:"description,omitempty"`
	Server       *OperationServer       `json:"server,omitempty"`
}

// RuntimeExpression parser for OpenAPI 3.0 runtime expressions
type RuntimeExpression struct {
	Expression string
	Source     string // url, method, statusCode, request, response
	Location   string // header, query, path, body
	Pointer    string // JSON pointer for body references
}

// RuntimeExpressionType represents different types of runtime expressions
type RuntimeExpressionType int

const (
	URLExpression RuntimeExpressionType = iota
	MethodExpression
	StatusCodeExpression
	RequestExpression
	ResponseExpression
)

// NewEnhancedCallback creates a new enhanced callback
func NewEnhancedCallback() *EnhancedCallback {
	return &EnhancedCallback{
		Expressions: make(map[string]*PathItem),
	}
}

// AddExpression adds a runtime expression to the callback
func (ec *EnhancedCallback) AddExpression(expression string, pathItem *PathItem) error {
	if err := ValidateRuntimeExpression(expression); err != nil {
		return fmt.Errorf("invalid runtime expression '%s': %w", expression, err)
	}

	ec.Expressions[expression] = pathItem
	return nil
}

// GetExpression returns the PathItem for a given runtime expression
func (ec *EnhancedCallback) GetExpression(expression string) (*PathItem, bool) {
	pathItem, exists := ec.Expressions[expression]
	return pathItem, exists
}

// GetAllExpressions returns all runtime expressions in the callback
func (ec *EnhancedCallback) GetAllExpressions() map[string]*PathItem {
	result := make(map[string]*PathItem)
	for expr, pathItem := range ec.Expressions {
		result[expr] = pathItem
	}
	return result
}

// NewEnhancedLink creates a new enhanced link
func NewEnhancedLink() *EnhancedLink {
	return &EnhancedLink{
		Parameters: make(map[string]interface{}),
	}
}

// SetOperationRef sets the operation reference for the link
func (el *EnhancedLink) SetOperationRef(operationRef string) *EnhancedLink {
	el.OperationRef = operationRef
	el.OperationId = "" // Clear operationId if operationRef is set
	return el
}

// SetOperationId sets the operation ID for the link
func (el *EnhancedLink) SetOperationId(operationId string) *EnhancedLink {
	el.OperationId = operationId
	el.OperationRef = "" // Clear operationRef if operationId is set
	return el
}

// AddParameter adds a parameter to the link
func (el *EnhancedLink) AddParameter(name string, expression interface{}) *EnhancedLink {
	el.Parameters[name] = expression
	return el
}

// SetRequestBody sets the request body for the link
func (el *EnhancedLink) SetRequestBody(body interface{}) *EnhancedLink {
	el.RequestBody = body
	return el
}

// SetDescription sets the description for the link
func (el *EnhancedLink) SetDescription(description string) *EnhancedLink {
	el.Description = description
	return el
}

// SetServer sets the server for the link
func (el *EnhancedLink) SetServer(server *OperationServer) *EnhancedLink {
	el.Server = server
	return el
}

// ParseRuntimeExpression parses OpenAPI 3.0 runtime expressions
func ParseRuntimeExpression(expr string) (*RuntimeExpression, error) {
	expr = strings.TrimSpace(expr)

	if expr == "" {
		return nil, fmt.Errorf("empty runtime expression")
	}

	re := &RuntimeExpression{
		Expression: expr,
	}

	// Parse different types of runtime expressions
	switch {
	case expr == "$url":
		re.Source = "url"
		return re, nil
	case expr == "$method":
		re.Source = "method"
		return re, nil
	case expr == "$statusCode":
		re.Source = "statusCode"
		return re, nil
	case strings.HasPrefix(expr, "$request."):
		return parseRequestExpression(expr, re)
	case strings.HasPrefix(expr, "$response."):
		return parseResponseExpression(expr, re)
	default:
		return nil, fmt.Errorf("invalid runtime expression format")
	}
}

// parseRequestExpression parses request-type runtime expressions
func parseRequestExpression(expr string, re *RuntimeExpression) (*RuntimeExpression, error) {
	re.Source = "request"

	// Remove "$request." prefix
	remaining := strings.TrimPrefix(expr, "$request.")

	// Parse different request parts
	switch {
	case strings.HasPrefix(remaining, "header."):
		re.Location = "header"
		headerName := strings.TrimPrefix(remaining, "header.")
		if headerName == "" {
			return nil, fmt.Errorf("header name cannot be empty")
		}
		re.Pointer = headerName
	case strings.HasPrefix(remaining, "query."):
		re.Location = "query"
		queryName := strings.TrimPrefix(remaining, "query.")
		if queryName == "" {
			return nil, fmt.Errorf("query parameter name cannot be empty")
		}
		re.Pointer = queryName
	case strings.HasPrefix(remaining, "path."):
		re.Location = "path"
		pathName := strings.TrimPrefix(remaining, "path.")
		if pathName == "" {
			return nil, fmt.Errorf("path parameter name cannot be empty")
		}
		re.Pointer = pathName
	case strings.HasPrefix(remaining, "body"):
		re.Location = "body"
		if len(remaining) > 4 && remaining[4] == '#' {
			// JSON pointer in body
			re.Pointer = remaining[5:] // Remove "body#/"
		}
	default:
		return nil, fmt.Errorf("invalid request expression: %s", remaining)
	}

	return re, nil
}

// parseResponseExpression parses response-type runtime expressions
func parseResponseExpression(expr string, re *RuntimeExpression) (*RuntimeExpression, error) {
	re.Source = "response"

	// Remove "$response." prefix
	remaining := strings.TrimPrefix(expr, "$response.")

	// Parse different response parts
	switch {
	case strings.HasPrefix(remaining, "header."):
		re.Location = "header"
		headerName := strings.TrimPrefix(remaining, "header.")
		if headerName == "" {
			return nil, fmt.Errorf("header name cannot be empty")
		}
		re.Pointer = headerName
	case strings.HasPrefix(remaining, "body"):
		re.Location = "body"
		if len(remaining) > 4 && remaining[4] == '#' {
			// JSON pointer in body
			re.Pointer = remaining[5:] // Remove "body#/"
		}
	default:
		return nil, fmt.Errorf("invalid response expression: %s", remaining)
	}

	return re, nil
}

// ValidateRuntimeExpression validates a runtime expression according to OpenAPI 3.0.3 spec
func ValidateRuntimeExpression(expr string) error {
	// Basic validation patterns for runtime expressions
	validExpressions := []string{
		`^\$url$`,
		`^\$method$`,
		`^\$statusCode$`,
		`^\$request\.header\.[a-zA-Z0-9_-]+$`,
		`^\$request\.query\.[a-zA-Z0-9_-]+$`,
		`^\$request\.path\.[a-zA-Z0-9_-]+$`,
		`^\$request\.body(#/.+)?$`,
		`^\$response\.header\.[a-zA-Z0-9_-]+$`,
		`^\$response\.body(#/.+)?$`,
	}

	for _, pattern := range validExpressions {
		matched, err := regexp.MatchString(pattern, expr)
		if err != nil {
			return fmt.Errorf("regex error: %w", err)
		}
		if matched {
			return nil
		}
	}

	return fmt.Errorf("invalid runtime expression format")
}

// ValidateCallback validates a callback according to OpenAPI 3.0.3 rules
func (ec *EnhancedCallback) Validate() error {
	if len(ec.Expressions) == 0 {
		return fmt.Errorf("callback must contain at least one expression")
	}

	for expr := range ec.Expressions {
		if err := ValidateRuntimeExpression(expr); err != nil {
			return fmt.Errorf("invalid expression '%s': %w", expr, err)
		}
	}

	return nil
}

// ValidateLink validates a link according to OpenAPI 3.0.3 rules
func (el *EnhancedLink) Validate() error {
	// Either operationRef or operationId must be present, but not both
	if el.OperationRef == "" && el.OperationId == "" {
		return fmt.Errorf("link must have either operationRef or operationId")
	}

	if el.OperationRef != "" && el.OperationId != "" {
		return fmt.Errorf("link cannot have both operationRef and operationId")
	}

	// Validate parameter expressions
	for paramName, paramValue := range el.Parameters {
		if paramName == "" {
			return fmt.Errorf("parameter name cannot be empty")
		}

		// If parameter value is a string, it might be a runtime expression
		if strValue, ok := paramValue.(string); ok {
			if strings.HasPrefix(strValue, "$") {
				if err := ValidateRuntimeExpression(strValue); err != nil {
					return fmt.Errorf("invalid parameter expression for '%s': %w", paramName, err)
				}
			}
		}
	}

	return nil
}
