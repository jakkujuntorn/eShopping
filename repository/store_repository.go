package repository

import (
	"fmt"
	"myShopEcommerce/models"
	"gorm.io/gorm"
)

type store_Repository struct {
	db *gorm.DB
}

type IRepository_Store interface {
	CreateStore_Repo(store_Struct *models.Store) (*models.Store, error)
	GetStore_Repo(id_User int) (*models.Store, []models.Product, error)
	UpdateStore_Repo(id_User int, newStore *models.Store) error
	DeleteStore_Repo()
}

func NewStore_Repository(database *gorm.DB) IRepository_Store {
	return &store_Repository{
		db: database,
	}
}

// CreateStore implements IRepository_Store
func (sr *store_Repository) CreateStore_Repo(dataStore *models.Store) (*models.Store, error) {

	tx := sr.db.Table("store").Create(dataStore)
	if tx.Error != nil {
		return &models.Store{}, tx.Error
	}
	tx = sr.db.Table("store").Where("id_user=?", dataStore.IdUser).Find(&dataStore)
	if tx.Error != nil {
		return &models.Store{}, tx.Error
	}
	return dataStore, tx.Error
}

// DeleteStore implements IRepository_Store
func (*store_Repository) DeleteStore_Repo() {
	panic("unimplemented")
}

// GetStore implements IRepository_Store
func (sr *store_Repository) GetStore_Repo(id_User int) (*models.Store, []models.Product, error) {
	dataStore := models.Store{}
	allProduct := []models.Product{}

	tx := sr.db.Table("myshop.store").Where("id_user=?", id_User).Find(&dataStore)
	if tx.Error != nil {
		return &dataStore, allProduct, tx.Error
	}
	tx = sr.db.Table("myshop.products").Where("id_store=?", id_User).Find(&allProduct)
	if tx.Error != nil {
		return &dataStore, allProduct, tx.Error
	}
	fmt.Println("")
	return &dataStore, allProduct, tx.Error
}

// UpdateStore implements IRepository_Store
func (sr *store_Repository) UpdateStore_Repo(id_user int, newDataStore *models.Store) error {
	storeData := models.Store{}
	// หาข้อมูลที่จะอัพเดท
	tx := sr.db.Table("store").Where("id_user=?", id_user).First(&storeData)
	if tx.Error != nil {
		return tx.Error
	}

	// อัพเดท ทับข้อมูลเก่า
	tx = sr.db.Table("store").Where("id_user=?", id_user).Updates(newDataStore)
	return tx.Error

}
