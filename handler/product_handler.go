package handler

import (
	_ "fmt"

	"myShopEcommerce/models"
	"myShopEcommerce/service"
	"myShopEcommerce/util"
	"strconv"

	"github.com/gin-gonic/gin"

	_ "github.com/Pallinder/go-randomdata"
)

type product_handler struct {
	productService service.IProduct_Service
}

type IProduct_handler interface {
	CreateProduct_Service(*gin.Context)
	SearchProduct_Service(*gin.Context)
	GetMyProduct_Service(*gin.Context)
	UpdateProduct_Service(*gin.Context)
	DeleteProduct_Service(*gin.Context)
}

func NewProduct_Handler(productService *service.IProduct_Service) IProduct_handler {
	return &product_handler{
		productService: *productService,
	}
}

// CreateProduct_Service implements IProduct_handler
func (ph *product_handler) CreateProduct_Service(c *gin.Context) {
	productRequest := models.ProductRequest{}
	product := models.Product{}

	// ดึงค่า id จาก token
	idUser, errGetid := util.GetidToken(c)
	if errGetid != nil {
		c.JSON(401, util.Error_Custom(401, errGetid.Error(), "handler_product_createProduct_getidtoken"))
	}
	product.IdStore = idUser
	product.ID = uint(idUser)

	// ShouldBindJSON
	errShouldBind := c.ShouldBindJSON(&productRequest)
	if errShouldBind != nil {
		c.JSON(400, util.Error_Custom(400, errShouldBind.Error(), "handler_product_createProduct_ShouldBindJSON"))
		return
	}
	// validate data
	errValidate := ph.productService.ValidateData_serviice(productRequest)
	if errValidate != nil {
		// ลองทำเพื่อให้ เช็ค error ว่า เป็น error ประเภทไหน ****
		util.Handler_Error(c, errValidate)
		return
	}

	// ****** ปั้นข้อมูลใหม่ ********
	// ProductRequest ไปเป็น Product
	product.Title = productRequest.Title
	product.Price = productRequest.Price
	product.Quantity = productRequest.Quantity
	product.Description = productRequest.Description
	product.Category = productRequest.Category

	// to Service
	errCreateDB := ph.productService.CreateProduct_Service(&product)
	if errCreateDB != nil {
		// ใช้เช็ค type error
		util.Handler_Error(c, errCreateDB)
		return
	}

	c.JSON(200, gin.H{
		"Success": "true",
		"Message": "Create Product Success",
	})
}

// DeleteProduct_Service implements IProduct_handler
func (*product_handler) DeleteProduct_Service(*gin.Context) {
	panic("unimplemented")
}

// GetProduct_Service implements IProduct_handler
func (ph *product_handler) GetMyProduct_Service(c *gin.Context) {
	// ดึงค่า id จาก token
	id_Store, errGetid := util.GetidToken(c)
	if errGetid != nil {
		c.JSON(401, util.Error_Custom(401, errGetid.Error(), "handler_product_createProduct_getidtoken"))
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

	dataProduct, errGetMyProduct := ph.productService.GetMyProduct_Service(id_Store, &pageination)
	if errGetMyProduct != nil {
		// ให้เช็ค type error
		util.Handler_Error(c, errGetMyProduct)
		return
	}

	c.JSON(200, gin.H{
		"Success": "true",
		"Data":    dataProduct,
	})
}

// SearchProduct_Service implements IProduct_handler
func (ph *product_handler) SearchProduct_Service(c *gin.Context) {
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
		// เช็ค type error
		util.Handler_Error(c, errSearch)
		return
	}

	c.JSON(200, models.PaginationResponse("true",
		total,
		pageination.Page,
		pageination.Limit,
		dataSearch,
	))
}

// UpdateProduct_Service implements IProduct_handler
func (ph *product_handler) UpdateProduct_Service(c *gin.Context) {
	// ดึงค่า id จาก token
	id_Store, errGetid := util.GetidToken(c)
	if errGetid != nil {
		c.JSON(401, util.Error_Custom(401, errGetid.Error(), "handler_product_createProduct_getidtoken"))
	}

	// get id product form param
	idProduct, errGetParam := util.GetParam(c,"idproduct")
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
		// เช็ค type error
		util.Handler_Error(c, errValidation)
		return
	}

	// to Service
	errUpdate := ph.productService.UpdateProduct_Service(idProduct, id_Store, &productUpdate)
	if errUpdate != nil {
		// เช็ค type error
		util.Handler_Error(c, errUpdate)
		return
	}

	c.JSON(200, gin.H{
		"Success": "true",
		"Message": "Update Product Success",
	})

}
