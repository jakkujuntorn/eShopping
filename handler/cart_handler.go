package handler

import (
	"fmt"
	"myShopEcommerce/models"
	"myShopEcommerce/service"
	"myShopEcommerce/util"
	"strings"

	"github.com/gin-gonic/gin"
)

type cart_handler struct {
	cartService service.ICart_Service
}

type ICart_Handler interface {
	CreateCart_Handler(*gin.Context)
	EditCartOrder_Handler(*gin.Context)
	GetCartByUserTYpe_Handler(*gin.Context)
	GetCartById_Handler(*gin.Context)
	EditStatusCart_Handler(*gin.Context)
	DeleteCart_Handler(*gin.Context)
}

func NewCart_Handler(cartService service.ICart_Service) ICart_Handler {
	return &cart_handler{
		cartService: cartService,
	}
}

// CreateCart_Handler implements ICart_Handler
func (ch *cart_handler) CreateCart_Handler(c *gin.Context) {

	cartRequest := models.CartRequest{}

	//***** ดึงค่า id user จาก token *******
	idUser, errGetId := util.GetidToken(c)
	if errGetId != nil {
		c.JSON(401, util.Error_Custom(401, errGetId.Error(), "handler_cart_createcart_getidtoken"))
		return
	}

	// ShouldBind
	errShould := c.ShouldBindJSON(&cartRequest)
	if errShould != nil {
		c.JSON(400, util.Error_Custom(400, errShould.Error(), "handler_cart_createcart_ShouldBindJSON"))
		return
	}

	// Validati data cart
	errValidate := ch.cartService.ValidateData_serviice(&cartRequest)
	if errValidate != nil {
		c.JSON(400, util.Error_Custom(400, errValidate.Error(), "handler_cart_createcart_validator"))
		return
	}

	// ปั้น data ใหม่ โดยแยก idStore  สินค้าร้านตัวเองเท่านั้น  ที่จะอยู่ร่วมกัน
	groupByIdStore := make(map[int][]models.Product_CartItem) // is array แต่ แบบนี้จะมี index เป็นค่าที่เรากำหนดได้

	// ************  แยกกลุ่ม product  ตาม idStore *********
	for _, data := range cartRequest.Product {
		// เอาค่า index เดิมใสค่าเพิ่มลงไป
		// ค่า idStore เดิมจะใส ข้อมูลที่  idStore  เท่ากันลงไป
		groupByIdStore[data.Id_Store] = append(groupByIdStore[data.Id_Store], data)
	}

	// fmt.Println(groupIdStore)
	// idStore 20 มีสิ้นค้า 2 ชิ้น idStore 100 มีสินค้า 1 ชิ้น
	// output map[20:[{0 20 21 5 500 } {0 20 27 10 200 }] 100:[{0 100 108 1 7500 }]]

	// =============== ปั้นอีก แบบ =================
	// ข้อมูลไม่ถูกเพราะ idStore ค่าเดียวกัน  แต่ค่าอื่นไม่เหมือนกัน เลยมองว่าข้อมูลคนละชุด
	// groupIdStore2 := make(map[interface{}][]models.CartItem)
	// for _, data := range cartRequest.Product {
	// 	groupIdStore2[data] = append(groupIdStore2[data], data)
	// }

	// ข้อมูลมี  3ชุด เพราะอะไร
	// map[
	//	{20 21 5 500 }:[{20 21 5 500 }]  {20 27 10 200 }:[{20 27 10 200 }]
	//  {100 108 1 7500 }:[{100 108 1 7500 }]
	//  ]
	// fmt.Println(groupIdStore2)

	//******************** อันนี้ไม่ต้องใช้แล้ว ********************
	// cartOrders := []models.FormatCartData{}
	// //********* ปั้นข้อมูลใหม่ เพื่อใส idUser กับ status *******
	// // อาจไม่ต้องทำ idUser ส่งแยก
	// // status ให้ DB ทำค่า default
	// for idStore, cartItem := range groupIdStore {
	// 	cartOrder := models.FormatCartData{
	// 		Id_User:   idUser, // มาจาก token
	// 		Id_Store:  idStore,
	// 		// Status:    "pending",
	// 		CartItems: cartItem,
	// 	}
	// 	cartOrders = append(cartOrders, cartOrder)
	// }
	// fmt.Println(cartOrders)
	// {id_user, idStore,ststus(pending) [{idCart, idStore, idProduct, quantity, price }
	//ouput [{12 20 pending [{0 20 21 5 500 } {0 20 27 10 200 }]} / {12 100 pending [{0 100 108 1 7500 }]}]

	//************************ DB ****************************
	// groupByIdStore คือ product ที่แยกกลุ่มไว้แล้ว
	errCreateCArt := ch.cartService.CreateCart_Service_V2(groupByIdStore, idUser)
	if errCreateCArt != nil {
		c.JSON(500, util.Error_Custom(500, errCreateCArt.Error(), "handler_cart_createcart_createcart_service"))
		return
	}

	c.JSON(201, gin.H{
		"Success": "true",
		"Message": "Create Cart Success",
	})
}

