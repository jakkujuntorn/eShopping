package repository

import (
	"errors"
	"fmt"

	"myShopEcommerce/models"
)

func (cr *cart_repository) ReStoreStock_Postgres(idProduct, quantity int) error {
	cr.mt.Lock()
	// เช็ค product ******************
	product := models.Product{}
	tx := cr.mysql.Table("myshop.products").Where("id_product=?", idProduct).Find(&product)
	if tx.Error != nil {
		return tx.Error
	}

	// คืนค่า ******************
	myQuery := `UPDATE  myshop.products SET  quantity= quantity+? WHERE  id_product = ?`
	tx = cr.mysql.Table("myshop.products").Exec(myQuery, quantity, idProduct)
	if tx.Error != nil {
		return tx.Error // .First(&carts) ถึงจะได้ error record not found
	}

	// product.Quantity = product.Quantity + quantity
	// // save product after cut stock
	// tx = cr.mysql.Table("myshop.products").Where("id_product=?", product.IdProduct).Updates(&product)
	// if tx.Error != nil {
	// 	return errors.New("can not update prpduct after cut stock")
	// }
	cr.mt.Unlock()
	return nil
}
func (cr *cart_repository) CheckStock_Postgres(idProduct, quantity int, status string) error {
	// txr := transation
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
	return nil
}

// Transaction for Postgres **********
// ตรงนี้ที่ต้องใช้ struct เพราะ ต้องใช้ db
func (cr *cart_repository) Transaction_Postgres(fn func(*cart_repository) error) error {
	fmt.Println("")
	dbPostgres := cr.postgres.Begin()
	dbMYSQL := cr.mysql.Begin()

	nRepo := NewCart_Repo(dbMYSQL, dbPostgres)

	// func ต้องได้ struct ถึงจะเรียกใช้ recive Func ของมันได้
	// ภายใน func นี้ถ้ามีการเรียกใช้ query แล้วมี  error ออกมามันจะไม่ commit ให้
	// transcation ทั้หมดเลยถูกยกเลิก
	err := fn(nRepo)

	// คีร์หลักคือตรงนี้ ถ้าใน การทำ transcation มี error มันจะมาหยุดตรงนี้ และไม่มีถึ
	// การทำงานของ commit
	if err != nil {
		fmt.Println("Error in TRacscation")
		dbPostgres.Rollback()
		dbMYSQL.Rollback()
		return err
	}

	dbPostgres.Commit()
	dbMYSQL.Commit()
	return nil
}

//    ****************************     V1   *******************************
// func (cr *cart_repository) CreateCart_Repo_Postgres(cartOrder map[int][]models.Product_CartItem, idUser int) error {
// 	tx := cr.postgres.Begin()

// 	//************* Sent Data toDB ******************
// 	for idStore, product_InStore := range cartOrder {
// 		cartOrderDB := models.CartOrderDB{
// 			Id_User:  idUser,
// 			Id_Store: idStore, // คือค่า key value ของ map (cartOrder)
// 			Status:   "pending",
// 		}

// 		if err := tx.Table("carts").Create(&cartOrderDB).Error; err != nil {
// 			return err
// 		}

// 		// ********************แยกข้อมูล CartItemDB ******************
// 		cartItemDB := []models.CartItem_Postgres{}
// 		for _, data_CartItem := range product_InStore {

// 			// ดึง id เตียมมา map กับ catItem
// 			cartOrderforCartItem := models.CartOrderDB{}
// 			// ดึงค่าสุดท้ายขึ้นมา มันคือ ค่าที่พึ่งใสลงไป
// 			// ต้องเอา cart id มา map  กับ cart_items idCart
// 			if err := tx.Table("carts").Last(&cartOrderforCartItem).Error; err != nil {
// 				tx.Rollback()
// 				return err
// 			}

// 			// ปั้นข้อมูล CartItemDB
// 			cartItem := models.CartItem_Postgres{
// 				ID:         cartOrderforCartItem.Id_Store,
// 				Id_cart:    int(cartOrderforCartItem.ID), // เอามาจาก Db Cart
// 				Id_Store:   data_CartItem.Id_Store,
// 				Id_User:    idUser,
// 				Id_Product: data_CartItem.Id_Product,
// 				Quantity:   data_CartItem.Quantity,
// 				Price:      data_CartItem.Price,
// 			}
// 			cartItemDB = append(cartItemDB, cartItem)
// 		}

// 		if err := tx.Table("cart_items").Create(&cartItemDB).Error; err != nil {

// 			tx.Rollback()
// 			return err
// 		}
// 	}
// 	tx.Commit()
// 	return nil
// }

