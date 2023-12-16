package service

import (
	"errors"
	"fmt"
	"myShopEcommerce/models"
	"myShopEcommerce/repository"
	"myShopEcommerce/util"
	_ "sort"
)

type cart_service struct {
	cartRepo repository.ICart_Repositiry
}

type ICart_Service interface {
	CreateCart_Service(map[int][]models.Product_CartItem, int) error
	CreateCart_Service_V2(map[int][]models.Product_CartItem, int) error

	GetCartById_Service(idParam, idUser int, path string) (cartItemDB []models.CartItemDB, err error)
	GetCartByUserType_Service(idUser int, path string) (CartOrderResponse *models.CartOrderResponse, err error)

	EditCartOrder_Service(idCart, idUser int, newData *models.CartRequest) error
	EdiStatustCart_Service(idCart, idUser int, status *models.StatusCartUpdate) error

	DeleteCartForUser_Service(idUser, idCart int) error
	ValidateData_serviice(dataValidate interface{}) error
}

func NewCart_Service(cartRepo repository.ICart_Repositiry) ICart_Service {
	return &cart_service{
		cartRepo: cartRepo,
	}
}

func (cs *cart_service) ValidateData_serviice(dataValidate interface{}) error {
	errValidate := util.ValidateDataUser(dataValidate)
	return errValidate
}

// CreateCart_Service implements ICart_Service
// ทำไม parmiter ถึงเป็น map[int] ได้ เพราะ  ข้อมูลส่ง มาทั้งหมดรวมทั้ง key value ด้วย
// Funcนี้จะเข้าไปอ่าน DB  2 ครั้ง
// 1.สร้าง id cart ก่อน
// 2.เอา id cart มาสร้าง
func (cs *cart_service) CreateCart_Service(cartOrder map[int][]models.Product_CartItem, idUser int) error {
	// cartOrder 1 ร้าน 1 บิล
	// fmt.Println(cartOrder) // ข้มูลมาเป็น array แยกตาม idStore  แล้ว
	// output map[20:[{0 20 21 5 500 } {0 20 27 10 200 }] 100:[{0 100 108 1 7500 }]]

	// แยกข้อมูล FormatCartData ออกเป็น cartOrderDB และ CartItemsDB
	// **********  cartOrderDB ***************
	// cartOrderDB := models.CartOrderDB{}
	//********** แบบ 1 *********************************
	// for _, data_Order := range cartOrder {
	// 	fmt.Print("")
	// 	// fmt.Println(data) // {12 5 pending [{0 5 77 1 500 }]}
	// 	cartOrderDB = models.CartOrderDB{
	// 		Id_User:  data_Order.Id_User,
	// 		Id_Store: data_Order.Id_Store,
	// 		Status:   "pending",
	// 	}

	// แยกข้อมูล เป็น 2 struct CartOrderDB กับ CartItemDB
	//****************** แบบ 2  ลด loop ไป 1 loop ***********************

	//************* แยกข้อมูล CartOrderDB *****************
	// loop ตามจำนวน Store
	// fmt.Print(cartOrder) //map[20:[{20 21 5 500} {20 27 10 200}] 100:[{100 108 1 7500}]]
	//เช่น idStore 20 กับ 100 จะ loop 2 รอบ
	for idStore, product_InStore := range cartOrder {
		fmt.Print("")
		// fmt.Println(data) // {12 5 pending [{0 5 77 1 500 }]}
		//*********** ปั้นข้อมูล cartOrderDB ********************
		cartOrderDB := models.CartOrderDB{
			Id_User:  idUser,
			Id_Store: idStore, // คือค่า key value ของ map (cartOrder)
			Status:   "pending",
		}
		// fmt.Println(cartOrderDB)
		// create cart DB repository และส่ง id cart ออกมา เพื่อเอามา map กับ cartItems
		// จะมี Cart ตามจำนวนร้านค้า เช่น มี Store 20 กับ 100 จะมี  Cart 2 ฟิว

		// idCart := make(chan int)
		// go func() error {
		// 	fmt.Println(" ************** start Cart ********************************")
		// 	orderCartDB, err := cs.cartRepo.CreateCart_Repo(cartOrderDB)
		// 	if err != nil {
		// 		return err
		// 	}
		// 	fmt.Println(" start id ///////////////////////")
		// 	idCart <- int(orderCartDB.ID)
		// 	return nil
		// }()

		// select {
		// case magError := <-errCart:
		// 	return errors.New(magError)
		// }

		orderCartDB, errCart := cs.cartRepo.CreateCart_Repo(cartOrderDB)
		if errCart != nil {
			return errCart
		}
		_ = orderCartDB

		// ********************แยกข้อมูล CartItemDB ******************
		cartItemDB := []models.CartItemDB{}
		// loop ตาม สินค้า ของ store นั้นๆ เพื่อ map cart id ลงไป
		// เช่น map[20:[{20 21 5 500} {20 27 10 200}] 100:[{100 108 1 7500}]]
		// store 20 มีสินค้า 2 ชิ้น จะ loop 2 รอบ
		for _, data_CartItem := range product_InStore {
			// ปั้นข้อมูล CartItemDB
			cartItem := models.CartItemDB{
				Id_cart: int(orderCartDB.ID), // เอามาจาก Db Cart
				// Id_cart:    int(<-idCart), // เอามาจาก Db Cart
				Id_Store:   data_CartItem.Id_Store,
				Id_User:    idUser,
				Id_Product: data_CartItem.Id_Product,
				Quantity:   data_CartItem.Quantity,
				Price:      data_CartItem.Price,
			}
			cartItemDB = append(cartItemDB, cartItem)
		}

		// fmt.Println(cartItemDB)
		// cartItems DB / 1 product จะมี 1 ฟิว
		// เช่น store 20 มี สอนค้า 2 ชิ้น DB จะสร้าง 2 ฟิว
		// cartItemDB เป็น array ส่งเข้าไป DB มัน create ได้

		errCartItem := cs.cartRepo.CreateCartItem_Repo(&cartItemDB)
		if errCartItem != nil {
			// return errCartItem
			return errCartItem
		}

	}

	return nil
}

