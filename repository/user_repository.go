package repository

import (
	_ "fmt"
	"myShopEcommerce/models"

	"github.com/IBM/sarama"
	"gorm.io/gorm"
)

type user_Repository struct {
	db       *gorm.DB
	producer sarama.SyncProducer
}

type IUser_Repository interface {
	CreateUser_Repo(userRequest *models.UserRequest) (*models.UserRequest, error)
	GetAllUser_Repo(pagination *models.Pagination) (userRequest []models.UserRequest, count int, err error)
	UpdateUser_Repo(id_User int, userUpdate *models.UserUpdate) error
	Login_Repo(userName string) (userRequest *models.UserRequest, err error)
	GetDataByUsername_Repo(userName string) (userRequest *models.UserRequest, err error)
}

func NewUser_Repository(database *gorm.DB) IUser_Repository {
	return &user_Repository{db: database}
}

func (u *user_Repository) CreateUser_Repo(userData *models.UserRequest) (*models.UserRequest, error) {
	tx := u.db.Table("user").Create(userData)
	return userData, tx.Error

}

func (u *user_Repository) GetAllUser_Repo(page *models.Pagination) ([]models.UserRequest, int, error) {
	// fmt.Println("start get user")
	users := []models.UserRequest{}
	var count int64
	tx := u.db.Table("user").Limit(page.Limit).Offset(page.Offset).Count(&count).Find(&users)
	return users, int(count), tx.Error
}

func (u *user_Repository) UpdateUser_Repo(id_User int, newData *models.UserUpdate) error {
	userData := models.UserRequest{}

	// fmt.Println("repo")
	tx := u.db.Table("user").First(&userData, id_User)
	if tx.Error != nil {
		return tx.Error
	}

	tx = u.db.Table("user").Model(models.UserRequest{}).Where("id=?", id_User).Updates(newData)

	return tx.Error

}

// get user by username
func (u *user_Repository) Login_Repo(username string) (*models.UserRequest, error) {
	dataRequest := models.UserRequest{}
	// แก้ query ตรงนี้ไหม ******
	tx := u.db.Table("user").Where("username=?", username).First(&dataRequest)
	return &dataRequest, tx.Error
}

func (u *user_Repository) GetDataByUsername_Repo(username string) (*models.UserRequest, error) {
	dataRequest := models.UserRequest{}
	tx := u.db.Table("user").Where("username=?", username).Find(&dataRequest)
	return &dataRequest, tx.Error
}
