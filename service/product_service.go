package service

import (
	_ "fmt"
	"myShopEcommerce/models"
	"myShopEcommerce/repository"
	"myShopEcommerce/util"

	_ "github.com/google/uuid"
)

type product_service struct {
	productRepo repository.IRepository_Product
}

type IProduct_Service interface {
	CreateProduct_Service(product *models.Product) error
	SearchProduct_Service(product_Search *models.Product_Search, pagination *models.Pagination) (productResponse []models.ProductResponse, count_Product int, err error)
	GetMyProduct_Service(id_Store int, pagination *models.Pagination) (product_Response []models.ProductResponse, err error)
	UpdateProduct_Service(id_Product, id_Store int, productUpdate *models.ProductUpdate) error
	DeleteProduct_Service(id_Product int) error
	ValidateData_serviice(dataValidate interface{}) error
}

func NewProduct_Service(productRepo *repository.IRepository_Product) IProduct_Service {
	return &product_service{
		productRepo: *productRepo,
	}
}

func (ps *product_service) ValidateData_serviice(dataValidate interface{}) error {
	return util.ValidateDataUser(dataValidate)
}

// CreateProduct_Service implements IProduct_Service
func (ps *product_service) CreateProduct_Service(product *models.Product) error {
	err := ps.productRepo.CreateProduct_Repo(product)
	if err != nil {
		return err
	}
	return nil
}

// DeleteProduct_Service implements IProduct_Service
func (*product_service) DeleteProduct_Service(idProduct int) error {
	panic("unimplemented")
}

// GetProduct_Service implements IProduct_Service
func (ps *product_service) GetMyProduct_Service(idStore int, page *models.Pagination) ([]models.ProductResponse, error) {
	productResponse := []models.ProductResponse{}

	// get from DB
	productService, err := ps.productRepo.GetAllMyProduct_Repo(idStore, page)
	if err != nil {
		return productResponse, err
	}

	for _, data := range productService {
		product := models.ProductResponse{
			// IdStore:     data.IdStore,
			Title:       data.Title,
			Price:       data.Price,
			Quantity:    data.Quantity,
			Description: data.Description,
			Category:    data.Category,
		}
		productResponse = append(productResponse, product)
	}
	return productResponse, nil
}

// SearchProduct_Service implements IProduct_Service
func (ps *product_service) SearchProduct_Service(textSearch *models.Product_Search, page *models.Pagination) ([]models.ProductResponse, int, error) {
	productResponse := []models.ProductResponse{}

	// search
	dataSearch, count, errSearch := ps.productRepo.SearchProduct_Repo(textSearch, page)
	if errSearch != nil {
		return productResponse, 0, errSearch
	}
	//for loop product to productResponse
	for _, data := range dataSearch {
		product := models.ProductResponse{
			Title:       data.Title,
			Price:       data.Price,
			Quantity:    data.Quantity,
			Description: data.Description,
			Category:    data.Category,
		}
		productResponse = append(productResponse, product)
	}

	return productResponse, count, errSearch
}

// UpdateProduct_Service implements IProduct_Service
func (ps *product_service) UpdateProduct_Service(id_product int, id_store int, newProduct *models.ProductUpdate) error {
	errUpdate := ps.productRepo.UpdateProduct_Repo(id_product, id_store, newProduct)
	if errUpdate != nil {	
		return errUpdate
	}
	return nil
}
