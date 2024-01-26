package repository

import (
	"errors"
	"fmt"

	"myShopEcommerce/models"
	"time"
	"sync"
	"gorm.io/gorm"
)

type cart_repository struct {
	mysql *gorm.DB

	postgres *gorm.DB
	
	mt sync.Mutex
}

type ICart_Repositiry interface {
	// **********************************   mysql  **********************************
	//******** ไม่ได้ใช้ CreateCart_Repo *********
	CreateCart_Repo(cartOrder models.CartOrderDB) (cartOrderDB *models.CartOrderDB, err error)

	//********* ไม่ได้ใช้ CreateCartItem_Repo ********
	CreateCartItem_Repo(cartItem *[]models.CartItemDB) error

	CreateCart_Repo_V2(cartOrder map[int][]models.Product_CartItem, idUser int) error
	CheckStock(idProduct, quantity int, status string) error // ไม่ได้ใช้
	GetAllCartForUser_Repo(idUser int) (cartOrderDB []models.CartOrderDB, carItemst_User []models.CartItemDB, err error)
	GetAllCartForStore_Repo(idStore int) (cartOrderDB []models.CartOrderDB, cartItemDB_Struct []models.CartItemDB, err error)
	GetCartByIdForUser_Repo(idcart, idUser int) (carItemst_User []models.CartItemDB, err error)

	// อันนี้ต้องแก้เรื่องเวลาด้วย ****
	GetCartByIdForStore_Repo(idcart, idStore int) (carItemst_Store []models.CartItemDB, err error)
	EditCartOrder_Repo(cartItemDB []models.CartItemDB) (err error)
	EditStatusCartOrder_Repo(idCart, idUser int, status *models.StatusCartUpdate) (err error)
	DeleteCart_Repo(idUser, idcart int) error

	// ***********************  moongo ******************
	CreateCart_Repo_Mongo(cartOrder map[int][]models.Product_CartItem, idUser int) error
	EditCartOrder_Repo_Mongo(cartItemDB []models.CartItemDB) (err error)
	EditStatusCartOrder_Repo_Mongo(idCart, idUser int, status *models.StatusCartUpdate) (err error)

	// ****************************  Postgres SQL  ***********************************
	Transaction_Postgres(func(*cart_repository) error) error
	CreateCart_Repo_Postgres(cartOrder map[int][]models.Product_CartItem, idUser int) error
	EditCartOrder_Repo_Postgres(cartItemDB []models.CartItemDB) (err error)
	EditStatusCartOrder_Repo_Postgres(idCart, idUser int, status *models.StatusCartUpdate) (err error)
	GetAllCartForUser_Repo_Postgres(idUser int) (cartOrderDB []models.CartOrderDB, carItemst_User []models.CartItemDB, err error)
	GetAllCartForStore_Repo_Postgres(idStore int) (cartOrderDB []models.CartOrderDB, cartItemDB_Struct []models.CartItemDB, err error)
	GetCartByIdForUser_Repo_Postgres(idcart, idUser int) (carItemst_User []models.CartItemDB, err error)
	GetCartByIdForStore_Repo_Postgres(idcart, idStore int) (carItemst_Store []models.CartItemDB, err error)

	CheckStock_Postgres(idProduct, quantity int, status string) (error)
	ReStoreStock_Postgres(idProduct, quantity int) (error)
	DeleteCartForUser_Repo_Postgres(idUser, idCart int) error
	// Move_DeleteCart_Row([]models.Delete_carts,[]models.Delete_cart_items) error
}

func NewCart_Repo(mysql *gorm.DB, postgres *gorm.DB) *cart_repository {
	return &cart_repository{
		mysql: mysql,
		// monGo: monGo,
		postgres: postgres,
	}
}

var rollBack = make(chan string)
var commit = make(chan string)

// CreateCart_Repo implements ICart_Repositiry
// การตัด stock ตวร อยู่ตรงนี้ เพราะเป็นจุดแรกที่เริ่มทำงานด้าน  cart
func (cr *cart_repository) CreateCart_Repo(cart models.CartOrderDB) (cartOrder *models.CartOrderDB, err error) {
	// create Cart

	tx := cr.mysql.Begin()

	tx.Table("myshop.carts").Create(&cart)
	if tx.Error != nil {
		return &models.CartOrderDB{}, tx.Error
	}
	// สร้างตามจำนวน store
	// ดึงค่า cart เพื่อเอา id
	// cartOrder := models.CartOrderDB{}

	// tx = cr.db.Table("myshop.carts").
	tx.Table("myshop.carts").Last(&cartOrder)
	if tx.Error != nil {
		return &models.CartOrderDB{}, tx.Error
	}

	fmt.Println("************* wait Status **************")
	// มันรอตรงนี้ และไม่ทำงานต่อ
	// commit or rollback
	select {
	case <-rollBack:
		tx.Rollback()
	case <-commit:
		tx.Commit()
	}

	return cartOrder, tx.Error
}

