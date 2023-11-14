package service

import (
	"errors"
	"fmt"
	"myShopEcommerce/models"
	"myShopEcommerce/repository"
	"myShopEcommerce/util"
)

type store_Service struct {
	storeService repository.IRepository_Store
}

type IStore_Service interface {
	CreateStore_Service(dataStore *models.Store) (dataStore_Response *models.Store, err error)
	GetStore_Service(idUser int) (dataStore_Response *models.Store_Response, err error)
	UpdateStore_Service(id_User int, dataStore *models.Store) error
	DeleteStore_Service()
	ValidateData_serviice(dataValidate interface{}) error
}

func NewStore_Service(storeService repository.IRepository_Store) IStore_Service {
	return &store_Service{
		storeService: storeService,
	}
}

// Validate DAta
func (ss *store_Service) ValidateData_serviice(dataValidate interface{}) error {
	return util.ValidateDataUser(dataValidate)
}

// CreateStore_Service implements IStore_Service
func (ss *store_Service) CreateStore_Service(dataStore *models.Store) (*models.Store, error) {
	storeResponse := models.Store{}

	storeRepo, errCreate := ss.storeService.CreateStore_Repo(dataStore)
	if errCreate != nil {
		return &storeResponse, errCreate
	}
	return storeRepo, nil
}

// DeleteStore_Service implements IStore_Service
func (*store_Service) DeleteStore_Service() {
	panic("unimplemented")
}

// GetStore_Service implements IStore_Service
func (ss *store_Service) GetStore_Service(idUser int) (*models.Store_Response, error) {
	// dataStore_Service := models.Store{}
	store_Response := models.Store_Response{}
	dataStro_Repo, productsRepo, err := ss.storeService.GetStore_Repo(idUser)
	if err != nil {
		return &store_Response, err
	}

	// product ในร้าน
	products := []models.ProductResponse{}
	for _, data := range productsRepo {
		product := models.ProductResponse{
			Title:       data.Title,
			Price:       data.Price,
			Quantity:    data.Quantity,
			Description: data.Description,
			Category:    data.Category,
		}

		products = append(products, product)
	}

	if dataStro_Repo.IdUser == "" {
		// gorm ใช้ find ถึงไม่เจอก็จะไม่มี error
		// เลยต้องสร้าง error ขึ้นมาเอง
		return &store_Response, errors.New("row not found")
	}

	fmt.Println("")
	store_Response.IdUser = dataStro_Repo.IdUser
	store_Response.StoreName = dataStro_Repo.StoreName
	store_Response.Products = products
	return &store_Response, err
}

// UpdateStore_Service implements IStore_Service
func (ss *store_Service) UpdateStore_Service(id_User int, newDataStore *models.Store) error {

	err := ss.storeService.UpdateStore_Repo(id_User, newDataStore)
	if err != nil {
		return err
	}
	return nil
}
