package models

import "time"

type Product struct {
	Id         uint64     `json:"id"`
	Name       string     `json:"name"`
	MerchantId uint64     `json:"merchant_id"`
	CategoryId *uint64    `json:"category_id,omitempty"`
	Tags       []uint64   `json:"tags"`
	Images     []string   `json:"image_ids"`
	Sizes      []Sizes    `json:"sizes"`
	SaleDate   time.Time  `json:"sale_date"`
	EndDate    *time.Time `json:"end_date"`
	// TODO need to move these to test section
	RecursiveTest *Product          `json:"recursive_test"`
	SwagnoTagTest int               `json:"swagno_tag_test" swagno:"object"`
	TestMap       map[string]*Sizes `json:"test_map"`
	InterfaceTest interface{}       `json:"interface_test"`
}

type Products []Product

type MapTest map[string]Product

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

func (s ErrorResponse) GetDescription() string {
	return "Error Processing Request"
}

func (s ErrorResponse) GetReturnCode() string {
	return "400"
}

type MerchantPageResponse struct {
	Brochures    []Product `json:"products"`
	MerchantName string    `json:"merchantName"`
}

func (s MerchantPageResponse) GetDescription() string {
	return "Request Accepted"
}

func (s MerchantPageResponse) GetReturnCode() string {
	return "201"
}
