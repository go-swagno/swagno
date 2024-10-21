package models

import "time"

type Product struct {
	Id             uint64                    `json:"id"`
	Name           string                    `json:"name"`
	MerchantId     uint64                    `json:"merchant_id"`
	CategoryId     *uint64                   `json:"category_id,omitempty"`
	Tags           []uint64                  `json:"tags"`
	Images         []*string                 `json:"image_ids"`
	ImagesPtr      *[]string                 `json:"image_ids_ptr"`
	Sizes          []Sizes                   `json:"sizes"`
	SizePtrs       []*Sizes                  `json:"size_ptrs"`
	SaleDate       time.Time                 `json:"sale_date"`
	EndDate        *time.Time                `json:"end_date"`
	Complex        ComplexSuccessfulResponse `json:"complex"`
	Interface      interface{}               `json:"interface"`
	OmitEmpty      string                    `json:"omitemptytest,omitempty"`
	RequiredField  interface{}               `json:"required_field,omitempty" required:"true"`
	EmbeddedStruct EmbeddedStruct            `json:"embedded_struct"`
}

type EmbeddedStruct struct {
	Sizes
	OtherField int `json:"other_field"`
}

type Sizes struct {
	Size string `json:"size"`
}

type ProductPost struct {
	Name       string  `json:"name" example:"John Smith"`
	MerchantId uint64  `json:"merchant_id" example:"123456"`
	CategoryId *uint64 `json:"category_id,omitempty" example:"123"`
}

type ErrorResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

type MerchantPageResponse struct {
	Brochures    []Product `json:"products"`
	MerchantName string    `json:"merchantName"`
}
