package service

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"webSocket-be/model"
)

type Client struct {
	Socket *websocket.Conn
	Text   chan []byte
	User   model.User // 发送者信息
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
			err = c.Socket.WriteMessage(websocket.CloseMessage, []byte("读取消息错误"))
			if err != nil {
				log.Println("发送消息错误:", err)
				return
			}
			return
		}

		message := model.DialogMessage{}
		err = json.Unmarshal(msg, &message)
		if err != nil {
			log.Println("发送消息错误:", err)
			msg = []byte("解析消息错误")
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
	err := c.Socket.Close()
	if err != nil {
		log.Println("ID: "+c.User.Id+"	关闭连接错误:", err)
	}
}
