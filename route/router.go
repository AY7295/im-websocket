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
	token := c.Query("token")
	if token == "" {
		service.SuccessResponse(c, "pong")
	}

	user, err := service.VerifyToken(token)
	if err != nil {
		service.ErrorResponse(c, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"message": user,
	})
}
