package repository

import (
	"fmt"
	"myShopEcommerce/models"

	"gorm.io/gorm"
)

type product_Repository struct {
	db *gorm.DB
}

type IRepository_Product interface {
	CreateProduct_Repo(product *models.Product) error
	SearchProduct_Repo(productSearch *models.Product_Search, paginaTion *models.Pagination) ([]models.Product, int, error)
	GetAllMyProduct_Repo(id_Store int, pagination *models.Pagination) ([]models.Product, error)
	UpdateProduct_Repo(id_Product, id_Store int, productUpdate *models.ProductUpdate) error
	DeleteProduct_Repo(id_Product int) error
}

func NewProduct_Repository(database *gorm.DB) IRepository_Product {
	return &product_Repository{db: database}
}

// CreateStore implements IRepository_Product
func (pr *product_Repository) CreateProduct_Repo(product *models.Product) error {
	tx := pr.db.Table("products").Create(product)
	return tx.Error
}

func (pr *product_Repository) SearchProduct_Repo(text *models.Product_Search, page *models.Pagination) ([]models.Product, int, error) {
	product := []models.Product{}
	var count int64
	// custom Query
	title := fmt.Sprintf("%s%s%s%s%s", "", "%", text.Title, "%", "")
	category := fmt.Sprintf("%s%s%s%s%s", "", "%", text.Category, "%", "")

	tx := pr.db.Table("products").
		Where("title LIKE ?", title).
		Where("price BETWEEN ? and ?", text.PriceStart, text.PriceEnd).
		Where("category LIKE ?", category).
		Limit(page.Limit).
		Offset(page.Offset).
		Count(&count).
		Find(&product)

	fmt.Println(count) // ค่านี้ คำนวณยังไง **********
	// fmt.Print(tx.RowsAffected) // จำนวนที่แสดงต่อหน้า

	return product, int(count), tx.Error
}

// DeleteStore implements IRepository_Product
func (*product_Repository) DeleteProduct_Repo(idProduct int) error {
	panic("unimplemented")
}

// GetStore implements IRepository_Product
func (pr *product_Repository) GetAllMyProduct_Repo(id_Store int, pagination *models.Pagination) ([]models.Product, error) {
	product := []models.Product{}
	tx := pr.db.Table("products").Where("id_store=?", id_Store).Limit(pagination.Limit).Offset(pagination.Offset).Find(&product)
	return product, tx.Error
}

// UpdateStore implements IRepository_Product
func (pr *product_Repository) UpdateProduct_Repo(id_Product int, id_store int, newDataProduct *models.ProductUpdate) error {
	product := models.Product{}

	// หาข้อมูลขึ้นมาก่อน เอา  id_store และ id_product
	// use first
	tx := pr.db.Table("products").First(&product, "id_product=? and id_store=?", id_Product, id_store)
	if tx.Error != nil {
		fmt.Println("Repo: ", tx.Error.Error())
		return tx.Error
	}

	// อัพเดทข้อมูลทับลงไป
	tx = pr.db.Table("products").Where("id_product=? and id_store=?", id_Product, id_store).Updates(newDataProduct)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
