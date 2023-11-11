package models

type SuccessfulResponse struct {
	ID string `json:"ID"`
}

func (s SuccessfulResponse) GetDescription() string {
	return "Request Accepted"
}

func (s SuccessfulResponse) GetReturnCode() string {
	return "201"
}

type UnsuccessfulResponse struct {
	ErrorMsg1 string `json:"error_msg1"`
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
