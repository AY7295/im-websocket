package route

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"webSocket-be/config"
	"webSocket-be/service"
)

func NewRouter(m *service.ClientManager) *gin.Engine {
	gin.SetMode(viper.GetString("utils.gin_mode"))

	e := gin.Default()
	e.Use(auth)

	e.GET("chat", m.WS)

	return e
}

func auth(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		service.ErrorResponse(c, "receiver is empty")
	}
	user, err := service.VerifyToken(token)
	if err != nil {
		config.Logfile.Println(err)
		service.ErrorResponse(c, "bad Authorization err: "+err.Error())
		return
	}

	c.Set("user", user)
	c.Next()

}
