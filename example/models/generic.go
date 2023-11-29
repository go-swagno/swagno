package models

type PostBody struct {
	Name string `json:"name" example:"John Smith"`
	ID   uint64 `json:"id" example:"123456"`
}

type EmptySuccessfulResponse struct{}

func (s EmptySuccessfulResponse) GetDescription() string {
	return "OK"
}

func (s EmptySuccessfulResponse) GetReturnCode() string {
	return "200"
}

type SuccessfulResponse struct {
	ID string `json:"ID" example:"1234-1234-1234-1234"`
}

func (s SuccessfulResponse) GetDescription() string {
	return "Request Accepted"
}

func (s SuccessfulResponse) GetReturnCode() string {
	return "201"
}

type UnsuccessfulResponse struct {
	ErrorField1 string `json:"error_msg1"`
}

func (u UnsuccessfulResponse) GetDescription() string {
	return "Bad Request"
}

func (u UnsuccessfulResponse) GetReturnCode() string {
	return "400"
}

type PageNotFound struct {
	ErrorMsg2 string `json:"error_msg2"`
}

func (u PageNotFound) GetDescription() string {
	return "Page Not Found"
}

func (u PageNotFound) GetReturnCode() string {
	return "404"
}
