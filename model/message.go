package model

import (
	"time"
)

type DialogMessage struct {
	Id        string    `json:"_id"`       // 无实际意义
	Text      string    `json:"text"`      // 消息内容
	CreatedAt time.Time `json:"createdAt"` // 创建时间
	User      User      `json:"user"`      // 接收者
	Image     string    `json:"image"`     // 图片地址
	Sent      bool      `json:"sent"`      // 是否发送
	Received  bool      `json:"received"`  // 是否接收到
	Pending   bool      `json:"pending"`   // 正在发送中
}
