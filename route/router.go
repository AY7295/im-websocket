package route

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"webSocket-be/model"
)

func NewRouter(m *model.ClientManager) *gin.Engine {
	gin.SetMode(viper.GetString("utils.gin_mode"))
	e := gin.Default()

	e.GET("/ping", pong)

	e.GET("chat/ws", m.WS)

	return e
}

func pong(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
