package service

import (
	"errors"
	"fmt"
	"myShopEcommerce/models"
	"myShopEcommerce/repository"
	"myShopEcommerce/util"
	"time"
)

type user_Service struct {
	userRepo repository.IUser_Repository
}

type IUser_Sevice interface {
	CreateUser_Service(userData *models.UserRequest) (userResponse *models.UserResponse, err error)
	GetAllUser_Service(pagination *models.Pagination) (userResponse []models.UserResponse, count_User int, err error)
	UpdateUser_Service(id_User int, userUpdate *models.UserUpdate) error
	Login_Service(dataLogin *models.UserLogin) (userResponse *models.UserResponse, err error)
	GetDataByUsername_Service(userName string) (userResponse *models.UserResponse, err error)
	ValidateData_serviice(dataValidate interface{}) error
}

func NewUser_Service(userRepo repository.IUser_Repository) IUser_Sevice {
	return &user_Service{userRepo: userRepo}
}

// Validate data struct
func (us *user_Service) ValidateData_serviice(dataValidate interface{}) error {
	return util.ValidateDataUser(dataValidate)
}

// CreateUserService implements IUser_Sevice
func (us *user_Service) CreateUser_Service(userData *models.UserRequest) (*models.UserResponse, error) {

	dataResponse := models.UserResponse{}
	// hash password
	hasPasswword, err := util.HasPasword(userData.Password)
	if err != nil {
		return &dataResponse, err
	}

	userData.Password = hasPasswword
	// ตรงนี้เาต้องใสเองไหม ***********
	userData.CreatedAt = time.Now()

	// to Repo
	userRepo, errCreate := us.userRepo.CreateUser_Repo(userData)

	// ปั้น data response ใหม่
	newUser := models.UserResponse{
		Username:    userRepo.Username,
		Firstname:   userRepo.Firstname,
		Lastname:    userRepo.Lastname,
		Phonenumber: userRepo.Phonenumber,
		// ******* ใสค่าซ้อนใน struct *******
		// Address: models.Address{
		// 	Address: userRepo.Address.Address,
		// 	City:    userRepo.Address.City,
		// 	Zipcode: userRepo.Address.Zipcode,
		// },

		// ***** ใสค่าใน struct ที่ซ้อนกันแบบนี้ก็ได้ *****
		Address: userRepo.Address,
	}

	return &newUser, errCreate
}

// GetDataUserService implements IUser_Sevice
func (us *user_Service) GetAllUser_Service(page *models.Pagination) ([]models.UserResponse, int, error) {
	users := []models.UserResponse{}
	// to Repo
	dataUser, count, errGetAllUser := us.userRepo.GetAllUser_Repo(page)
	if errGetAllUser != nil {
		return users, 0, errGetAllUser
	}

	// loop data user
	// ตรงนี้มีวิธีที่ไม่ต้อง loop ไหม
	for _, value := range dataUser {
		users = append(users, models.UserResponse{
			Id:          value.Id,
			Username:    value.Username,
			Firstname:   value.Firstname,
			Lastname:    value.Lastname,
			Phonenumber: value.Phonenumber,
			Address:     value.Address,
		})
	}
	return users, count, nil
}

// UpdateUserService implements IUser_Sevice
func (us *user_Service) UpdateUser_Service(id_User int, newData *models.UserUpdate) error {
	// หาข้อมูลที่จะ update ก่อน
	// update new data
	errUpdateUser := us.userRepo.UpdateUser_Repo(id_User, newData)
	if errUpdateUser != nil {
		return errUpdateUser
	}
	return nil
}

func (us *user_Service) Login_Service(dataLogin *models.UserLogin) (*models.UserResponse, error) {
	// dataRepo := models.UserRequest{}

	dataResponse := models.UserResponse{}

	// เช็คว่ามี user ไหม
	dataRepo, err := us.userRepo.Login_Repo(dataLogin.Username)
	if err != nil {
		// fmt.Errorf มันคือ error 
		// ตัวอักษรใน error ไม่ควรใช้ตัวพิมพ์ใหญ่ มันจะ error สีเหลือง ****
		return &dataResponse, fmt.Errorf("username, password is incorrect")
		// return &dataResponse, err
	}

	passwordSalt := fmt.Sprintf(dataLogin.Password + "Russy")

	//***************** CompairPassword ******
	// ส่ง password in repo  กับ password ที่ login เข้ามา ไปเทียบกัน
	err = util.CompairPassword(dataRepo.Password, passwordSalt)
	if err != nil {
		return &dataResponse, errors.New("username, password is incorrect")
	}

	// ปั้นข้อมูลใหม่
	dataResponse = models.UserResponse{
		Id:          dataRepo.Id,
		Username:    dataRepo.Username,
		Firstname:   dataRepo.Firstname,
		Lastname:    dataRepo.Lastname,
		Phonenumber: dataRepo.Phonenumber,
		Role:        dataRepo.Role,
		Address:     dataRepo.Address,
	}

	return &dataResponse, nil
}

func (us *user_Service) GetDataByUsername_Service(username string) (*models.UserResponse, error) {

	dataRepo, errGetByUsername := us.userRepo.GetDataByUsername_Repo(username)

	dataResponse := models.UserResponse{
		Id:          dataRepo.Id,
		Username:    dataRepo.Username,
		Firstname:   dataRepo.Firstname,
		Lastname:    dataRepo.Lastname,
		Phonenumber: dataRepo.Phonenumber,
		Role:        dataRepo.Role,
		Address:     dataRepo.Address,
	}

	return &dataResponse, errGetByUsername
}

// ทำอะไรได้บ้าง
// type dataType interface {
// 	models.UserRequest | models.UserUpdate | models.UserLogin | models.Store
// }
