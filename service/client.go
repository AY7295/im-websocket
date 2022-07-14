package service

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"webSocket-be/model"
)

type Client struct {
	Socket     *websocket.Conn
	Text       chan []byte
	User       model.User // 发送者信息
	ReceiverId string     // 接收者id
}

func (c *Client) Read(manager *ClientManager) {
	defer c.Unregister(manager)
	for {
		c.Socket.PongHandler()
		select {
		case message := <-c.Text:
			err := c.Socket.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Println("发送消息错误:", err)
				return
			}
		}
	}
}

func (c *Client) Write(manager *ClientManager) {
	defer c.Unregister(manager)
	for {
		_, msg, err := c.Socket.ReadMessage()
		if err != nil {
			log.Println("读取消息错误:", err)
			return
		}

		message := model.DialogMessage{}
		err = json.Unmarshal(msg, &message)
		if err != nil {
			log.Println("反序列化消息失败: ", err)
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
