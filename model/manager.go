package model

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"sync"
	"webSocket-be/service"
)

type ClientManager struct {
	*sync.RWMutex
	Clients    map[string]*Client
	Unread     map[string][]DialogMessage
	Broadcast  chan *Broadcast
	Reply      chan *Client
	Register   chan *Client
	Unregister chan *Client
}

func NewManager() *ClientManager {
	return &ClientManager{
		Clients:    make(map[string]*Client),
		Broadcast:  make(chan *Broadcast),
		Reply:      make(chan *Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}

}

func (manager *ClientManager) WS(c *gin.Context) {

	token := c.GetHeader("Authorization")

	client, err := service.VerifyRequest(token, c)
	if err != nil {
		service.ErrorResponse(c, err.Error())
		return
	}

	manager.Register <- client

	go client.Read(manager)
	go client.Write(manager)

}

func (manager *ClientManager) Start() {
	for {
		select {
		case client := <-manager.Register:
			//manager.Lock() // 加锁
			manager.Clients[client.User.Id] = client
			messages := make([]DialogMessage, 0)
			// tip 发送未读消息
			for id, unread := range manager.Unread {
				if id == client.User.Id {
					for _, msg := range unread {
						messages = append(messages, msg)
					}
					text, err := json.Marshal(messages)
					if err != nil {
						log.Println("发送未读消息错误:", err)
						return
					}
					client.Text <- text
					delete(manager.Unread, id)
				}
			}
			//manager.Unlock() // 解锁

		case client := <-manager.Unregister:
			//manager.Lock() // 加锁
			_, ok := manager.Clients[client.User.Id]

			if ok {
				close(client.Text)
				delete(manager.Clients, client.User.Id)
			}
			//manager.Unlock() // 解锁

		case broadcast := <-manager.Broadcast:
			ToId := broadcast.Message.User.Id
			online := false

			message, err := json.Marshal(broadcast.Message)
			if err != nil {
				log.Println("json.Marshal error:", err)
				continue
			}

			// 在线直接发送()
			//manager.Lock() // 加锁
			for id, client := range manager.Clients {
				if id != ToId {
					continue
				}
				client.Text <- message
				online = true
			}

			// 不在线, 先存入未读消息
			if !online {
				manager.Unread[ToId] = append(manager.Unread[ToId], broadcast.Message)
			}
			//manager.Unlock() // 解锁
		}
	}
}
