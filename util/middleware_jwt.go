package util

import (
	"fmt"
	"log"
	"myShopEcommerce/models"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	jwt "github.com/golang-jwt/jwt/v4"
)

// load signature
func loadSignature() string {
	// load config ENV
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Can not Load ENV File")
	}
	signature := os.Getenv("SIGNATURE")
	return signature
}

type customClaims struct {
	Role string `json:"role,omitempty"`
	jwt.StandardClaims
}

// ปั้น claims  payload
func claimsValue(userResponse *models.UserResponse) *customClaims {
	customClaims := customClaims{}
	// *****  Role สร้างขึ้นเอง *****
	customClaims.Role = userResponse.Role

	//******* jwt.StandardClaims **********
	// id
	customClaims.Id = strconv.Itoa(userResponse.Id)
	// เวลาหมดอายุ
	customClaims.ExpiresAt = time.Now().Add(30 * time.Minute).Unix()
	// น่าจะชื่อ user
	customClaims.Audience = userResponse.Username

	return &customClaims
}

func CreateToken(userRespnse *models.UserResponse) (string, error) {
	// load signature
	signature := loadSignature()

	// ปั้น claims   หรือ payload โดย userRespnse
	jwtClaims := claimsValue(userRespnse)

	jToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	_ = jToken

	// ********* รวม func แบบนี้เลยก็ได้ *******
	jToken2, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims).SignedString([]byte(signature))
	if err != nil {
		return "", err
	}
	_ = jToken2

	// ทำงานเหมือนกัน แต่จะใส payload ตรงนี้แทน
	jToken3 := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		// payload
		// "role": userRespnse.Role,
		// "jti":  strconv.Itoa(userRespnse.Id),
		// "exp":  time.Now().Add(30 * time.Minute).Unix(),
		// "aud":  userRespnse.Username,

		// payload
		"Role":      userRespnse.Role,
		"Id":        strconv.Itoa(userRespnse.Id),
		"ExpiresAt": time.Now().Add(30 * time.Minute).Unix(),
		"Audience":  userRespnse.Username,
	})
	_ = jToken3
	t, err3 := jToken3.SignedString([]byte(signature))
	_ = t
	_ = err3
	// fmt.Println("SignedString :",t)

	// ************  สรุป *************
	// 1.StandardClaims // ปั้น payload ปั้น StandardClaims (จะปั้นแยกหรือเอามาใส)
	// cliams := jwt.StandardClaims{
	// 	Issuer:    strconv.Itoa(user.Id),
	// 	ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // ระยะเวลา  เอาเวลาปัจจุบัน+เวลาอนาคตที่จะหมด เช่น อันนี้ +24 ช.ม. ****
	// }

	// 2.NewWithClaims(method, claims หรือ payload)
	// 3.SignedString([byte(signature)])

	// return Func นี้ออกไปเลยเพราะมี string, error
	// return jToken.SignedString([]byte(signature))

	// รวม func แบบนี้ก็ได้เพราะ signedString มันมี string and error ***
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims).SignedString([]byte(signature))

	// return t, err3
}

