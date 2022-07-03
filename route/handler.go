package route

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"webSocket-be/config"
	"webSocket-be/model"
	"webSocket-be/service"
)

type Handler struct {
	Config  config.Config        `json:"config"`
	Manager *model.ClientManager `json:"-"`
}

func (h *Handler) WS(c *gin.Context) {

	token := c.GetHeader("Authorization")

	client, err := service.VerifyRequest(token, c)
	if err != nil {
		service.ErrorResponse(c, err.Error())
		return
	}

	h.Manager.Register <- client

	go client.Read(h.Manager)
	go client.Write(h.Manager)

}

func NewHandler(path string) (*Handler, error) {
	h := &Handler{}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(file, h)
	if err != nil {
		return nil, err
	}
	h.Manager = &model.ClientManager{
		Clients:    make(map[string]*model.Client),
		Broadcast:  make(chan *model.Broadcast),
		Reply:      make(chan *model.Client),
		Register:   make(chan *model.Client),
		Unregister: make(chan *model.Client),
	}
	return h, nil
}
