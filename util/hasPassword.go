package util

import (
	_"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HasPasword(password string) (string, error) {
	// saltPassword:= fmt.Sprintf(password+"Russy")
	
	bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(password+"Russy"), bcrypt.DefaultCost)
	return string(bcryptPassword), err
}

func CompairPassword(hash_Password, user_Password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash_Password), []byte(user_Password))
	return err
}