func (cr *cart_repository) CreateCartItem_Repo(cartItems *[]models.CartItemDB) error {

	// ตัด stock
	// ดึงค่าตาม id product
	product := models.Product{}
	for _, v := range *cartItems {
		tx := cr.mysql.Table("myshop.products").Where("id_product=?", v.Id_Product).Find(&product)
		if tx.Error != nil {
			rollBack <- "rollback"
			return tx.Error
		}

		// check stock
		if product.Quantity < v.Quantity {
			rollBack <- "rollback"
			errorText := fmt.Sprintf("%v not enough", product.Title)
			return errors.New(errorText)
		}

		//  ตัดยอดตาม quantity
		product.Quantity = product.Quantity - v.Quantity

		// fmt.Println(product)
		// save product after cut stock
		tx = cr.mysql.Table("myshop.products").Where("id_product=?", product.IdProduct).Updates(&product)
		if tx.Error != nil {
			rollBack <- "rollback"
			return errors.New("can not update prpduct after cut stock")
		}
	}

	// create cart_items
	tx := cr.mysql.Table("myshop.cart_items").Create(&cartItems)
	if tx.Error != nil {
		rollBack <- "rollback"
		return tx.Error
	}
	fmt.Print("")

	// all seccress ต้องแจ้งไปบอก myshop.carts ให้ commit
	commit <- "commit"
	// last return
	return tx.Error
}

