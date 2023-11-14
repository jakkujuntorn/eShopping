package models

import (
	"gorm.io/gorm"
)

type ProductRequest struct {
	Title       string `json:"title" db:"title" validate:"required,notblank,haveSpace"`
	Price       int    `json:"price" db:"price" validate:"required,notblank,haveSpace"`
	Quantity    int    `json:"quantity" db:"quantity" validate:"required,notblank,haveSpace"`
	Description string `json:"description" db:"description" validate:"required,notblank,checkScript"`
	Category    string `json:"category" db:"category" validate:"required,notblank,haveSpace"`
}

type ProductUpdate struct {
	Title       string `json:"title" db:"title" validate:"required,notblank"`
	Price       int    `json:"price" db:"price" validate:"required,notblank"`
	Quantity    int    `json:"quantity" db:"quantity" validate:"required,notblank"`
	Description string `json:"description" db:"description" validate:"required,notblank,checkScript"`
	gorm.Model
}

// type product_pagination struct {
// 	Total    int
// 	Per_Page int
// }

type ProductResponse struct {
	IdStore     int    `json:"-"`
	IdProduct   int    `json:"id_product"`
	Title       string `json:"title"`
	Price       int    `json:"price"`
	Quantity    int    `json:"quantity"`
	Description string `json:"description"`
	Category    string `json:"category"`
}

type Product struct {
	// id product DB สร้างเอง
	// IdProduct   int    `json:"id_product" db:"id_product"`
	// id_store ดึงมากจาก id user
	// Id          int    `json:"id" db:"id" validate:"required"`
	IdProduct   int    `json:"id_product"  db:"id_product"`
	IdStore     int    `json:"id_store" db:"id_store" validate:"required,notblank,haveSpace"`
	Title       string `json:"title" db:"title" validate:"required,notblank,haveSpace,checkScript"`
	Price       int    `json:"price" db:"price" validate:"required,notblank"`
	Quantity    int    `json:"quantity" db:"quantity" validate:"required,notblank"`
	Description string `json:"description" db:"description" validate:"required,notblank,checkScript"`
	Category    string `json:"category" db:"category" validate:"required,notblank,haveSpace"`
	gorm.Model
}

type Picture_Product struct {
	IdProduct int    `json:"id_product" db:"id_product" validate:"required,notblank,haveSpace"`
	URL       string `json:"url" db:"url"`
}

type Product_Search struct {
	Title      string
	PriceStart int
	PriceEnd   int
	Category   string
}
