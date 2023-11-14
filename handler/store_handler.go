package handler

import (
	_ "fmt"
	"myShopEcommerce/models"
	"myShopEcommerce/service"
	"myShopEcommerce/util"
	"strconv"

	"github.com/gin-gonic/gin"
)

type handler_Handler struct {
	storeHandler service.IStore_Service
}

type IStore_Handler interface {
	CreateStore_Handler(*gin.Context)
	GetStore_Handler(*gin.Context)
	UpdateStore_Handler(*gin.Context)
	DeleteStore_Handler()
}

func NewStore_Handler(storeHandler service.IStore_Service) IStore_Handler {
	return &handler_Handler{
		storeHandler: storeHandler,
	}
}

// CreateStore_Service implements IStore_Handler
func (hh *handler_Handler) CreateStore_Handler(c *gin.Context) {
	dataStore := models.Store{}

	// ดึงค่า id จาก token
	idToken, errGetid := util.GetidToken(c)
	if errGetid != nil {
		c.JSON(401, util.Error_Custom(401, errGetid.Error(), "handler_store_createStore_getidtoken"))
	}

	// ShouldBindJSON
	err := c.ShouldBindJSON(&dataStore)
	if err != nil {
		c.JSON(400, util.Error_Custom(400, err.Error(), "handler_store_createStore_shouldBinjson"))
		return
	}

	// เช็ค request ข้อมูลใน struct ใช้ จาก interface Service
	// เจอปัญหานี้  "validator: (nil *interface {})"
	// เพราะ ตรง validator ไปใส pointer
	errValidate := hh.storeHandler.ValidateData_serviice(dataStore)
	if errValidate != nil {
		c.JSON(400, util.Error_Custom(400, errValidate.Error(), "handler_store_createStore_validator"))
		return
	}

	dataStore.IdUser = strconv.Itoa(idToken)
	dataStore.IdStore = strconv.Itoa(idToken)

	//to service
	dataService, errCreate := hh.storeHandler.CreateStore_Service(&dataStore)
	if errCreate != nil {
		c.JSON(500, util.Error_Custom(500, errCreate.Error(), "handler_store_createStore_createstore_service"))
		return
	}

	c.JSON(201, gin.H{
		"Success": "true",
		"Data":    dataService,
	})
}

// DeleteStore_Service implements IStore_Handler
func (*handler_Handler) DeleteStore_Handler() {
	panic("unimplemented")
}

// GetStore_Service implements IStore_Handler
func (hh *handler_Handler) GetStore_Handler(c *gin.Context) {

	// ดึงค่าจาก token
	idUser, errGetid := util.GetidToken(c)
	if errGetid != nil {
		c.JSON(401, util.Error_Custom(401, errGetid.Error(), "handler_store_getStore_getIdtoken"))
	}

	// to DB
	dataStore_Service, errGetStore := hh.storeHandler.GetStore_Service(idUser)
	if errGetStore != nil {
		c.JSON(500, util.Error_Custom(500, errGetStore.Error(), "handler_store_getStore_service"))
		return
	}

	c.JSON(200, gin.H{
		"Success": "true",
		"Data":    dataStore_Service,
	})
}

// UpdateStore_Service implements IStore_Handler
func (hh *handler_Handler) UpdateStore_Handler(c *gin.Context) {
	// ดึงค่าจาก token
	idUser, errGetid := util.GetidToken(c)
	if errGetid != nil {
		c.JSON(401, util.Error_Custom(401, errGetid.Error(), "handler_store_updateStore_getIdtoken"))
	}

	// ดึงข้อมูล
	dataStore := models.Store{}
	errShoulddBind := c.ShouldBindJSON(&dataStore)
	if errShoulddBind != nil {
		c.JSON(400, util.Error_Custom(400, "Can not ShouldBindJSON", "handler_store_updateStore_shouldBindJson"))
		return
	}

	// เช็ค request ข้อมูลใน struct ใช้ จาก interface Service
	errValidate := hh.storeHandler.ValidateData_serviice(dataStore)
	if errValidate != nil {
		c.JSON(400, util.Error_Custom(400, errValidate.Error(), "handler_store_updateStore_validator"))
		return
	}

	// to service
	errUpdate := hh.storeHandler.UpdateStore_Service(idUser, &dataStore)
	if errUpdate != nil {
		c.JSON(500, util.Error_Custom(500, errUpdate.Error(), "handler_store_update_updateStore_service"))
		return
	}

	c.JSON(200, gin.H{
		"Success": "true",
		"message": "Update Success",
	})
}
