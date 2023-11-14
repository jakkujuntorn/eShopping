package service

import (
	"encoding/json"
	"fmt"

	"myShopEcommerce/models"
	"myShopEcommerce/repository"
	"myShopEcommerce/util"
	"time"

	"github.com/go-redis/redis"
	"github.com/spf13/viper"

	_ "github.com/google/uuid"
)

type product_service_redis struct {
	productRepo repository.IRepository_Product
	redisClient *redis.Client
}

type IProduct_Service_Redis interface {
	
	CreateProduct_Service(product *models.Product) error
	SearchProduct_Service(searchText *models.Product_Search, pagination *models.Pagination) ([]models.ProductResponse, int, error)
	GetAllMyProduct_Service(idStore int, pagination *models.Pagination) ([]models.ProductResponse, string, error)
	UpdateProduct_Service(idProduct int, idStore int, newProduct *models.ProductUpdate) error
	DeleteProduct_Service(idUser int) error
	ValidateData_serviice(dataValidate interface{}) error
}

// func NewProduct_Service_Redis(productRepo *repository.IRepository_Product_redis, redisClient *redis.Client) IProduct_Service_Redis {
// 	return &product_service_redis{
// 		productRepo: *productRepo,
// 		redisClient: redisClient,
// 	}
// }

func NewProduct_Service_Redis(productRepo *repository.IRepository_Product, redisClient *redis.Client) IProduct_Service_Redis {
	return &product_service_redis{
		productRepo: *productRepo, 
		redisClient: redisClient,
	}
}

func (ps *product_service_redis) ValidateData_serviice(dataValidate interface{}) error {
	errValidate := util.ValidateDataUser(dataValidate)
	return errValidate
}

// CreateProduct_Service implements IProduct_Service
func (ps *product_service_redis) CreateProduct_Service(product *models.Product) error {
	errCreate := ps.productRepo.CreateProduct_Repo(product)
	if errCreate != nil {
		return errCreate
	}
	return nil
}

// DeleteProduct_Service implements IProduct_Service
func (*product_service_redis) DeleteProduct_Service(idProduct int) error {
	panic("unimplemented")
}

// GetProduct_Service implements IProduct_Service
func (ps *product_service_redis) GetAllMyProduct_Service(idStore int, page *models.Pagination) ([]models.ProductResponse, string, error) {
	productResponse := []models.ProductResponse{}

	// ************ Get from Redis ********
	if productsJson, err := ps.redisClient.Get(viper.GetString("redis.redis_key")).Result(); err == nil {
		// 1.2ถ้า json err == nil แสดงว่ามีค่าจาก redis
		// 1.3 แปลงค่ามาเป็น struct ด้วย json
		if json.Unmarshal([]byte(productsJson), &productResponse); err == nil {
			// แปลงข้อมูล จาก productsJson สู่ products
			fmt.Println(" ******Service Redis DB *******")
			return productResponse, "Redis", nil
		}
	}

	// *********** Get from DB **************
	productService, errAllProduct := ps.productRepo.GetAllMyProduct_Repo(idStore, page)
	if errAllProduct != nil {

		return productResponse, "", errAllProduct
	}
	// ปั้น ข้อมูลใหม่
	for _, data := range productService {
		product := models.ProductResponse{
			// IdStore:     data.IdStore,
			IdProduct:   data.IdProduct,
			Title:       data.Title,
			Price:       data.Price,
			Quantity:    data.Quantity,
			Description: data.Description,
			Category:    data.Category,
		}
		productResponse = append(productResponse, product)
	}

	//******************* Set in Redis *******************
	if data, err := json.Marshal(productResponse); err == nil {
		// json.Marshal retun []byte
		// ไม่ต้องจัดการ error ก็ได้ เพราะถ้า set ไม่ได้ มันจะไปอ่านจาก DB เอง
		// er := s.redisClient.Set(context.Background(), key, string(data), time.Second*10)
		// _ = er

		ps.redisClient.Set(viper.GetString("redis.redis_key"), string(data), time.Second*10)
	}
	fmt.Println(" ******Service  DB *******")
	return productResponse, "DB", nil
}

// SearchProduct_Service implements IProduct_Service
func (ps *product_service_redis) SearchProduct_Service(textSearch *models.Product_Search, page *models.Pagination) ([]models.ProductResponse, int, error) {
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
	return productResponse, count, nil
}

// UpdateProduct_Service implements IProduct_Service
func (ps *product_service_redis) UpdateProduct_Service(id_product int, id_store int, newProduct *models.ProductUpdate) error {
	errUpdate := ps.productRepo.UpdateProduct_Repo(id_product, id_store, newProduct)
	if errUpdate != nil {
		return errUpdate
	}
	return nil
}
