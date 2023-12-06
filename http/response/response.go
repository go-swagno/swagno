package response

// Info is an interface for response information.
type Info interface {
	GetDescription() string
	GetReturnCode() string
}

type CustomResponseType struct {
	Model       any
	ReturnCode  string
	Description string
}

func (c CustomResponseType) GetDescription() string {
	return c.Description
}
func (c CustomResponseType) GetReturnCode() string {
	return c.ReturnCode
}

func New(model any, returnCode string, description string) Info {
	return CustomResponseType{
		Model:       model,
		ReturnCode:  returnCode,
		Description: description,
	}
}
