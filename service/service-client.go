package service

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

func VerifyRequest(receiverId, token string, c *gin.Context) (*Client, error) {

	user, err := VerifyToken(token)
	if err != nil {
		return nil, err
	}

	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}).Upgrade(c.Writer, c.Request, nil) // 升级成ws协议
	if err != nil {
		config.Logfile.Println(fmt.Errorf("setup websocket connection err: %w", err))
		http.NotFound(c.Writer, c.Request)
		return nil, err
	}

	return &Client{
		Socket:     conn,
		Text:       make(chan []byte),
		User:       user.(*model.User),
		ReceiverId: receiverId,
	}, nil

}
