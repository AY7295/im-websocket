package service

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"webSocket-be/config"
	"webSocket-be/model"
)

type Client struct {
	Socket     *websocket.Conn
	Text       chan []byte
	User       *model.User // 发送者信息
	ReceiverId string      // 接收者id
}

func (c *Client) Write(manager *ClientManager) {
	defer c.Unregister(manager)
	for {
		c.Socket.PongHandler()
		select {
		case message := <-c.Text:
			err := c.Socket.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				config.Logfile.Println(fmt.Errorf("user %s send message err: %w", c.User.Id, err))
				return
			}
		}
	}
}

func (c *Client) Read(manager *ClientManager) {
	defer c.Unregister(manager)
	for {
		_, msg, err := c.Socket.ReadMessage()
		if err != nil {
			config.Logfile.Println(fmt.Errorf("user %s; read message err: %w", c.User.Id, err))
			return
		}

		message := model.DialogMessage{}
		err = json.Unmarshal(msg, &message)
		if err != nil {
			config.Logfile.Println(fmt.Errorf("user %s; unmarshal message err: %w", c.User.Id, err))
			return
		}
		manager.Broadcast <- &Broadcast{
			Client:  c,
			Message: message,
		}
	}
}

func (c *Client) Unregister(manager *ClientManager) {
	manager.Unregister <- c
}
