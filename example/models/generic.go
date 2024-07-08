package models

type PostBody struct {
	Name string `json:"name" example:"John Smith"`
	ID   uint64 `json:"id" example:"123456"`
}

type EmptySuccessfulResponse struct{}

type Response struct {
	Status string `json:"status"`
	Data   any    `json:"data,omitempty"`
	Errors any    `json:"errors,omitempty"`
}
type SuccessfulResponse struct {
	ID string `json:"ID" example:"1234-1234-1234-1234"`
}

type UnsuccessfulResponse struct {
	ErrorField1 string `json:"error_msg1"`
}

type PageNotFound struct {
	ErrorMsg2 string `json:"error_msg2"`
}