// GetCartByUserTYpe_Handler implements ICart_Handler
func (ch *cart_handler) GetCartByUserTYpe_Handler(c *gin.Context) {
	//***** ดึงค่า id user จาก token *******
	idUserType, errGetId := util.GetidToken(c)
	if errGetId != nil {
		c.JSON(401, util.Error_Custom(401, errGetId.Error(), "handler_cart_getallcartforstore_getidtoken"))
		return
	}

	// ดึง path เพื่อแยก type
	//ว่าจะเป็น user เพื่อเช็ค ว่าตัวเองสั่งอะไร
	// หรือ store ว่ามีใครสั่งอะไรมา
	fullpath := c.FullPath()
	fmt.Println(fullpath)
	path := strings.Split(fullpath, "/")

	// to service
	cartOrder_Store, err := ch.cartService.GetCartByUserType_Service(idUserType, path[2])
	if err != nil {
		c.JSON(500, util.Error_Custom(500, err.Error(), "handler_cart_getallcartforstore_getallcartforstore_service"))
		return
	}

	c.JSON(200, gin.H{
		"Data": cartOrder_Store,
		"Path": path[2],
	})
}

// Get Cart By id
func (ch *cart_handler) GetCartById_Handler(c *gin.Context) {
	//***** ดึงค่า id user จาก token *******
	idUserType, errGetId := util.GetidToken(c)
	if errGetId != nil {
		c.JSON(401, util.Error_Custom(401, errGetId.Error(), "handler_cart_getallcartforstore_getidtoken"))
		return
	}

	// Get Param
	idParam, errGetParam := util.GetParam(c, "id")
	if errGetParam != nil {
		c.JSON(404, util.Error_Custom(404, errGetParam.Error(), "handler_product_createProduct_getParam"))
		return
	}

	// get path
	fullpath := c.FullPath()
	// fmt.Println(fullpath)
	path := strings.Split(fullpath, "/")

	cartResonse, errGetById := ch.cartService.GetCartById_Service(idParam, idUserType, path[2])
	if errGetById != nil {
		c.JSON(404, util.Error_Custom(404, errGetById.Error(), "handler_product_createProduct_GetCartById_Service"))
		return
	}

	c.JSON(200, gin.H{
		"Data": cartResonse,
		"Path": path[2],
	})

}


