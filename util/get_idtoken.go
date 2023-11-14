package util

import (
	"errors"
	_ "fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetUserToken(c *gin.Context) (string, error) {
	aud, ok := c.Get("aud") // ดึงค่า username
	if ok {
		aud_String, _ := aud.(string)
		return aud_String, nil
	} else {
		// return "", Error_Custom(401, "You not login", "user_handler_GetUserByName")
		return "", errors.New("can not get token")
	}
}

func GetidToken(c *gin.Context) (int, error) {
	//********** ใช้ชื่อแบบย่อ *************
	idToken, _ := c.Get("jti")
	// aud, _ := c.Get("aud")

	//************ ใช้ชื่อแบบเต็ม **************
	// idToken, errGet := c.Get("Id")
	// if errGet {
	// รับ token ได้ให้แปลงค่า tokenFF
	// } else {
	// 	return 0, errors.New("can not get token")
	// 	// return 0, Error_Custom(401, "You not login", "user_handler_GetUserByName")
	// }

	// เช็ค token ว่าเป็น string รึป่าว
	// ถ้ามันไม่ ok จะข้ามไป
	idString, ok := idToken.(string)
	if ok {
		// แปลงค่าได้
		id, errStrconv := strconv.Atoi(idString)
		if errStrconv != nil {
			return 0, errStrconv
		} else {
			return id, nil
		}
	} else {
		// แปลค่าไม่ได้
		return 0, errors.New("cannot change token to string")
	}

}