func (cr *cart_repository) CreateCart_Repo_V2(cartOrder map[int][]models.Product_CartItem, idUser int) error {

	// fmt.Println("********************************")
	// ร้าน 36 มี 2 สินค้า
	// ร้าน 22 มี 1 สินค้า

	// map[
	//22:[{22 22 22 22}]
	//36:[{36 63 1 18} {36 64 1 18}]
	//]

	// fmt.Println(cartOrder)
	// fmt.Println("*******************************")

	// *********************** ตัด stock *************
	// for idStore, product_InStore := range cartOrder {
	// 	_ = idStore
	// 	// 22 กับ 36
	// 	// fmt.Println(idStore)
	// 	// fmt.Println(product_InStore)
	// 	// 	36
	// 	// [{36 63 1 18} {36 64 1 18}]
	// 	// 22
	// 	// [{22 22 22 22}]
	// 	for i, v := range product_InStore { // เอาไว้แยก product แต่ละตัว
	// 		fmt.Print(i) //index  ธรรมดา
	// 		fmt.Print(v) // สินค้า แยกแต่ละ สินค้าเลย

	// 		// สินค้าแยกที่ละตัว พร้อมเอาไปตัด stock
	// 		errStock := cr.CheckStock(v.Id_Product, v.Quantity, "cut")
	// 		if errStock != nil {
	// 			return errStock
	// 		}
	// 	}
	// }

	tx := cr.mysql.Begin()
	//************* Sent Data toDB ******************
	for idStore, product_InStore := range cartOrder {
		cartOrderDB := models.CartOrderDB{
			Id_User:  idUser,
			Id_Store: idStore, // คือค่า key value ของ map (cartOrder)
			Status:   "pending",
		}

		if err := tx.Table("myshop.carts").Create(&cartOrderDB).Error; err != nil {
			return err
		}

		// ********************แยกข้อมูล CartItemDB ******************
		cartItemDB := []models.CartItemDB{}
		for _, data_CartItem := range product_InStore {

			// การเช็ค stock กับตัด stock
			errStock := cr.CheckStock(data_CartItem.Id_Product, data_CartItem.Quantity, "cut")
			if errStock != nil {
				tx.Rollback()
				return errStock
			}
			// ดึง id เตียมมา map กับ catItem
			cartOrderforCartItem := models.CartOrderDB{}
			// ดึงค่าสุดท้ายขึ้นมา มันคือ ค่าที่พึ่งใสลงไป
			// ต้องเอา cart id มา map  กับ cart_items idCart
			if err := tx.Table("myshop.carts").Last(&cartOrderforCartItem).Error; err != nil {
				tx.Rollback()
				return err
			}

			// ปั้นข้อมูล CartItemDB
			cartItem := models.CartItemDB{
				Id_cart: int(cartOrderforCartItem.ID), // เอามาจาก Db Cart
				// Id_cart:    int(<-idCart), // เอามาจาก Db Cart
				Id_Store:   data_CartItem.Id_Store,
				Id_User:    idUser,
				Id_Product: data_CartItem.Id_Product,
				Quantity:   data_CartItem.Quantity,
				Price:      data_CartItem.Price,
			}
			cartItemDB = append(cartItemDB, cartItem)
		}

		if err := tx.Table("myshop.cart_items").Create(&cartItemDB).Error; err != nil {
			// เปลี่ยน logic ใหม่ ไม่ต้องใช้ตรงนี้แล้ว
			//******** error ต้องคืนค่า stock ***********
			// for _, product_InStore := range cartOrder {
			// 	for _, v := range product_InStore { // เอาไว้แยก product แต่ละตัว
			// 		errStock := cr.CheckStock(v.Id_Product, v.Quantity, "return")
			// 		if errStock != nil {
			// 			return errStock
			// 		}
			// 	}
			// }

			//
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

func (cr *cart_repository) CheckStock(idProduct, quantity int, status string) error {
	if status == "cut" {
		cr.mt.Lock()
		product := models.Product{}
		tx := cr.mysql.Table("myshop.products").Where("id_product=?", idProduct).Find(&product)
		if tx.Error != nil {
			return tx.Error
		}
		// check stock
		if product.Quantity < quantity {
			errorText := fmt.Sprintf("%v not enough", product.Title)
			return errors.New(errorText)
		}

		//  ตัดยอดตาม quantity
		product.Quantity = product.Quantity - quantity
		// save product after cut stock
		tx = cr.mysql.Table("myshop.products").Where("id_product=?", product.IdProduct).Updates(&product)
		if tx.Error != nil {
			return errors.New("can not update prpduct after cut stock")
		}
		cr.mt.Unlock()
	}

	if status == "return" {
		cr.mt.Lock()
		// เช็ค product
		product := models.Product{}
		tx := cr.mysql.Table("myshop.products").Where("id_product=?", idProduct).Find(&product)
		if tx.Error != nil {
			return tx.Error
		}
		// คืนค่า
		product.Quantity = product.Quantity + quantity
		// save product after cut stock
		tx = cr.mysql.Table("myshop.products").Where("id_product=?", product.IdProduct).Updates(&product)
		if tx.Error != nil {
			return errors.New("can not update prpduct after cut stock")
		}
		cr.mt.Unlock()
	}

	return nil
}

// GetCart_Repo implements ICart_Repositiry
// Get All For Store
func (cr *cart_repository) GetAllCartForStore_Repo(idStore int) (carts_Store []models.CartOrderDB, carItemst_Store []models.CartItemDB, err error) {

	// ดึง cart ที่เวลา create เกิน 31 นาที **
	// ปัญหาต่อมา อถ้าข้ามวัน มันจะไม่ดึงข้อมูลขึ้นมา
	tx := cr.mysql.Table("carts").Raw("SELECT * FROM carts WHERE timestampadd(minute, 30,created_at) - Now() < 0  AND id_store = ?", idStore).Scan(&carts_Store)
	if tx.Error != nil {
		return []models.CartOrderDB{}, []models.CartItemDB{}, tx.Error
	}

	// ดึง cartItems ที่เวลา create เกิน 31 นาที **
	tx = cr.mysql.Table("cart_items").Raw("SELECT * FROM cart_items WHERE timestampadd(minute, 30,created_at) - Now() < 0  AND id_store = ?", idStore).Scan(&carItemst_Store)
	if tx.Error != nil {
		return []models.CartOrderDB{}, []models.CartItemDB{}, tx.Error
	}

	return carts_Store, carItemst_Store, tx.Error
}

// GetCart_Repo implements ICart_Repositiry
// Get All For User
func (cr *cart_repository) GetAllCartForUser_Repo(idUser int) (carts_User []models.CartOrderDB, carItemst_User []models.CartItemDB, err error) {
	tx := cr.mysql.Table("myshop.carts").Where("id_user=?", idUser).Find(&carts_User)
	if tx.Error != nil {
		return []models.CartOrderDB{}, []models.CartItemDB{}, tx.Error
	}
	tx = cr.mysql.Table("myshop.cart_items").Where("id_user=?", idUser).Find(&carItemst_User)
	if tx.Error != nil {
		return []models.CartOrderDB{}, []models.CartItemDB{}, tx.Error
	}

	return carts_User, carItemst_User, tx.Error
}

func (cr *cart_repository) GetCartByIdForUser_Repo(idCart, idUser int) (cartItems []models.CartItemDB, err error) {
	tx := cr.mysql.Table("myshop.cart_items").Where("id_user=?", idUser).Where("id_cart=?", idCart).Find(&cartItems)
	return cartItems, tx.Error
}

func (cr *cart_repository) GetCartByIdForStore_Repo(idCart, idStore int) (carItemst_Store []models.CartItemDB, err error) {
	tx := cr.mysql.Table("myshop.cart_items").Where("id_store=?", idStore).Where("id_cart=?", idCart).Find(&carItemst_Store)
	if tx.Error != nil {
		return carItemst_Store, tx.Error
	}
	return carItemst_Store, tx.Error
}

// UpdateCart_Repo implements ICart_Repositiry
// change status for store
func (cr *cart_repository) EditCartOrder_Repo(cartItems []models.CartItemDB) (err error) {
	cartOrder := models.CartOrderDB{}
	// เช้คว่ามี cart id ของเรารึป่าว
	tx := cr.mysql.Table("myshop.carts").Where("id=?", cartItems[0].Id_cart).Where("id_user=?", cartItems[0].Id_User).First(&cartOrder)
	if tx.Error != nil {
		return tx.Error // .First(&carts) ถึงจะได้ error record not found
	}

	//  เช็คเวลา *****
	// เพิ่มเวลา 30 นาทีตรงนี้ สามารถทำให้เช็คว่าเวลาที่ create กับเวลาปัจจุบันมันต่างกันจริง
	createTime := cartOrder.CreatedAt.Add(30 * time.Minute).Unix()

	// ต้องเช็คเวลาก่อน่ว่าการแกไขครั้งนี้ กิน 30 นาทีรึยัง
	if createTime > time.Now().Unix() {
		fmt.Println("In Time")
		// ลบ อันเดิมทิ้ง
		// tx := cr.db.Table("myshop.carts").Where("id=?", cartItems[0].Id_cart).Where("id_user=?", cartItems[0].Id_User).Delete(&cartOrder)

		// if tx.Error != nil {
		// 	return tx.Error
		// }

		// // Create ใหม่
		// tx = cr.db.Table("myshop.cart_items").Create(&cartItems)

		// if tx.Error != nil {
		// 	return tx.Error
		// }

		return nil

	} else {
		// fmt.Println("out time")
		return errors.New("time out for edit order")
	}

}

// Edit status cart  for store *******
//****** ถ้า แกไข ต้องคืนของเข้า store ด้วย ********
func (cr *cart_repository) EditStatusCartOrder_Repo(idCart, idUser int, status *models.StatusCartUpdate) (err error) {
	cartDB := models.CartOrderDB{}

	// เช็คเวลา จะโชว์ order ที่ เกิน 30 นาที ขึ้นไป
	// ให้ user ตัดสินใจ 30 นาทีก่อน
	tx := cr.mysql.Table("carts").Raw("SELECT * FROM carts WHERE timestampadd(minute, 30,created_at) - Now() < 0 AND id_store = ? AND id=?", idUser, idCart).First(&cartDB)
	if tx.Error != nil {
		return tx.Error // .First(&carts) ถึงจะได้ error record not found
	}
	// แก้ status
	cartDB.Status = status.Status

	//save status
	tx = cr.mysql.Table("carts").Save(&cartDB)
	if tx.Error != nil {
		fmt.Println(tx.Error.Error())
		return tx.Error // .First(&carts) ถึงจะได้ error record not found
	}

	return nil
}

func (cr *cart_repository) DeleteCart_Repo(idUser, idcart int) error {
	// get idcart ออกมาก่อนว่ามีไหม
	carts := []models.CartOrderDB{}
	tx := cr.mysql.Table("myshop.carts").Where("id=?", idcart).Where("id_user=?", idUser).First(&carts)
	if tx.Error != nil {
		return tx.Error // .First(&carts) ถึงจะได้ error record not found
	}

	//  update ถ้าไม่มี มันไม่ error ******
	// cartOrder := models.CartOrderDB{}
	// tx = cr.db.Table("myshop.carts").Where("id=?", idcart).Where("id_user=?", idUser).Delete(&cartOrder)

	// cartItems := models.CartItemDB{}
	// tx = cr.db.Table("myshop.cart_items").Where("id_cart=?", idcart).Where("id_user=?", idUser).Find(&cartItems)
	return tx.Error
}
