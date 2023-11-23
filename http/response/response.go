package response

// Info is an interface for response information.
type Info interface {
	GetDescription() string
	GetReturnCode() string
}

// ErrorResponses is an interface for error responses.
type ErrorResponses interface {
	GetErrors() []Info
}
