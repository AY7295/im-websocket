package route

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"webSocket-be/service"
)

func NewRouter(m *service.ClientManager) *gin.Engine {
	gin.SetMode(viper.GetString("utils.gin_mode"))
	e := gin.Default()

	e.GET("/ping", pong)

	e.GET("chat", m.WS)

	return e
}

func pong(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
