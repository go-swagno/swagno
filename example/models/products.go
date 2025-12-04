package models

import (
	"database/sql"
	"time"
)

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

	OptionString  sql.NullString  `json:"option_string"`
	OptionInt16   sql.NullInt16   `json:"option_int16"`
	OptionInt32   sql.NullInt32   `json:"option_int32"`
	OptionInt64   sql.NullInt64   `json:"option_int64"`
	OptionFloat64 sql.NullFloat64 `json:"option_float64"`
	OptionByte    sql.NullByte    `json:"option_byte"`
	OptionBool    sql.NullBool    `json:"option_bool"`
	OptionTime    sql.NullTime    `json:"option_time"`
}

type ErrorResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

type MerchantPageResponse struct {
	Brochures    []Product `json:"products"`
	MerchantName string    `json:"merchantName"`
}
