package repository

import (
	"context"
	"errors"
	"fmt"
	"myShopEcommerce/models"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

// ****************** Mongo *******************
// Create cart
func (cr *cart_repository) CreateCart_Repo_Mongo(cartOrder map[int][]models.Product_CartItem, idUser int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// ********** Transaction ************
	sess, err := cr.monGo.Client().StartSession()
	if err != nil {
		fmt.Println(err)
		return err
	}

	defer sess.EndSession(ctx)

	// err = sess.StartTransaction(options.Transaction().SetWriteConcern(writeconcern.New(writeconcern.WMajority())))
	err = sess.StartTransaction()
	if err != nil {
		fmt.Println("StartTransaction:", err.Error())
		return err
	}

	myShopdatabase := sess.Client().Database("myshop")

	errUseSession := sess.Client().UseSession(ctx, func(sc mongo.SessionContext) error {

		cartOrderDB := models.CartOrderDB_Mongo{
			Id_User:  idUser,
			Id_Store: idUser, // คือค่า key value ของ map (cartOrder)
			Status:   "pending",
			CreateAt: time.Now().Add(7 * time.Hour),
			UpdateAt: time.Now().Add(7 * time.Hour),
		}

		// _, err = cr.monGo.Collection("carts").InsertOne(ctx, cartOrderDB)
		_, err = myShopdatabase.Collection("carts").InsertOne(ctx, cartOrderDB)
		if err != nil {
			fmt.Println("InsertOne:", err.Error())
			// Something went wrong, abort the transaction
			_ = sess.AbortTransaction(ctx)
			return err
		}

		// Simulate an error for testing purposes
		errInTransaction := errors.New("Error Test")
		if errInTransaction != nil {
			fmt.Println("Have Error:", errInTransaction.Error())
			er := sess.AbortTransaction(ctx)
			if er != nil {
				fmt.Println("AbortTransaction :", er)
			}
			return errInTransaction
		}
		return nil
	})

	if errUseSession != nil {
		return errUseSession
	}

	// Everything went well, commit the transaction
	if err := sess.CommitTransaction(ctx); err != nil {
		fmt.Println("CommitTransaction:", err.Error())
		return err
	}

	return nil

	// sess.CommitTransaction(ctx)

	// for idStore, _ := range cartOrder {
	// 	cartOrderDB := models.CartOrderDB_Mongo{
	// 		Id_User:  idUser,
	// 		Id_Store: idStore, // คือค่า key value ของ map (cartOrder)
	// 		Status:   "pending",
	// 		CreateAt: time.Now().Add(7 * time.Hour),
	// 		UpdateAt: time.Now().Add(7 * time.Hour),
	// 	}
	// 	err = sess.StartTransaction()
	// 	if err != nil {
	// 		return err
	// 	}

	// 	//******** insert carts ********
	// 	// resultCarts, err := cr.monGo.Collection("carts").InsertOne(ctx, &cartOrderDB)
	// 	_, errIntranscation = sess.Client().Database("myshop").Collection("carts").InsertOne(ctx, &cartOrderDB)
	// 	// resultCarts, err := cr.monGo.Collection("carts").InsertOne(sess, &cartOrderDB)
	// 	if errIntranscation != nil {
	// 		fmt.Println("carts")
	// 		// sess.AbortTransaction(ctx)
	// 		return errIntranscation
	// 	}

	// 	errIntranscation = errors.New("Test Error")
	// 	// ตรงนี้ต้องเรียกใช้ Abort หรือ  commit
	// 	if errIntranscation != nil {
	// 		// แต่เข้ามาตรงนี้  Transaction ก็ไม่ยกเลิก
	// 		fmt.Println("Have Error in transaction")

	// 		err := sess.AbortTransaction(ctx)
	// 		if err != nil {
	// 			fmt.Println(err.Error())
	// 			return err
	// 		}
	// 		fmt.Println("EndSession")
	// 		 sess.EndSession(ctx)

	// 	} else {
	// 		err = sess.CommitTransaction(ctx)
	// 		if err != nil {
	// 			fmt.Println("CommitTransaction")
	// 			return err
	// 		}
	// 		sess.EndSession(ctx)
	// 	}

	// }

	// *********************************************************
	// transactionFunc := func(sess mongo.SessionContext) (interface{}, error) {

	// 	ss := sess.Client()
	// 	  .StartSession()

	// 	for idStore, product_InStore := range cartOrder {
	// 		cartOrderDB := models.CartOrderDB_Mongo{
	// 			Id_User:  idUser,
	// 			Id_Store: idStore, // คือค่า key value ของ map (cartOrder)
	// 			Status:   "pending",
	// 			CreateAt: time.Now().Add(7 * time.Hour),
	// 			UpdateAt: time.Now().Add(7 * time.Hour),
	// 		}

	// 		//******** insert carts ********
	// 		// resultCarts, err := cr.monGo.Collection("carts").InsertOne(ctx, &cartOrderDB)
	// 		resultCarts, err := ss.Database("myshop").Collection("carts").InsertOne(ctx, &cartOrderDB)
	// 		// resultCarts, err := cr.monGo.Collection("carts").InsertOne(sess, &cartOrderDB)
	// 		if err != nil {
	// 			fmt.Println("carts")

	// 			return nil, err
	// 		}

	// 		// sess.AbortTransaction(ctx)
	// 		// จำลอง error ******************
	// 		// return nil, errors.New("TEst Error")

	// 		// ********************แยกข้อมูล CartItemDB ******************
	// 		cartItemDB := []models.CartItemDB_Mongo{}

	// 		for _, data_CartItem := range product_InStore {

	// 			// การเช็ค stock กับตัด stock
	// 			errIntranscation := cr.CheckStock(data_CartItem.Id_Product, data_CartItem.Quantity, "cut")
	// 			if errIntranscation != nil {
	// 				fmt.Println("CheckStock")

	// 				return nil, errIntranscation
	// 			}

	// 			// ปั้นข้อมูล CartItemDB
	// 			cartItem := models.CartItemDB_Mongo{
	// 				Id_cart:    resultCarts.InsertedID.(primitive.ObjectID).Hex(), // แปลง  _id  เอาแต่ค่า string
	// 				Id_Store:   data_CartItem.Id_Store,
	// 				Id_User:    idUser,
	// 				Id_Product: data_CartItem.Id_Product,
	// 				Quantity:   data_CartItem.Quantity,
	// 				Price:      data_CartItem.Price,
	// 				CreateAt:   time.Now().Add(7 * time.Hour),
	// 				UpdateAt:   time.Now().Add(7 * time.Hour),
	// 			}
	// 			cartItemDB = append(cartItemDB, cartItem)
	// 		}

	// 		// แปลง object
	// 		newVale := make([]interface{}, len(cartItemDB))
	// 		for i, v := range cartItemDB {
	// 			newVale[i] = v
	// 		}

	// 		// resultCartItems, errIntranscation := ssessionCotext.Client().Database("myshop").Collection("cart_items").InsertMany(ssessionCotext, newVale)
	// 		// resultCartItems, errIntranscation := ss.Database("myshop").Collection("cart_items").InsertMany(sess, newVale)
	// 		resultCartItems, errIntranscation := ss.Database("myshop").Collection("cart_items").InsertMany(ctx, newVale)

	// 		if errIntranscation != nil {
	// 			fmt.Println("cart_items")

	// 			return nil, errIntranscation
	// 		}
	// 		_ = resultCartItems

	// 	}

	// 	return nil, nil
	// }

	// result, err := sess.WithTransaction(ctx, transactionFunc)

	// if err != nil {
	// 	fmt.Println("WithTransaction :",err.Error())
	// 	er := sess.AbortTransaction(ctx)
	// 	if er != nil {
	// 		return er
	// 	} else {
	// 		err := sess.CommitTransaction(ctx)
	// 		if err != nil {
	// 			fmt.Println("CommitTransaction")
	// 			return err
	// 		}
	// 	}

	// 	fmt.Println("In This")
	// 	return err
	// }
	// _ = result

	// *********************************  WithTransaction *******************************
	// _, errWithTransaction := sess.WithTransaction(ctx, func(ct mongo.SessionContext) (interface{}, error) {

	// 	// cc := ct.Client()

	// 	// se, _ := cc.StartSession()
	// 	// se.Client().Database("myshop")
	// 	// ***********************************************
	// 	for idStore, product_InStore := range cartOrder {

	// 		cartOrderDB := models.CartOrderDB_Mongo{
	// 			Id_User:  idUser,
	// 			Id_Store: idStore, // คือค่า key value ของ map (cartOrder)
	// 			Status:   "pending",
	// 			CreateAt: time.Now().Add(7 * time.Hour),
	// 			UpdateAt: time.Now().Add(7 * time.Hour),
	// 		}

	// 		//******** insert carts ********
	// 		// resultCarts, errInsertCart := cr.monGo.Collection("carts").InsertOne(ctx, &cartOrderDB)
	// 		resultCarts, errIntranscation := ct.Client().Database("myshop").Collection("carts").InsertOne(ctx, &cartOrderDB)
	// 		if errIntranscation != nil {
	// 			return nil, errIntranscation
	// 		}

	// 		// ********************แยกข้อมูล CartItemDB ******************
	// 		cartItemDB := []models.CartItemDB_Mongo{}

	// 		for _, data_CartItem := range product_InStore {

	// 			// การเช็ค stock กับตัด stock
	// 			errIntranscation = cr.CheckStock(data_CartItem.Id_Product, data_CartItem.Quantity, "cut")
	// 			if errIntranscation != nil {
	// 				return nil, errIntranscation
	// 			}

	// 			// ปั้นข้อมูล CartItemDB
	// 			cartItem := models.CartItemDB_Mongo{
	// 				Id_cart:    resultCarts.InsertedID.(primitive.ObjectID).Hex(), // แปลง  _id  เอาแต่ค่า string
	// 				Id_Store:   data_CartItem.Id_Store,
	// 				Id_User:    idUser,
	// 				Id_Product: data_CartItem.Id_Product,
	// 				Quantity:   data_CartItem.Quantity,
	// 				Price:      data_CartItem.Price,
	// 				CreateAt:   time.Now().Add(7 * time.Hour),
	// 				UpdateAt:   time.Now().Add(7 * time.Hour),
	// 			}
	// 			cartItemDB = append(cartItemDB, cartItem)
	// 		}

	// 		// แปลง object
	// 		newVale := make([]interface{}, len(cartItemDB))
	// 		for i, v := range cartItemDB {
	// 			newVale[i] = v
	// 		}

	// 		//*********** to MongoDB ********
	// 		// ctxx,_:=context.WithTimeout(context.Background(), 1*time.Second)
	// 		// time.Sleep(5 *time.Second)

	// 		// resultCartItems, errInsertCartItems := cr.monGo.Collection("cart_items").InsertMany(ctx, newVale)
	// 		resultCartItems, errIntranscation := ct.Client().Database("myshop").Collection("cart_items").InsertMany(ctx, newVale)
	// 		// fmt.Println(errInsertCartItems)

	// 		if errIntranscation != nil {
	// 			return nil, errIntranscation
	// 		}

	// 		_ = resultCartItems

	// 	}

	// 	// if errIntranscation != nil {
	// 	// 	// แต่เข้ามาตรงนี้  Transaction ก็ไม่ยกเลิก
	// 	// 	fmt.Println("Have Error in transaction")
	// 	// 	err := sess.AbortTransaction(ctx)
	// 	// 	if err != nil {
	// 	// 		fmt.Println(err.Error())
	// 	// 		return nil, err
	// 	// 	}
	// 	// 	fmt.Println("EndSession")
	// 	// 	sess.EndSession(ctx)
	// 	// 	return nil, nil

	// 	// } else {
	// 	// 	err = sess.CommitTransaction(ctx)
	// 	// 	if err != nil {
	// 	// 		fmt.Println("CommitTransaction")
	// 	// 		return nil, err
	// 	// 	}
	// 	// 	sess.EndSession(ctx)
	// 	// 	return nil, nil
	// 	// }

	// 	return nil, nil
	// })

	// if errWithTransaction != nil {
	// 	err := sess.AbortTransaction(ctx)
	// 	if err != nil {
	// 		fmt.Println(err.Error())
	// 		return err
	// 	}
	// 	fmt.Println(errWithTransaction.Error())
	// 	return errWithTransaction
	// }
	//***********  End WithTransaction ************

	return nil
}

//**************** For User ***********
func (cr *cart_repository) EditCartOrder_Repo_Mongo(cartItemDB []models.CartItemDB) (err error) {
	return nil
}

//*************** For Store **************
func (cr *cart_repository) EditStatusCartOrder_Repo_Mongo(idCart, idUser int, status *models.StatusCartUpdate) (err error) {
	return nil
}
