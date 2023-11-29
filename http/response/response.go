package response

// Info is an interface for response information.
type Info interface {
	GetDescription() string
	GetReturnCode() string
}
