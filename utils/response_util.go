package utils

import (
	"github.com/gin-gonic/gin"
)

func SuccessResponse(c *gin.Context, statusCode int, result interface{}) {
	c.JSON(statusCode, gin.H{
		"errorCode":    200,
		"errorMessage": "success",
		"result":       result,
	})
}

func ErrorResponse(c *gin.Context, statusCode int, errorCode int, message string) {
	c.JSON(statusCode, gin.H{
		"errorCode":    errorCode,
		"errorMessage": message,
		"result":       nil,
	})
}
