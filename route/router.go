package route

import (
	"github.com/gin-gonic/gin"
)

func NewRouter(h *Handler) *gin.Engine {
	gin.SetMode(h.Config.GinMode)
	e := gin.Default()

	e.GET("/ping", pong)

	e.GET("chat/ws", h.WS)

	return e
}

func pong(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
