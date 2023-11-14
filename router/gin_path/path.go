package ginpath

import (
	"github.com/gin-gonic/gin"
	"myShopEcommerce/handler"
	"myShopEcommerce/util"
)

func GinRouter_User(handlerUser handler.IUser_Handler,
	handlerStore handler.IStore_Handler,
	handlerCart handler.ICart_Handler,

	// handlerProduct handler.IProduct_handler,
	handlerProduct handler.IProduct_handler_redis,
) {
	r := gin.Default()

	// 3 คำสั่งนี้ = gin.DEfault()  ******
	// rr:=gin.New()
	// rr.Use(gin.Logger(),gin.Recovery())
	// rr.Use(gin.Recovery())

	//********* Router User*******
	//********* Register *********
	r.POST("/register", handlerUser.Register_Handler)

	//********** Login *******
	r.POST("/login", handlerUser.Login_Handler)

	//************ Search
	r.POST("/search/", handlerProduct.SearchProduct_Service)

	// ใช้ VerifyToken ขึ้น token Unexpected ******
	r.Use(util.VerifyToken)

	// ***********  path ด้านล่างนี้ต้องผ่านการ login ก่อน ***********
	// ******** User Router ************
	userPath := r.Group("")
	{
		// GetAll Users for admin ******
		userPath.GET("/users", util.VerifyRole, handlerUser.GetAllUser_Handler)

		// GetByUserName Users for admin ******
		// userPath.GET("/users", util.VerifyRole, handlerUser.GetDataUser_Handler)

		// Get My user ******
		userPath.GET("/users/:username", handlerUser.GetUserByUsername)

		// update *********
		userPath.PATCH("/user/:id", handlerUser.UpdateUser_Handler)

	}

	//**************** get user protect path **********
	// protechPath :=r.Group("",util.VerifyToken)
	// protechPath.GET("/user", handlerFunc.GetDataUser_Handler)
	// protechPath.POST("/user", handlerFunc.CreateUser_Handler)
	// protechPath.PATCH("/user/:id", handlerFunc.UpdateUser_Handler)

	// สร้าง Group API

	//********* Router Store *******
	storePath := r.Group("")
	{
		// Get  Store
		storePath.GET("/store/:username", handlerStore.GetStore_Handler)

		// Create Store
		storePath.POST("/store/:id", handlerStore.CreateStore_Handler)

		// Update Store
		storePath.PATCH("/store/:id", handlerStore.UpdateStore_Handler)

	}

	// ********** Router Product ************
	// *****  product ของร้าน
	productPath := r.Group("")
	{
		// Get Product All in my store
		productPath.GET("/:storename/products", handlerProduct.GetAllMyProduct_Service)
		// Create Product in my store
		productPath.POST("/:storename/product", handlerProduct.CreateProduct_Service)
		// Update Product in my store
		productPath.PATCH("/:storename/product/:idproduct", handlerProduct.UpdateProduct_Service)
		// Delete Product in my store
		productPath.DELETE("/:storename/product")
	}

	// ***************  Router Cart *******************
	cartPath := r.Group("")
	{
		// For User เราสั่งอะไรไปบ้าง  เอาทั้งหมด ทำงานไดเปกติ
		cartPath.GET("/cart/cartuser", handlerCart.GetCartByUserTYpe_Handler)
		// For Store  มีใครสั่งสินค่าเราบ้าง เอาทั้งหมด ทำงานไดเปกติ
		cartPath.GET("/cart/cartstore", handlerCart.GetCartByUserTYpe_Handler)

		// เลือก เฉพาะ cart ตาม id   User
		cartPath.GET("/cart/cartuserbyid/:id", handlerCart.GetCartById_Handler)

		// เลือก เฉพาะ cart ตาม id   Store
		cartPath.GET("/cart/cartstorebyid/:id", handlerCart.GetCartById_Handler)

		//Create Cart  เราสั่งสินค้า
		cartPath.POST("/cart/create", handlerCart.CreateCart_Handler)

		//Edit  Cart for user
		// user จะแก้ไขได้ ต้องไม่เกิน ครึ่ง ช.ม. จาก การ create
		// เปลี่ยน จำนวน - เแลี่ยนจำนวนแล้วเซฟทับ
		// ลบ สินค้าออก - เอาสินค้าบางชิ้นออก
		// ยกเลิก order - เอา แค่ idcart มาลบ หรือ จะเอาแค่ ยกเลิก order เพราะส่วนอื่นจะทำที่ font end
		// เอา id cart มาหาก็น่าจะพอ แต่ user ต้องเช็คก่อนว่าเป็นเจ้าของจริงๆ หา id จาก token
		// id token,idcart มาเช็ค  *********
		cartPath.PATCH("/cart/editusercart/:idcart", handlerCart.EditCartOrder_Handler)
		
		// Edit cart status for Store
		// ร้านค้าเขามายืนยันการสั่งสินค้า หลัง 30 นาที
		cartPath.PATCH("/cart/storeeditstatuscart/:idcart", handlerCart.EditStatusCart_Handler)

		// Cancle Order by User
		//จะทำได้เมื่ไม่เกิน 30 นาที หลัง create
		// cartPath.PATCH("/cart/carteditstatusstore/:idcart", handlerCart.EditStatusCart_Handler)
		cartPath.DELETE("/cart/deletecart/:idcart", handlerCart.DeleteCart_Handler)

	}

	r.Run(":8000")

}
