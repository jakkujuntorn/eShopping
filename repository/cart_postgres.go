package repository

import (
	"fmt"
	"myShopEcommerce/models"
)

func (cr *cart_repository) CreateCart_Repo_Postgres(cartOrder map[int][]models.Product_CartItem, idUser int) error {
	tx := cr.postgres.Begin()

	//************* Sent Data toDB ******************
	for idStore, product_InStore := range cartOrder {
		cartOrderDB := models.CartOrderDB{
			Id_User:  idUser,
			Id_Store: idStore, // คือค่า key value ของ map (cartOrder)
			Status:   "pending",
		}

		if err := tx.Table("carts").Create(&cartOrderDB).Error; err != nil {
			return err
		}

		// ********************แยกข้อมูล CartItemDB ******************
		cartItemDB := []models.CartItem_Postgres{}
		for _, data_CartItem := range product_InStore {

			// ดึง id เตียมมา map กับ catItem
			cartOrderforCartItem := models.CartOrderDB{}
			// ดึงค่าสุดท้ายขึ้นมา มันคือ ค่าที่พึ่งใสลงไป
			// ต้องเอา cart id มา map  กับ cart_items idCart
			if err := tx.Table("carts").Last(&cartOrderforCartItem).Error; err != nil {
				tx.Rollback()
				return err
			}

			// ปั้นข้อมูล CartItemDB
			cartItem := models.CartItem_Postgres{
				ID:         cartOrderforCartItem.Id_Store,
				Id_cart:    int(cartOrderforCartItem.ID), // เอามาจาก Db Cart
				Id_Store:   data_CartItem.Id_Store,
				Id_User:    idUser,
				Id_Product: data_CartItem.Id_Product,
				Quantity:   data_CartItem.Quantity,
				Price:      data_CartItem.Price,
			}
			cartItemDB = append(cartItemDB, cartItem)
		}

		if err := tx.Table("cart_items").Create(&cartItemDB).Error; err != nil {

			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

// for user ให้เวลา 30 นาที
func (cr *cart_repository) EditCartOrder_Repo_Postgres(cartItemDB []models.CartItemDB) (err error) {
	panic("")
}

// for Store หลัง user  crete 30 นาที ถึงจะเข้าไปเปลี่ยน status ได้
func (cr *cart_repository) EditStatusCartOrder_Repo_Postgres(idCart, idUser int, status *models.StatusCartUpdate) (err error) {
	cartDB := models.CartOrderDB{}
	tx := cr.postgres.Table("carts").Raw("select * from carts  where now() > created_at+interval '30 minutes'  and id =? and  id_store=?", idCart, idUser).First(&cartDB)
	if tx.Error != nil {
		return tx.Error
	}
	cartDB.Status = status.Status

	tx = cr.postgres.Table("carts").Where("id=?", idCart).Updates(&cartDB)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (cr *cart_repository) GetAllCartForUser_Repo_Postgres(idUser int) (cartOrderDB []models.CartOrderDB, carItemst_User []models.CartItemDB, err error) {
	tx := cr.postgres.Table("carts").Where("id_user=?", idUser).Find(&cartOrderDB)
	if tx.Error != nil {
		return []models.CartOrderDB{}, []models.CartItemDB{}, tx.Error
	}

	tx = cr.postgres.Table("cart_items").Where("id_user=?", idUser).Find(&carItemst_User)
	if tx.Error != nil {
		return []models.CartOrderDB{}, []models.CartItemDB{}, tx.Error
	}

	return cartOrderDB, carItemst_User, nil
}

func (cr *cart_repository) GetAllCartForStore_Repo_Postgres(idStore int) (cartOrderDB []models.CartOrderDB, cartItemDB_Struct []models.CartItemDB, err error) {

	tx := cr.postgres.Table("carts").Raw("select * from carts  where now() > created_at+interval '30 minutes'  and id_store =?", idStore).Scan(&cartOrderDB)
	if tx.Error != nil {
		return []models.CartOrderDB{}, []models.CartItemDB{}, tx.Error
	}

	tx = cr.postgres.Table("cart_items").Raw("select * from cart_items  where now() > created_at+interval '30 minutes'  and id_store =?", idStore).Scan(&cartItemDB_Struct)
	if tx.Error != nil {
		return []models.CartOrderDB{}, []models.CartItemDB{}, tx.Error
	}

	return cartOrderDB, cartItemDB_Struct, nil
}

// ดึงตาม id_cart กับ  id_user ไม่สนเวลา
func (cr *cart_repository) GetCartByIdForUser_Repo_Postgres(idcart, idUser int) (carItemst_User []models.CartItemDB, err error) {
	// ดึง id cart จาก cart items
	tx := cr.postgres.Table("cart_items").Raw("select * from cart_items  where  id_cart =?  and id_user=?", idcart, idUser).Scan(&carItemst_User)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return carItemst_User, nil
}

func (cr *cart_repository) GetCartByIdForStore_Repo_Postgres(idCart, idStore int) (carItemst_Store []models.CartItemDB, err error) {
	fmt.Println("cart", idCart)
	fmt.Println("store:", idStore)
	tx := cr.postgres.Table("cart_items").Raw("select * from cart_items  where  id_cart =?  and id_store=?", idCart, idStore).Scan(&carItemst_Store)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return carItemst_Store, nil
}
