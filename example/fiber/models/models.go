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
	SaleDate   *time.Time `json:"sale_date"`
	EndDate    time.Time  `json:"end_date"`
}

type Sizes struct {
	Size string `json:"size"`
}

type ProductPost struct {
	Name       string  `json:"name"`
	MerchantId uint64  `json:"merchant_id"`
	CategoryId *uint64 `json:"category_id,omitempty"`
}

type ErrorResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

type MerchantPageResponse struct {
	Brochures    []Product `json:"products"`
	MerchantName string    `json:"merchantName"`
}