func (cs *cart_service) CreateCart_Service_V2(cartOrder map[int][]models.Product_CartItem, idUser int) error {
	//************ mySQL *************
	// err := cs.cartRepo.CreateCart_Repo_V2(cartOrder, idUser)

	// *************  Mongo ***************
	// err := cs.cartRepo.CreateCart_Repo_Mongo(cartOrder, idUser)
	// if err != nil {
	// 	return err
	// }

	// ********************  POstgres  **************
	err := cs.cartRepo.CreateCart_Repo_Postgres(cartOrder, idUser)
	if err != nil {
		return err
	}
	fmt.Println("Transaction Completed")
	return nil
}

func (cs *cart_service) GetCartById_Service(idParam, idUser int, path string) (cartItemDB []models.CartItemDB, err error) {

	if path == "cartstorebyid" {
		// ***********************  mysql  *************************************
		// cart_ItemStore, errCartRepo := cs.cartRepo.GetCartByIdForStore_Repo(idParam, idUser)
		// fmt.Println(len(cart_ItemStore))

		// *****************************  postgres  ***********************
		cart_ItemStore, errCartRepo := cs.cartRepo.GetCartByIdForStore_Repo_Postgres(idParam, idUser)

		if errCartRepo != nil {
			return []models.CartItemDB{}, errCartRepo
		}

		if len(cart_ItemStore) == 0 {
			return []models.CartItemDB{}, errors.New("record not found")
		}

		return cart_ItemStore, nil
	}

	if path == "cartuserbyid" {
		// ***********************   mysql  ****************************************
		// cart_ItemUser, errCartOrder := cs.cartRepo.GetCartByIdForUser_Repo(idParam, idUser)
		// if errCartOrder != nil {
		// 	return []models.CartItemDB{}, errCartOrder
		// }

		// ***************************  postgres  *************************
		cart_ItemUser, errCartOrder := cs.cartRepo.GetCartByIdForUser_Repo_Postgres(idParam, idUser)
		if errCartOrder != nil {
			return []models.CartItemDB{}, errCartOrder
		}

		if len(cart_ItemUser) == 0 {
			return []models.CartItemDB{}, errors.New("record not found")
		}

		return cart_ItemUser, nil
	}

	return []models.CartItemDB{}, nil
}

