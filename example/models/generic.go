package models

type PostBody struct {
	Name string `json:"name" example:"John Smith"`
	ID   uint64 `json:"id" example:"123456"`
}

type EmptySuccessfulResponse struct{}

func (s EmptySuccessfulResponse) Description() string {
	return "OK"
}

func (s EmptySuccessfulResponse) ReturnCode() string {
	return "200"
}

type SuccessfulResponse struct {
	ID string `json:"ID" example:"1234-1234-1234-1234"`
}

func (s SuccessfulResponse) Description() string {
	return "Request Accepted"
}

func (s SuccessfulResponse) ReturnCode() string {
	return "201"
}

type UnsuccessfulResponse struct {
	ErrorField1 string `json:"error_msg1"`
}

func (u UnsuccessfulResponse) Description() string {
	return "Bad Request"
}

func (u UnsuccessfulResponse) ReturnCode() string {
	return "400"
}

type PageNotFound struct {
	ErrorMsg2 string `json:"error_msg2"`
}

func (u PageNotFound) Description() string {
	return "Page Not Found"
}

func (u PageNotFound) ReturnCode() string {
	return "404"
}