func VerifyToken(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		c.AbortWithStatusJSON(500, Error_Custom(500, "Login Please", "validater token"))
		return
	}
	tokenString := strings.TrimPrefix(token, "Bearer ")

	//****************** แบบ 1 ***************
	// 2อันนี้ใช้ต่างกันยังไง ใน youtube ชอบใช้อันนี้
	// tokenWithClaims, _ := jwt.ParseWithClaims(tokenString, &customClaims{}, func(t *jwt.Token) (interface{}, error) {
	// 	_, ok := t.Method.(*jwt.SigningMethodHMAC)
	// 	if !ok {
	// 		return nil, fmt.Errorf("Unexpected signing method : %v", t.Header["alg"])
	// 	}
	// 	return []byte(loadSignature()), nil
	// })

	// ถ้าเช็ค errv ตรงนี้ ก็ได้
	// if errParseWithClaims != nil {
	// 	c.AbortWithStatusJSON(404, gin.H{
	// 		"error": "token Unexpected (ParseWithClaims)",
	// 	})
	// 	return
	// }

	// ทำตาม youtube ศึกษาดูว่ามันทำงานยังไง ****
	// ต้องใช้ jwt.ParseWithClaim

	// claims, ok := tokenWithClaims.Claims.(*customClaims)
	// if !ok {
	// 	c.AbortWithStatusJSON(500, gin.H{
	// 		"message": "token is invalid",
	// 	})
	// }
	// // เช็คเรื่อง ระยะเวลา ว่าหมดรึยัง
	// if claims.ExpiresAt < time.Now().Local().Unix() {
	// 	c.AbortWithStatusJSON(500, gin.H{
	// 		"message": "token is already expired",
	// 	})
	// }

	// ทำแบบนี้ต้องทำตาม youtube ถึงจะได้ เพราะ tokens.Claims.(*customClaims)********
	// c.Set("jti", claims.Id)

	// *****************************************************
	//*************** แบบ 2 *****************
	tokens, errParse := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		// เช็ค เช็ค  method jwt .ใช้ SigningMethodHMAC รึป่าว ***
		// func นี้แยกออกมาได้
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method : %v", t.Header["alg"])
		}
		// ถ้าใช้จะคืน ลายเซ็นออกไป เพื่อเอาไปถอดรหัส
		return []byte(loadSignature()), nil
	})

	if errParse != nil {
		// c.AbortWithStatusJSON(404, gin.H{
		// 	// "error": errParse.Error(),// ผิด กับ หมดเวลา
		// 	// "error": Error_Custom(404, errParse.Error(), "check Token"), // ผิด กับ หมดเวลา
		// })
		c.AbortWithStatusJSON(500, Error_Custom(500, errParse.Error(), "validater token"))
		return
	}
	_ = tokens
	// fmt.Println("Token:::::::",tokens.Claims)
	// fmt.Println(tokens.Claims) // aud: / exp: ระยะเวลา/ ค่าที่ใสเพิ่มจะอยู่ในนี้ด้วย
	// fmt.Println(tokens.Header) // alg:HS256 / typ:JWT
	// fmt.Println(tokens.Method.(jwt.SigningMethod)) // &{HS256 SHA-256}
	// fmt.Println(tokens.Raw) // token
	// fmt.Println(tokens.Signature) // signature ที่เป็น []byte
	// fmt.Println(tokens.Valid) // true
	// fmt.Println("Header alg",tokens.Header["alg"])
	// fmt.Println("Header typ",tokens.Header["typ"])

	// tokens.Valid เช็ค ค่าต่างๆใน token
	// c,ok:= tokens.Claims.(jwt.MapClaims)
	// if ok {
	// 	fmt.Println("",c["aud"])
	// }

	//******************** Set แบบชื่อ ย่อ ************
	if claims_, ok := tokens.Claims.(jwt.MapClaims); ok && tokens.Valid {
		// ดึงค่าและ เซตค่า aud
		c.Set("aud", claims_["aud"]) // aud  เป็นค่า json ของ StandardClaims

		// ดึงค่าจาก token
		role := claims_["role"] // role  เป็นค่า json ของ customClaims เราสร้างขึ้นมาเอง
		// เซตค่าที่ดึงได้
		c.Set("role", role)

		// ดึงค่าและ เซตค่า jit
		c.Set("jti", claims_["jti"]) // jti  เป็นค่า json ของ StandardClaims
	}

	//******************** Set แบบชื่อ เต็ม ************
	// if claims_, ok := tokens.Claims.(jwt.MapClaims); ok && tokens.Valid {
	// 	// ดึงค่าและ เซตค่า aud
	// 	c.Set("Audience", claims_["Audience"]) // aud  เป็นค่า json ของ StandardClaims

	// 	// ดึงค่าจาก token
	// 	role := claims_["Role"] // role  เป็นค่า json ของ customClaims เราสร้างขึ้นมาเอง
	// 	// เซตค่าที่ดึงได้
	// 	c.Set("Role", role)

	// 	// ดึงค่าและ เซตค่า jit
	// 	c.Set("Id", claims_["Id"]) // jti  เป็นค่า json ของ StandardClaims
	// }

	c.Next()

}
