package service

import "github.com/gin-gonic/gin"

func ErrorResponse(c *gin.Context, data any) {
	c.JSON(400, gin.H{
		"message": data,
	})
}

func SuccessResponse(c *gin.Context, data any) {
	c.JSON(200, gin.H{
		"message": data,
	})
}
