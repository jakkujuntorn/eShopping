package models

import (
	"time"

	"gorm.io/gorm"
)

// รูปแบบข้อมูลที่ font end ส่งมา
type CartRequest struct {
	// Id_User int        `json:"idUser" validate:"required"`
	Product []Product_CartItem `json:"product" validate:"required"`
}

// รูปแบบข้อมูลที่ font end ส่งมา ใช้สำหรับ แบ่งกลุ่ม idStore
type Product_CartItem struct {
	// Id_cart    int `json:"idCart"` // id carrt เอามาจาก Cart  id ต้องรอให้สร้าง cart ก่อน
	Id_Store   int `json:"idStore"  validate:"required"` // ใสใน DB เพื่อให้รู็ว่ามาจากร้านอะไร
	Id_Product int `json:"idProduct" validate:"required"`
	Quantity   int `json:"quantity" validate:"required"`
	Price      int `json:"price" validate:"required"`
}

//*****************  DB ****************
// ถ้าไม่ใช้ gorm model create at จะใสเวลาให้ไหม
type CartOrderDB struct {
	Id_User  int    `json:"idUser" db:"id_user" validate:"required"`
	Id_Store int    `json:"idStore" db:"id_store" validate:"required"`
	Status   string `json:"status" db:"status" validate:"required"`
	gorm.Model
}

type CartOrderDB_Mongo struct {
	Id_User  int       `json:"idUser" bson:"id_user" validate:"required"`
	Id_Store int       `json:"idStore" bson:"id_store" validate:"required"`
	Status   string    `json:"status" bson:"status" validate:"required"`
	CreateAt time.Time `json:"create_at" bson:"create_at"`
	UpdateAt time.Time `json:"update_at"bson:"update_at"`
	DeleteAt time.Time `json:"delete_at" bson:"delete_at"`
}

type CartItemDB_Mongo struct {
	Id_cart    string    `json:"idCart"  db:"id_cart" validate:"required"` // id carrt เอามาจาก Cart  id ต้องรอให้สร้าง cart ก่อน
	Id_Store   int       `json:"idStore" db:"id_store"  validate:"required"`
	Id_User    int       `json:"idUser"  db:"id_user" validate:"required"`
	Id_Product int       `json:"idProduct" db:"id_product" validate:"required"`
	Quantity   int       `json:"quantity" db:"quantity" validate:"required"`
	Price      int       `json:"price" db:"price" validate:"required"`
	CreateAt   time.Time `json:"create_at" bson:"create_at"`
	UpdateAt   time.Time `json:"update_at"bson:"update_at"`
	DeleteAt   time.Time `json:"delete_at" bson:"delete_at"`
}

// ใส id_user ลงไปด้วย
// เพราะ ตอนดึงจะดึงค่าตาม id_user ขึ้นมาทั้งหมด
// แบ่ง cartitems ตาม id_cart ด้วย map([int]models.cartItems)
type CartItemDB struct {
	Id_cart    int `json:"idCart"  db:"id_cart" validate:"required"` // id carrt เอามาจาก Cart  id ต้องรอให้สร้าง cart ก่อน
	Id_Store   int `json:"idStore" db:"id_store"  validate:"required"`
	Id_User    int `json:"idUser"  db:"id_user" validate:"required"`
	Id_Product int `json:"idProduct" db:"id_product" validate:"required"`
	Quantity   int `json:"quantity" db:"quantity" validate:"required"`
	Price      int `json:"price" db:"price" validate:"required"`
	gorm.Model
}

type CartItem_Postgres struct {
	ID         int
	Id_cart    int `json:"idCart"  db:"id_cart" validate:"required"` // id carrt เอามาจาก Cart  id ต้องรอให้สร้าง cart ก่อน
	Id_Store   int `json:"idStore" db:"id_store"  validate:"required"`
	Id_User    int `json:"idUser"  db:"id_user" validate:"required"`
	Id_Product int `json:"idProduct" db:"id_product" validate:"required"`
	Quantity   int `json:"quantity" db:"quantity" validate:"required"`
	Price      int `json:"price" db:"price" validate:"required"`
	gorm.Model
}

type Delete_GetCart_Postgres struct {
	ID       int    `json:"id" `
	Id_User  int    `json:"idUser"`
	Id_Store int    `json:"idStore"`
	Status   string `json:"status"`
	gorm.Model
}

type Delete_GetCartItem_Postgres struct {
	ID         int `json:"id" `
	Id_cart    int `json:"idCart"  db:"id_cart" validate:"required"` // id carrt เอามาจาก Cart  id ต้องรอให้สร้าง cart ก่อน
	Id_Store   int `json:"idStore" db:"id_store"  validate:"required"`
	Id_User    int `json:"idUser"  db:"id_user" validate:"required"`
	Id_Product int `json:"idProduct" db:"id_product" validate:"required"`
	Quantity   int `json:"quantity" db:"quantity" validate:"required"`
	Price      int `json:"price" db:"price" validate:"required"`
	gorm.Model
}

type Delete_carts struct {
	ID          int    `json:"id"`
	Id_UserDel  int    ` db:"id_user_del"`
	Id_StoreDel int    `json:"id_store_del" db:"id_store_del"`
	StatusDel   string `json:"status_del" db:"status_del"`
	gorm.Model
}

type Delete_cart_items struct {
	ID         int
	Id_cartDel    int `json:"id_cartDel" db:"id_cartDel"`
	Id_StoreDel   int `json:"id_storeDel" db:"id_storeDel"`
	Id_UserDel    int `json:"id_userDel" db:"id_userDel"`
	Id_ProductDEl int `json:"id_productDel" db:"id_productDel"`
	QuantityDel   int `json:"quantityDel" db:"quantityDel"`
	PriceDel      int `json:"priceDel" db:"priceDel"`
	gorm.Model
}

//*********** Resonse to font end
// for user
type CartOrderResponse struct {
	Id_User int `json:"idUser" db:"id_user" validate:"required"`
	Cart    []OrderCartDetail
}

type OrderCartDetail struct {
	Id_cart int
	Product []CartItemDB // gorm.Model ไม่ต้องส่งค่าออกมาก้ได้
	Total   int
	Status  string
}

// for store ไม่ใช้ และ
// type CartOrderResponse_Store struct {
// 	Id_Store int `json:"idUser" db:"id_user" validate:"required"`
// 	Order    []OrderCare_Store
// }

// type OrderCare_Store struct {
// 	Id_cart int
// 	Product []CartItemDB // gorm.Model ไม่ต้องส่งค่าออกมาก้ได้
// 	Total   int
// 	Status  string
// }

type StatusCartUpdate struct {
	Status string `json:"status" db:"status" validate:"required"`
}
