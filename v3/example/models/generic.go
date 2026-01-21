package models

import "time"

type PostBody struct {
	Name string `json:"name" example:"John Smith"`
	ID   uint64 `json:"id" example:"123456"`
}

type EmptySuccessfulResponse struct{}

type SuccessfulResponse struct {
	ID           string            `json:"ID" example:"1234-1234-1234-1234" desc:"Unique identifier"`
	IDs          []int             `json:"IDs" example:"[1,2,3,4]" desc:"List of IDs"`
	Interface    any               `json:"interface" example:"{\"key\":\"value\"}" desc:"Generic interface field"`
	Map          map[string]string `json:"map" example:"{\"key1\":\"value1\",\"key2\":\"value2\"}" desc:"Map field"`
	OptionalInfo *string           `json:"optional_info,omitempty" example:"Some optional info" desc:"An optional info field"`
	Struct       struct{}          `json:"struct" desc:"A nested struct field"`
	OptionalIDs  *[]int            `json:"optional_IDs,omitempty" example:"[5,6,7,8]" desc:"An optional list of IDs"`
	Time         time.Time         `json:"time" example:"2023-10-05T14:48:00Z" desc:"Timestamp field"`
}

type UnsuccessfulResponse struct {
	ErrorField1 string `json:"error_msg1"`
}

type PageNotFound struct {
	ErrorMsg2 string `json:"error_msg2"`
}
