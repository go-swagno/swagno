package models

type Object struct {
	Name string `json:"name" example:"John Smith"`
}

type Nested struct {
	Objects *[]Object `json:"objects,omitempty"`
	Strings *[]string `json:"strings,omitempty"`
}

type Deeply struct {
	Nested Nested `json:"nested"`
}

type ComplexSuccessfulResponse struct {
	Data *Deeply `json:"deeply"`
}
