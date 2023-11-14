package util

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetParam(c *gin.Context, paramText string) (int, error) {
	param := c.Param(paramText)
	idProduct, err := strconv.Atoi(param)
	if err != nil {
		return 0, err
	}
	return idProduct, nil
}