// ********************* Delete Cart for User * ******************
func (cr *cart_repository) DeleteCartForUser_Repo_Postgres(idUser, idCart int) error {
	deleteCart := []models.Delete_GetCart_Postgres{}
	deleteCart_items := []models.Delete_GetCartItem_Postgres{}

	errTransaction := cr.Transaction_Postgres(func(d *cart_repository) error {
		// ********************** ดึงข้อมูล ที่จะลบ ******************************
		// เช็ค cart ว่าเวลาไม่เกิน 30 นาที  ถ้าเกิน 30 นาที ห้ามลบ
		// ดึงค่าใน carts *******************
		tx := d.postgres.Table("carts").Raw("select * from carts  where now() < created_at+interval '30 minutes'  and id =? and id_user=? and deleted_at IS NULL", idCart, idUser).Find(&deleteCart)
		if tx.Error != nil {
			// return tx.Error
			fmt.Println(tx.Error.Error())
			return errors.New("can not delete cart time out or record not found")
		}

		// เช็คว่า id cart ที่จะลบหรือไม่
		// เช็คที่ carts ที่เดียวก็ได้เพราะ ถ้าใน carts มี ใน cart_items  มันก็ต้องมี
		if len(deleteCart) == 0 {
			return errors.New("Can not Delete record not found")
		}

		fmt.Println("Start Get cart_items")
		// ดึงค่าใน cart_items **************************
		tx = d.postgres.Table("cart_items").Raw("select * from cart_items  where now() < created_at+interval '30 minutes'  and id_cart =? and id_user=?", deleteCart[0].ID, deleteCart[0].Id_User).Scan(&deleteCart_items)
		if tx.Error != nil {
			// return tx.Error
			fmt.Println(tx.Error.Error())
			return errors.New("can not delete cart_items time out")
		}

		// *************************** คืนค่าให้ product  Mysql *********************
		// คืนค่า stock ให้ product ที่ลบ ใน mysql
		// ตรงนี้ for เพื่อคืนค่าที่ละ product
		// จะทำ transcation ยังไง

		// myQuery := `UPDATE  myshop.products SET  quantity= quantity+? WHERE  id_product = ?`

		for _, v := range deleteCart_items {
			// อันนี้ไม่ใช้ เอาไปสร้าง func แยกออกไปแล้ว
			// tx := cr.mysql.Table("myshop.products").Exec(myQuery, v.Quantity, v.Id_Product)
			// if tx.Error != nil {
			// 	return tx.Error // .First(&carts) ถึงจะได้ error record not found
			// }

			err := d.ReStoreStock_Postgres(v.Id_Product, v.Quantity)
			if err != nil {
				return err
			}
		}

		//************************* ตรงนี้ที่ทำให้มันสร้าง เพิ่ม *************
		// ปัญหาที่เจอ มันสร้าง 2 table
		// เพราะ  ข้อมูลใน cloumn  ทั้ง 22 อันมันเหมือนกัน gorm แยกไม่ได้ มันเลยสร้าง 2 table
		// วิธีแก้ สร้าง table ใหม่ ที่ชื่อใน cloumn ต่างกัน

		for _, v := range deleteCart {
			sqlQuery := `INSERT INTO deleted_carts (id,"id_userDel","idstoreDel","statusDel",created_at, udpated_at) VALUES (?,?,?,?,?,? )`
			txx := d.postgres.Table("deleted_carts").Exec(sqlQuery, v.ID, v.Id_User, v.Id_Store, v.Status, v.CreatedAt, v.UpdatedAt)
			if txx.Error != nil {
				return txx.Error
			}
		}

		// บันทึกลงใน delete_cart_items
		for _, v := range deleteCart_items {

			// d := models.Delete_cart_items{
			// 	ID:            v.ID,
			// 	Id_cartDel:    v.Id_cart,
			// 	Id_StoreDel:   v.Id_Store,
			// 	Id_UserDel:    v.Id_User,
			// 	Id_ProductDEl: v.Id_Product,
			// 	QuantityDel:   v.Quantity,
			// 	PriceDel:      v.Price,
			// }
			// _ = d

			sqlQuery1 := `INSERT INTO public.deleted_cart_items(id, "id_cartDel", "id_userDel", "id_storeDel", "id_productDel", "quantityDel", "priceDel ",created_at, updated_at)VALUES (?, ?, ?, ?, ?, ?, ?,?,?)`
			// sqlQuery2 := `INSERT INTO deleted_cart_items (id,"id_cartDel","id_userDel,"id_storeDel","id_productDel","quantityDel","priceDel",) VALUES (?,?,?,?,?,?,?)`
			// _= sqlQuery2
			// txx := cr.postgres.Table("deleted_cart_items").Exec(sqlQuery1, &d)

			txx := d.postgres.Table("deleted_cart_items").Exec(sqlQuery1, v.ID, v.Id_cart, v.Id_User, v.Id_Store, v.Id_Product, v.Quantity, v.Price, v.CreatedAt, v.UpdatedAt)
			if txx.Error != nil {
				return txx.Error
			}
		}

		// *************** ใช้ท่านี้ไม่ได้ มัน เรียง ตัวแปรผิด ***************
		// tx = cr.postgres.Table("delete_carts").Create(&d_CreateCarts)
		// if tx.Error != nil {
		// 	return tx.Error
		// }
		// บันทึกลงใน delete_cart_items
		// tx = cr.postgres.Table("delete_cart_items").Create(&d_Create_Items)
		// if tx.Error != nil {
		// 	return tx.Error
		// }

		//  *******************ลบออกจาก carts จริงๆ จาก postgres ************************

		tx = d.postgres.Table("carts").Delete(&models.Delete_GetCart_Postgres{}, "id=?", idCart)
		if tx.Error != nil {
			return tx.Error
		}

		tx = d.postgres.Table("cart_items").Delete(&models.Delete_GetCartItem_Postgres{}, "id_cart=?", idCart)
		if tx.Error != nil {
			return tx.Error
		}

		// กรณีไม่มี error เลย
		return nil
	})

	// error ตรงนี้จะได้ค่ามาจาก การทำงานใน Transaction_Postgres()
	if errTransaction != nil {
		fmt.Println("errTransaction:", errTransaction.Error())
		return errTransaction
	}

	return nil
}

