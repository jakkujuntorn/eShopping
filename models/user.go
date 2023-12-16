package models

import (
	"gorm.io/gorm"
	// "github.com/go-playground/validator"
)

type UserRequest struct {
	Id          int    `json:"id"`
	Username    string `db:"username" json:"username" validate:"required,notblank,haveSpace" gorm:"unique"`
	Password    string `db:"password" json:"password" validate:"required,notblank,haveSpace"`
	Firstname   string `db:"firstname" json:"firstname" validate:"required,notblank"`
	Lastname    string `db:"lastname" json:"lastname" validate:"required,notblank"`
	Phonenumber string `db:"phonenumber" json:"honenumber" validate:"required"`
	Role        string `db:"role" json:"role"`
	Address  `json:"address"`
	// Address Address `json:"address"` // ถ้าใช้แบบนี้จะใสค่าแบบไหน
	gorm.Model
}

type UserResponse struct {
	Id          int    `json:"id"`
	Username    string `json:"username"`
	Firstname   string `json:"firstname"`
	Lastname    string `json:"lastname"`
	Phonenumber string `json:"phonenumber"`
	Role        string
	Address
}

type UserUpdate struct {
	Firstname   string `json:"firstname"`
	Lastname    string `json:"lastname"`
	Phonenumber string `json:"phonenumber"`
	Address
}

type UserLogin struct {
	Username string ` json:"username" db:"username" validate:"required,notblank,haveSpace"`
	Password string ` json:"password" db:"password" validate:"required,notblank,haveSpace"`
}

type Address struct {
	Address string `db:"address" json:"address" validate:"required,notblank"`
	City    string `db:"city" json:"city" validate:"required,notblank"`
	// ZipCode int  ตัวใหญ่กลางชื่อก็ไม่ได้
	Zipcode int `db:"zipcode" json:"zipcode" validate:"required,notblank"`
}