// UpdateCartStatus_Handler implements ICart_Handler
// สำหรับ  user เข้ามาแก้ไข order
func (ch *cart_handler) EditCartOrder_Handler(c *gin.Context) {
	newDataCart := models.CartRequest{}

	// get Param
	idCart, errGetParam := util.GetParam(c, "idcart")
	if errGetParam != nil {
		c.JSON(404, util.Error_Custom(404, errGetParam.Error(), "handler_product_createProduct_getParam"))
		return
	}

	//***** ดึงค่า id user จาก token *******
	idUser, errGetId := util.GetidToken(c)
	if errGetId != nil {
		c.JSON(401, util.Error_Custom(401, errGetId.Error(), "handler_product_createproduct_getidtoken"))
		return
	}

	// get data font end
	errShould := c.ShouldBindJSON(&newDataCart)
	if errShould != nil {
		c.JSON(400, util.Error_Custom(400, errShould.Error(), "handler_cart_editcartoder_shouldbindjson"))
		return
	}

	// ปั้นข้อมูลใหม่
	// ใช้ ข้อมูลตอนที่ดึง cart ของ store ส่งกลับมาก็ได้ จะง่ายมาก แต่ทำใน post man ไม่ได้
	// newDataCartStatus := models.StatusCartUpdate{
	// 	Id:       idCart,               // เอาจาก param
	// 	Id_Store: idStore,              // มากับ token
	// 	Status:   cartOrderData.Status, // font end
	// }

	// Validati data cart
	errValidate := ch.cartService.ValidateData_serviice(&newDataCart)
	if errValidate != nil {
		c.JSON(500, util.Error_Custom(500, errValidate.Error(), "handler_cart_editcartoder_validtor"))
		return
	}

	// fmt.Println(newDataCartStatus)
	// ส่งข้อมูลใหม่เข้าไปเลย โดยใช้ idCart เดิม
	// แล้วเอาอันเดิมลบออกให่หมด

	errUpdate := ch.cartService.EditCartOrder_Service(idUser, idCart, &newDataCart)
	if errUpdate != nil {
		c.JSON(500, util.Error_Custom(500, errUpdate.Error(), "handler_cart_editcartoder_editcartorder_service"))
		return
	}

	c.JSON(200, gin.H{
		"message": "Update Order Succress",
	})
}

//********** ต้องเป็นคนที่มีร้านค้าเท่านั้น *******
// หลังจาก 30 น่ที ถึงจะแก้ได้
func (ch *cart_handler) EditStatusCart_Handler(c *gin.Context) {

	status := models.StatusCartUpdate{}
	// get Param
	idCart, errGetParam := util.GetParam(c, "idcart")
	if errGetParam != nil {
		c.JSON(404, util.Error_Custom(404, errGetParam.Error(), "handler_product_createProduct_getParam"))
		return
	}

	//***** ดึงค่า id user จาก token *******
	idUser, errGetId := util.GetidToken(c)
	if errGetId != nil {
		c.JSON(500, util.Error_Custom(500, errGetId.Error(), "handler_cart_editstatuscart_getidtoken"))
		return
	}

	// ShouldBind
	errShould := c.ShouldBindJSON(&status)
	if errShould != nil {
		c.JSON(400, util.Error_Custom(400, errShould.Error(), "handler_cart_createcart_ShouldBindJSON"))
		return
	}

	// Validati data cart
	errValidate := ch.cartService.ValidateData_serviice(&status)
	if errValidate != nil {
		c.JSON(500, util.Error_Custom(500, errValidate.Error(), "handler_cart_editcartoder_validtor"))
		return
	}

	// to DB
	errEditCart := ch.cartService.EdiStatustCart_Service(idCart, idUser, &status)
	if errEditCart != nil {
		c.JSON(500, util.Error_Custom(500, errEditCart.Error(), "handler_cart_editstatuscart_editcartorder_service"))
		return
	}

	c.JSON(201, gin.H{
		"message": "Update Status Order Succress",
	})
}

func (ch *cart_handler) DeleteCart_Handler(c *gin.Context) {

	// get Param
	idCart, errGetParam := util.GetParam(c, "idcart")
	if errGetParam != nil {
		c.JSON(404, util.Error_Custom(404, errGetParam.Error(), "handler_product_createProduct_getParam"))
		return
	}

	idUser, errGetId := util.GetidToken(c)
	if errGetId != nil {
		c.JSON(401, util.Error_Custom(401, errGetId.Error(), "handler_cart_deletecart_getidtoken"))
		return
	}

	errDelete := ch.cartService.DeleteCart_Service(idUser, idCart)
	if errDelete != nil {
		c.JSON(500, util.Error_Custom(500, errDelete.Error(), "handler_cart_deletecart_deletecart_servie"))
		return
	}

	c.JSON(201, gin.H{
		"message": "delete Order Succress",
	})

}