// GetCartById_Service implements ICart_Service
func (cs *cart_service) GetCartByUserType_Service(idUser int, path string) (*models.CartOrderResponse, error) {

	catrResponse := models.CartOrderResponse{}

	if path == "cartstore" {
		//  ******************  mysql  ***********************
		// carts, cart_items, errCartRepo := cs.cartRepo.GetAllCartForStore_Repo(idUser)
		// if errCartRepo != nil {
		// 	return &catrResponse, errCartRepo
		// }

		// ************************  postgres  ***********************
		carts, cart_items, errCartRepo := cs.cartRepo.GetAllCartForStore_Repo_Postgres(idUser)
		if errCartRepo != nil {
			return &catrResponse, errCartRepo
		}

		groupCart_Items := make(map[int][]models.CartItemDB)
		// จัดกลุ่ม  cart_items
		for _, data := range cart_items {
			groupCart_Items[int(data.Id_cart)] = append(groupCart_Items[int(data.Id_cart)], data)
		}

		// ปั้นข้อมูลให่ ก่อน response
		catrResponse.Id_User = idUser
		// map carts กับ cart_items
		for _, data := range carts {
			// หาราคารวมของสินค้าทั้งหมด
			totalPrice := util.TotalPrice(groupCart_Items[int(data.ID)])
			// ปั้นข้อมูลใหม่สำหรับ CartOrderResponse_Store.Order
			OrderCare_Store := models.OrderCartDetail{
				Id_cart: int(data.ID),
				Product: groupCart_Items[int(data.ID)],
				Total:   totalPrice,
				Status:  data.Status,
			}
			catrResponse.Cart = append(catrResponse.Cart, OrderCare_Store)
		}

		return &catrResponse, nil
	}

	if path == "cartuser" {

		// **********************   mysql  ***********************
		// cartOrderDB, cart_ItemUser, errCartOrder := cs.cartRepo.GetAllCartForUser_Repo(idUser)
		// if errCartOrder != nil {
		// 	return &catrResponse, errCartOrder
		// }

		// **************************  postgres  ******************************
		cartOrderDB, cart_ItemUser, errCartOrder := cs.cartRepo.GetAllCartForUser_Repo_Postgres(idUser)
		if errCartOrder != nil {
			return &catrResponse, errCartOrder
		}

		cartItems := make(map[int][]models.CartItemDB)
		// 1.1 ************* แบ่งกบลุ่ม cart Items ตาม id_cart **************
		for _, cartItem := range cart_ItemUser {
			cartItems[int(cartItem.Id_cart)] = append(cartItems[int(cartItem.Id_cart)], cartItem) //  ยังไม่เข้าใจการทำงานตรงนี้
		}

		catrResponse.Id_User = idUser
		for _, cartOrder := range cartOrderDB {
			totalPrice := util.TotalPrice(cartItems[int(cartOrder.ID)])

			// ปั้นข้อมูล orderCartDeteil
			// OrderCartDetail เป็น array  เพราะ 1 user มีหลาย order
			orderCartDeteil := models.OrderCartDetail{
				Id_cart: int(cartOrder.ID),
				Product: cartItems[int(cartOrder.ID)], // cartItems ถูกจัดกลุ่มมาแล้ว แค่เอาข้อมูลตาม key ที่ myshop.carts.ID
				Total:   totalPrice,
				Status:  cartOrder.Status,
			}

			catrResponse.Cart = append(catrResponse.Cart, orderCartDeteil)
		}

		return &catrResponse, nil
	}

	return &models.CartOrderResponse{}, nil
}

// Edit cart order  V2
func (cs *cart_service) EditCartOrder_Service(idUser, idCart int, newData *models.CartRequest) error {
	//ปั้นข้อมูลใหม่
	cartItems := []models.CartItemDB{}

	for _, v := range newData.Product {
		cartItem := models.CartItemDB{
			Id_cart:  idCart,
			Id_Store: v.Id_Store,
			Id_User:  idUser,
			Quantity: v.Quantity,
			Price:    v.Price,
		}
		cartItems = append(cartItems, cartItem)
	}

	// ส่งเข้า db
	//  ทำไมตรงนี้ใช้ pointer ไม่ได้
	err := cs.cartRepo.EditCartOrder_Repo(cartItems)
	return err
}

func (cs *cart_service) EdiStatustCart_Service(idCart, idUser int, status *models.StatusCartUpdate) error {
	// **************************  mysql  ***************************
	// return cs.cartRepo.EditStatusCartOrder_Repo(idCart, idUser, status)

	// ****************************  postgres  ***************************
	return cs.cartRepo.EditStatusCartOrder_Repo_Postgres(idCart, idUser, status)
}

func (cs *cart_service) DeleteCartForUser_Service(idUser, idCart int) error {
	// err := cs.cartRepo.DeleteCart_Repo(idUser, idCart)
	// if err != nil {
	// 	// fmt.Println(err.Error())
	// 	return err
	// }

	// postgres **********************
	err := cs.cartRepo.DeleteCartForUser_Repo_Postgres(idUser, idCart)
	if err != nil {

		return err
	}

	return nil
}
