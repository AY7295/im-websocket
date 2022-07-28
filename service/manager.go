package service

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"sync"
	"webSocket-be/config"
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

	receiver := c.Query("receiver")
	if receiver == "" {
		ErrorResponse(c, "receiver is empty")
		return
	}

	client, err := NewClient(receiver, c)
	if err != nil {
		ErrorResponse(c, "server err: "+err.Error())
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
				err := client.Socket.Close()
				if err != nil {
					config.Logfile.Println(fmt.Errorf("close socket connection err: %w", err))
				}
				close(client.Text)
			}

		case broadcast := <-manager.Broadcast:

			message, err := json.Marshal(broadcast.Message)
			if err != nil {
				config.Logfile.Println(fmt.Errorf("marshal message err: %w", err))
				continue
			}

			online := false
			// descp 在 当前对话框 直接发送
			if v, ok := manager.Hubs.Load(broadcast.Client.ReceiverId); ok && v.(*Client).ReceiverId == broadcast.Client.User.Id {
				v.(*Client).Text <- message
				online = true
			}

			// descp 不在 当前对话框 极光推送
			if !online {
				err = model.ZAddWithContext(broadcast.Client.ReceiverId, broadcast.Message)
				if err != nil {
					config.Logfile.Println(fmt.Errorf("ZAddWithContext failed err: %w", err))
				}
				err = model.NewJPush(broadcast.Message.User.Name, broadcast.Message.Text, []string{broadcast.Client.ReceiverId}).POST()
				if err != nil {
					config.Logfile.Println(fmt.Errorf("JPush failed err: %w", err))
				}
			}

		}
	}
}
