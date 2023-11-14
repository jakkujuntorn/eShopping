package handler

import (

	"myShopEcommerce/models"
	"myShopEcommerce/service"
	"myShopEcommerce/util"
	"strconv"

	_ "github.com/Pallinder/go-randomdata"
	"github.com/gin-gonic/gin"
)

type user_handler struct {
	userService service.IUser_Sevice
}

type IUser_Handler interface {
	Register_Handler(*gin.Context)
	GetAllUser_Handler(*gin.Context)
	UpdateUser_Handler(*gin.Context)
	Login_Handler(*gin.Context)
	GetUserByUsername(*gin.Context)
}

func NewUser_Handler(userService service.IUser_Sevice) IUser_Handler {
	return &user_handler{
		userService: userService,
	}
}

// func NewUser_Handler_Kafka(userService service.IUser_Sevice) IUser_Handler {
// 	return &user_handler{
// 		userService: userService,
// 	}
// }

// CreateUser_Handler implements IUser_Handler
func (uh *user_handler) Register_Handler(c *gin.Context) {

	var dataUser models.UserRequest

	err := c.ShouldBindJSON(&dataUser)
	if err != nil {
		c.JSON(400, util.Error_Custom(500, err.Error(), "register-ShouldBindJson"))
		return
	}

	// validator data
	errValidater := uh.userService.ValidateData_serviice(dataUser)
	if errValidater != nil {
		c.JSON(400, util.Error_Custom(400, errValidater.Error(), "register-Validator"))
		return
	}

	dataRegister, err := uh.userService.CreateUser_Service(&dataUser)
	if err != nil {
		c.JSON(500, util.Error_Custom(500, err.Error(), "register-CreateTo Service"))
		return
	}

	// register แล้วจะ auto login เลยดีไหม
	//น่าจะอยู่ font end ว่าให้ login ไปเบยรึป่าว
	c.JSON(201, gin.H{
		"Success": "success",
		"Data":    dataRegister,
	})

	//******** custom response **********
	// c.JSON(200, util.Data_Response(dataRegister))
}

// GetDataUser_Handler implements IUser_Handler
func (uh *user_handler) GetAllUser_Handler(c *gin.Context) {

	page := c.DefaultQuery("page", "0")
	limit := c.DefaultQuery("limit", "5")
	pages, _ := strconv.Atoi(page)
	limits, _ := strconv.Atoi(limit)

	// pagiantion
	pageination := models.PaginationDB(pages, limits, (pages-1)*limits)

	dataUser, total, err := uh.userService.GetAllUser_Service(&pageination)
	if err != nil {
		c.JSON(500, util.Error_Custom(500, err.Error(), "user_handler_GetDataUser"))
		return
	}

	// ใช้ find
	if len(dataUser) == 0 {
		c.JSON(404, util.Error_Custom(404, "page more User", "user_handler_GetDataUser"))
		return
	}

	c.JSON(200, models.PaginationResponse(
		"true",
		total,
		pageination.Page,
		pageination.Limit,
		dataUser,
	))
}

// UpdateUser_Handler implements IUser_Handler
func (uh *user_handler) UpdateUser_Handler(c *gin.Context) {
	var dataUpdate = models.UserUpdate{}

	// ดึงค่าจาก token
	idUser, errGetid := util.GetidToken(c)
	if errGetid != nil {
		c.JSON(401, util.Error_Custom(401, errGetid.Error(), "handler_store_updateStore_getIdtoken"))
	}

	err := c.ShouldBindJSON(&dataUpdate)
	if err != nil {
		c.JSON(400, util.Error_Custom(400, "Can not ShouldBindJSON", "user_handleruUpdateUser"))
	}

	// validator data
	errValidate := uh.userService.ValidateData_serviice(&dataUpdate)
	if errValidate != nil {
		c.JSON(400, util.Error_Custom(400, errValidate.Error(), "user_handler_updateUser"))
		return
	}

	// to service
	errUpdate := uh.userService.UpdateUser_Service(idUser, &dataUpdate)
	if errUpdate != nil {
		c.JSON(500, util.Error_Custom(500, errUpdate.Error(), "user_handler_updateUser"))
		return
	}

	c.JSON(200, gin.H{
		"Success": "true",
		"Message": "Update Success",
	})

}

func (uh *user_handler) Login_Handler(c *gin.Context) {
	dataLogin := models.UserLogin{}
	errShouldBind := c.ShouldBindJSON(&dataLogin)
	if errShouldBind != nil {
		c.JSON(400, util.Error_Custom(400, errShouldBind.Error(), "user_handler_login"))
		return
	}

	// Validator data
	errValidate := uh.userService.ValidateData_serviice(dataLogin)
	if errValidate != nil {
		c.JSON(400, util.Error_Custom(400, errValidate.Error(), "user_handler_validator"))
		return
	}

	// to DB
	dataResponse, errLogin := uh.userService.Login_Service(&dataLogin)
	if errLogin != nil {
		c.JSON(500, util.Error_Custom(500, errLogin.Error(), "user_handler_login"))
		return
	}

	// สร้างแบบนี้ไม่ error
	// var userResponse = models.UserResponse{}
	// สร้างแบบนี้ error เหลือง ว่าไมได้ใช้งาตวแปรนี้
	userResponse := models.UserResponse{}

	userResponse = *dataResponse

	//********** Create Token ********
	token, err := util.CreateToken(&userResponse)
	if err != nil {
		c.JSON(500, util.Error_Custom(500, "Can not Create Token", "user_handler_login"))
		return
	}

	//*********** ทำอะไร ***********
	// c.SetCookie("Authorization", token, 3600*24*30, "", "", false, true)

	c.JSON(200, gin.H{
		"Success":  "true",
		"Date":     dataResponse,
		"jwtToken": token,
	})
}

func (uh *user_handler) GetUserByUsername(c *gin.Context) {
	// Get User Token
	userToken, errGetUser := util.GetUserToken(c)
	if errGetUser != nil {
		c.JSON(401, util.Error_Custom(401, errGetUser.Error(), "handler_cart_deletecart_getidtoken"))
		return
	}

	// to Service
	dataResponse, errGetUser := uh.userService.GetDataByUsername_Service(userToken)
	if errGetUser != nil {
		c.JSON(500, util.Error_Custom(500, errGetUser.Error(), "user_handler_GetUserByName"))
		return
	}

	c.JSON(200, gin.H{
		"Success": "true",
		"Data":    dataResponse,
	})
}
