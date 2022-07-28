package service

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"webSocket-be/config"
	"webSocket-be/model"
)

func NewClient(receiverId string, c *gin.Context) (*Client, error) {

	user, ok := c.Get("user")
	if !ok {
		return nil, errors.New("bad Authorization")
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
