package models

import (
	"gorm.io/gorm"
)

type Store struct {
	IdUser    string `db:"id_user" `
	IdStore string  `db:"id_store" `
	StoreName string `json:"storename" db:"store_name" validate:"required,notblank"`

	gorm.Model
}

type Store_Response struct {
	IdUser    string `db:"id_user"`
	StoreName string `json:"storename"`
	Products  []ProductResponse
	gorm.Model
}
