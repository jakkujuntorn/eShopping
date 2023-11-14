package util

import (
	"fmt"
	// "log"
	// "os"
	// "strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv"

	jwt "github.com/golang-jwt/jwt/v4"
)

type JwtClaims struct {
	ComapnyId string
	Username  string
	Roles     []int
	// VerifyExpiresAt, VerifyIssuer อยู่ใน StandardClaims 
	jwt.StandardClaims
}

// เพิ่ม Func Valid ก็ confrom ตาม claims ได้แล้ว
// func (jw *JwtClaims)Valid() error {
// 	return nil
// }

const ip = "192.168.0.107"

// ทำงานอย่างไร
func (claims JwtClaims) Valid2() error {
	var now = time.Now().UTC().Unix()
	if claims.VerifyExpiresAt(now, true) && claims.VerifyIssuer(ip, true) {
		return nil
	}
	return fmt.Errorf("Token is invalid")
}

func GenrateToken(claims *JwtClaims, expirationTime time.Time) (string, error) {
	//*** Jwt StandardClaims **
	// set ค่า claims
	claims.ExpiresAt = expirationTime.Unix()
	claims.IssuedAt = time.Now().UTC().Unix()
	claims.Issuer = ip

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenSign, err := token.SignedString([]byte("signature"))
	if err != nil {
		return "", err
	}

	return tokenSign, nil
}

func Login(c *gin.Context) {
	// gin ทำงาน...
	// ถอดข้อมูล ShouldBindJSON()...

	//่ jwt
	// set ค่า StandardClaims
	var claims = JwtClaims{}
	claims.ComapnyId = "ComapnyId"
	claims.Username = "dataUser.Username"
	claims.Roles = []int{1, 2, 3}
	claims.Audience = c.Request.Header.Get("audience")

	var tokenCreateTime = time.Now().UTC()
	var expirationTime = tokenCreateTime.Add(time.Duration(10) * time.Minute)

	tokenString, err := GenrateToken(&claims, expirationTime)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{
			"Message": "Error is Gerating token",
		})
	}
	c.JSON(200, gin.H{
		"Message": "Login Success",
		"token":   tokenString,
	})

}
