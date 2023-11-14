package util

import (
	_"fmt"
	_ "myShopEcommerce/models"
	"strings"
	"unicode"

	"regexp"

	"github.com/go-playground/validator"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	validate.RegisterValidation("notblank", NotBlank)
	validate.RegisterValidation("haveSpace", CheckSpace)
	validate.RegisterValidation("checkScript", CheckScriptJava)
}

func ValidateDataUser(dataUser interface{}) error {
	err := validate.Struct(dataUser)
	if err != nil {
		return err
	}
	return nil
}

// custom Validate Not Blank
func NotBlank(fl validator.FieldLevel) bool {
	return strings.TrimSpace(fl.Field().String()) != ""
}

// check space
func CheckSpace(value validator.FieldLevel) bool {
	for _, v := range value.Field().String() {
		if unicode.IsSpace(v) {
			// fmt.Println("Is Space")
			return false
		}
	}
	return true
}

func CheckScriptJava(value validator.FieldLevel) bool {
	// regex
	pattern := `<script[^>]*>([\s\S]*?)<\/script>`
	result := regexp.MustCompile(pattern)

	// text:= `<script> let x;x = 6;
	// document.getElementById("demo").innerHTML = x;
	// </script>`

	// ถ้ามีจะได้ true ถ้า reture จะติด error
	//ถ้าไม่มีจะได้ false
	// rr := result.Match([]byte(value.Field().String()))
	// fmt.Println(rr)

	return !result.Match([]byte(value.Field().String()))
}

func insertErrorText(text []string) []string {
	var errLoop []string
	for _, value := range text {
		errLoop = append(errLoop, value)
	}
	return errLoop
}

// type dataUserType interface {
// 	models.UserRequest | models.UserUpdate | models.UserLogin | models.Store
// }

// func ValidateDataUser2[T dataUserType](dataUser *T) error {
// 	err := validate.Struct(dataUser)
// 	return err
// }

// ต้องรับค่า path มาใสด้วยว่า error จาก path ไหน
// เอาชื่อ struct มาใส น่าจะได้ + string เข้าไป
// func ValidateDataUser(dataUser interface{}) error {
// func ValidateDataUser[T dataUserType](dataUser interface{}) error {
// 	err := validate.Struct(dataUser)
// 	if err != nil {
// 		return Error_Custom(500, err.Error(), "Validate_Data")
// 	}
// 	return nil
// }
