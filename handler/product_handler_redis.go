package handler

import (
	_ "time"

	"myShopEcommerce/models"
	"myShopEcommerce/service"
	"myShopEcommerce/util"
	"strconv"

	_ "math/rand"

	_ "github.com/Pallinder/go-randomdata"
	"github.com/gin-gonic/gin"
	// "github.com/google/uuid"
	
)

type product_handler_redis struct {
	productService service.IProduct_Service_Redis
}

type IProduct_handler_redis interface {
	CreateProduct_Service(*gin.Context)
	SearchProduct_Service(*gin.Context)
	GetAllMyProduct_Service(*gin.Context)
	UpdateProduct_Service(*gin.Context)
	DeleteProduct_Service(*gin.Context)
}

func NewProduct_Handler_redis(productService *service.IProduct_Service_Redis) IProduct_handler_redis {
	return &product_handler_redis{
		productService: *productService,
	}
}

// CreateProduct_Service implements IProduct_handler
func (ph *product_handler_redis) CreateProduct_Service(c *gin.Context) {
	productRequest := models.ProductRequest{}
	product := models.Product{}

	// ดึงค่าจาก token
	idUser, errGetid := util.GetidToken(c)
	if errGetid != nil {
		c.JSON(401, util.Error_Custom(401, errGetid.Error(), "handler_store_updateStore_getIdtoken"))
	}

	product.IdStore = idUser
	product.ID = uint(idUser)

	// ShouldBindJSON
	errShouldBind := c.ShouldBindJSON(&productRequest)
	if errShouldBind != nil {
		c.JSON(400, util.Error_Custom(400, errShouldBind.Error(), "handler_create_product"))
		return
	}

	// validate data
	errValidate := ph.productService.ValidateData_serviice(productRequest)
	if errValidate != nil {
		c.JSON(400, util.Error_Custom(400, errValidate.Error(), "handler_create_product"))
		return
	}

	// ปั้นข้อมูลใหม่
	product.Title = productRequest.Title
	product.Price = productRequest.Price
	product.Quantity = productRequest.Quantity
	product.Description = productRequest.Description
	product.Category = productRequest.Category


	// DB
	// ที่ Error 1054 unknown column in field list
	// เพราะ Table มันไปทำ foreign key ไว้ อาจทำแบบผิดๆ มันเลยไม่ให้ใสค่า
	errCreateDB := ph.productService.CreateProduct_Service(&product)
	if errCreateDB != nil {
		c.JSON(500, util.Error_Custom(500, errCreateDB.Error(), "handle_createproduct"))
		return
	}

	c.JSON(201, gin.H{
		"Success": "true",
		"Message": "Create Product Success",
	})
}

// DeleteProduct_Service implements IProduct_handler
func (*product_handler_redis) DeleteProduct_Service(*gin.Context) {
	panic("unimplemented")
}

// GetProduct_Service implements IProduct_handler
func (ph *product_handler_redis) GetAllMyProduct_Service(c *gin.Context) {

	//***** ดึงค่า id user จาก token *******
	idStore, errGetId := util.GetidToken(c)
	if errGetId != nil {
		c.JSON(401, util.Error_Custom(401, errGetId.Error(), "handler_cart_createcart_getidtoken"))
		return
	}

	// รับค่า page and limit จาก Query
	page := c.DefaultQuery("page", "0")
	limit := c.DefaultQuery("limit", "5")
	// แปลง str to int
	pages, _ := strconv.Atoi(page)
	limits, _ := strconv.Atoi(limit)
	// ปั้น pagination
	pageination := models.Pagination{
		Page:   pages,
		Limit:  limits,
		Offset: (pages - 1) * limits,
	}

	dataProduct, status, errGetMyProduct := ph.productService.GetAllMyProduct_Service(idStore, &pageination)
	if errGetMyProduct != nil {
		c.JSON(500, util.Error_Custom(500, errGetMyProduct.Error(), "handle_product_getallmyproduct"))
		return
	}

	c.JSON(200, gin.H{
		"Success": "true",
		"Data":    dataProduct,
		"Status":  status,
	})
}

// SearchProduct_Service implements IProduct_handler
func (ph *product_handler_redis) SearchProduct_Service(c *gin.Context) {

	title := c.DefaultQuery("title", "")
	priceStart := c.DefaultQuery("priceStart", "0")
	priceEnd := c.DefaultQuery("priceEnd", "10000")
	category := c.DefaultQuery("category", "")

	ps, _ := strconv.Atoi(priceStart)
	pe, _ := strconv.Atoi(priceEnd)

	productSearch := models.Product_Search{
		Title:      title,
		PriceStart: ps,
		PriceEnd:   pe,
		Category:   category,
	}

	page := c.DefaultQuery("page", "0")
	limit := c.DefaultQuery("limit", "5")
	pages, _ := strconv.Atoi(page)
	limits, _ := strconv.Atoi(limit)

	pageination := models.Pagination{
		Page:   pages,
		Limit:  limits,
		Offset: (pages - 1) * limits,
	}

	// Search
	dataSearch, total, errSearch := ph.productService.SearchProduct_Service(&productSearch, &pageination)
	if errSearch != nil {
		c.JSON(500, util.Error_Custom(500, errSearch.Error(), "handler_product_searchproduct"))
		return
	}

	//customer responer
	c.JSON(200, models.PaginationResponse("true",
		total,
		pageination.Page,
		pageination.Limit,
		dataSearch,
	))
}

// UpdateProduct_Service implements IProduct_handler
func (ph *product_handler_redis) UpdateProduct_Service(c *gin.Context) {
	//***** ดึงค่า id user จาก token *******
	idStore, errGetId := util.GetidToken(c)
	if errGetId != nil {
		c.JSON(401, util.Error_Custom(401, errGetId.Error(), "handler_cart_createcart_getidtoken"))
		return
	}

	// get id product form param
	idProduct, errGetParam := util.GetParam(c, "idproduct")
	if errGetParam != nil {
		c.JSON(404, util.Error_Custom(404, errGetParam.Error(), "handler_product_createProduct_getParam"))
		return
	}

	productUpdate := models.ProductUpdate{}

	errShouldBindJSON := c.ShouldBindJSON(&productUpdate)
	if errShouldBindJSON != nil {
		c.JSON(400, util.Error_Custom(400, "Err ShouldBindJson", "product_handler_UpdateProduct"))
		return
	}

	// validation data
	errValidation := ph.productService.ValidateData_serviice(&productUpdate)
	if errValidation != nil {
		c.JSON(400, util.Error_Custom(400, errValidation.Error(), "handler_product_UpdateProduct"))
		return
	}

	errUpdate := ph.productService.UpdateProduct_Service(idProduct, idStore, &productUpdate)
	if errUpdate != nil {
		c.JSON(500, util.Error_Custom(500, errUpdate.Error(), "handler_product_UpdateProduct"))
		return
	}

	c.JSON(200, gin.H{
		"Success": "true",
		"Message": "Update Product Success",
	})

}
