package service

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"sync"
	"webSocket-be/model"
)

type ClientManager struct {
	Hubs       *sync.Map
	Broadcast  chan *Broadcast
	Register   chan *Client
	Unregister chan *Client
}

func NewManager() *ClientManager {
	return &ClientManager{
		Hubs:       &sync.Map{},
		Broadcast:  make(chan *Broadcast),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}

}

func (manager *ClientManager) WS(c *gin.Context) {

	token := c.GetHeader("Authorization")
	if token == "" {
		ErrorResponse(c, "token is empty")
	}

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

			manager.Hubs.Store(client.User.Id, client)

		case client := <-manager.Unregister:

			if _, ok := manager.Hubs.Load(client.User.Id); ok {
				manager.Hubs.Delete(client.User.Id)
				close(client.Text)
			} else {
				log.Println(client.User.Id + " is not exist in Hubs, delete failed")
			}

		case broadcast := <-manager.Broadcast:

			message, err := json.Marshal(broadcast.Message)
			if err != nil {
				log.Println(broadcast.Client.User.Id+" json.Marshal error:", err)
				continue
			}

			online := false
			// descp 在 当前对话框 直接发送
			if v, ok := manager.Hubs.Load(broadcast.Message.User.Id); ok {
				v.(Client).Text <- message
				online = true
			}

			// descp 不在 当前对话框 极光推送
			if !online {
				err = model.ZAddWithContext(broadcast.Message.User.Id, broadcast.Message)
				if err != nil {
					log.Println("ZAddWithContext error:", err)
				}
				err = model.NewJPush(broadcast.Message.User.Name, broadcast.Message.Text, []string{broadcast.Message.User.Id}).POST()
				if err != nil {
					log.Println("JPush error:", err)
				}
			}

		}
	}
}
