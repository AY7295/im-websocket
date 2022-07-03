package service

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"webSocket-be/model"
)

func VerifyRequest(token string, c *gin.Context) (*model.Client, error) {

	user, err := verifyToken(token)
	if err != nil {
		return nil, err
	}

	conn, err := (&websocket.Upgrader{
		CheckOrigin: checkOrigin, // 允许跨域
	}).Upgrade(c.Writer, c.Request, nil) // 升级成ws协议
	if err != nil {
		return nil, err
	}

	client := &model.Client{
		Socket: conn,
		Text:   make(chan []byte),
		User:   *user,
	}

	return client, nil
}

func checkOrigin(r *http.Request) bool {

	return true
}