//*********************   V2     ****************
func (cr *cart_repository) CreateCart_Repo_Postgres(cartOrder map[int][]models.Product_CartItem, idUser int) error {

	//************* Sent Data toDB ******************
	// func(d *cart_repository) ต้องรับ  struct ถึงจะได้ recive fun มาใช้งาน
	errTransaction := cr.Transaction_Postgres(func(d *cart_repository) error {
		// d ต้องได้ interface ของ repository ด้วย
		//จะใช้ งาน func อะไรให้ใช้จาก d *cart_repository

		//******************** สร้าง cart ID ***********************
		for idStore, product_InStore := range cartOrder {
			fmt.Println("User", idUser)
			fmt.Println("Store:", idStore)
			cartOrderDB := models.CartOrderDB{
				Id_User:  idUser,
				Id_Store: idStore, // คือค่า key value ของ map (cartOrder)
				Status:   "pending",
			}
			// สร้าง cart ID ***********************
			if err := d.postgres.Table("carts").Create(&cartOrderDB).Error; err != nil {
				return err
			}

			// ********************แยกข้อมูล CartItemDB ******************
			cartItemDB := []models.CartItem_Postgres{}
			// ****************สร้าง cart_items ****************
			for _, data_CartItem := range product_InStore {

				//******* Check Stoc and cut stock *******************
				// เช็คดูว่า มันทำ transcation จากตรงนี้ หรือใน Func Transaction_Postgres หรือทั้ง 2 จุด
				// ตรงนี้ก็ใช้  d ที่มาจาก d *cart_repository มันเลยมองว่า อยู่ใน transcation เดียวกัน
				errStock := d.CheckStock_Postgres(data_CartItem.Id_Product, data_CartItem.Quantity, "cut")
				if errStock != nil {
					return errStock
				}

				// ดึง id เตียมมา map กับ catItem
				cartOrderforCartItem := models.CartOrderDB{}
				// ดึงค่าสุดท้ายขึ้นมา มันคือ ค่าที่พึ่งใสลงไป
				// ต้องเอา cart id มา map  กับ cart_items idCart
				if err := d.postgres.Table("carts").Last(&cartOrderforCartItem).Error; err != nil {
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
			//สร้าง cart_items ****************
			if err := d.postgres.Table("cart_items").Create(&cartItemDB).Error; err != nil {
				return err
			}
		}
		// กรณีไม่มี error เลย
		return nil
	})

	// error ตรงนี้จะได้ค่ามาจาก การทำงานใน Transaction_Postgres()
	if errTransaction != nil {
		fmt.Println("errTransaction:", errTransaction.Error())
		return errTransaction
	}

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
	// fmt.Println("cart", idCart)
	// fmt.Println("store:", idStore)
	tx := cr.postgres.Table("cart_items").Raw("select * from cart_items  where  id_cart =?  and id_store=?", idCart, idStore).Scan(&carItemst_Store)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return carItemst_Store, nil
}
