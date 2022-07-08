package model

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"webSocket-be/model"
)

type ClientManager struct {
	*sync.RWMutex
	Clients    map[string]*Client
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

	client, err := VerifyRequest(token, c)
	if err != nil {
		ErrorResponse(c, err.Error())
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

			manager.Clients[client.User.Id] = client

		case client := <-manager.Unregister:

			_, ok := manager.Clients[client.User.Id]

			if ok {
				close(client.Text)
				delete(manager.Clients, client.User.Id)
			}

		case broadcast := <-manager.Broadcast:
			ToId := broadcast.Message.User.Id
			online := false

			message, err := json.Marshal(broadcast.Message)
			if err != nil {
				log.Println("json.Marshal error:", err)
				continue
			}

			// 在 当前对话框 直接发送()
			for id, client := range manager.Clients {
				if id != ToId {
					continue
				}
				client.Text <- message
				online = true
			}

			if !online {
				manager.Unread[ToId] = append(manager.Unread[ToId], broadcast.Message)
			}
			//manager.Unlock() // 解锁
		}
	}
}