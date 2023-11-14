package util

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// ******* Error *****
type error_Data struct {
	StatusCode int
	Code       int
	Message    []string
	Path       string
}

type Error_Response struct {
	Success   string
	ErrorData error_Data
	// confrom ตาม interface error error_Response เลยเป็น
	error
}

func Handler_Error(c *gin.Context, err error) {
	// ใช้ switch *******************
	switch e := err.(type) {
	case Error_Response:
		c.AbortWithStatusJSON(e.ErrorData.StatusCode, e)
	case error:
		c.AbortWithStatusJSON(500, gin.H{
			"data": e.Error(),
		})
	}

	// ใช้  if ***************
	// msg,ok := errUpdate.(util.Error_Response)
	// if ok {
	// 	fmt.Printf("Error REsponse: %v",msg.ErrorData)
	// 	c.AbortWithStatusJSON(404,gin.H{
	// 		"data":msg,
	// 	})
	// 	return
	// }
}

// Error_Custom return Error_Response
func Error_Custom(ststusCode int, message string, path string) Error_Response {

	// ทำ error array ตรงนี้ ***
	messSlice:= strings.Split(message,"\n")
	
	return Error_Response{
		Success: "false",
		ErrorData: error_Data{
			StatusCode: ststusCode,
			Code:       1150,
			Message:    messSlice, //  เป็น array
			Path:       path,
		},
	}
}
