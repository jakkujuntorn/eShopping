package util

import (
	"fmt"

	_ "errors"
	"github.com/gin-gonic/gin"
	"strings"
)

func VerifyRole(c *gin.Context) {
	// เอา username ไปหาข้อมูลใน db ว่าเป็น admin รึ ป่าว ***
	// ระหว่างดึงจาก gin กับ เอา token มาถอดใหม่อันไหนจะดีกว่ากัน ********
	// ใน youtube เอา token มาถิดใหม่ แล้วเช็คเวลาก่อนดึง สิทธ์จาก db ******
	role, _ := c.Get("role")

	// แปลง เป็น string
	string_Role, ok := role.(string)
	if !ok {
		fmt.Println("role not string")
	}

	// เอา role in token มาเช็ค สิทธิ์
	if strings.ToLower(string_Role) == "admin" {
		c.Next()
	} else {
		// ต้องใช้ AbortWithStatusJSON ไม่งั้นมันไม่หยุดถ้ามี error
		c.AbortWithStatusJSON(403, Error_Custom(403, "Forbidden Error", "validater token"))

		// c.JSON(500, Error_Custom(403, "Forbidden Error", "validater token"))
		// return
	}

}
