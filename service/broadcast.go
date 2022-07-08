package service

import "webSocket-be/model"

type Broadcast struct {
	Client  *Client
	Message model.DialogMessage
}
