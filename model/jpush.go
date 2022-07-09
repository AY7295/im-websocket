package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"log"
	"net/http"
)

const (
	titleFormat = "您有一条来自 %s 的新消息"
	msgFormat   = "内容: %s "
)

type JPush struct {
	Platform     []string     `json:"platform"`
	Audience     Audience     `json:"audience"`
	Options      Options      `json:"options"`
	Notification Notification `json:"notification"`
}

type Audience struct {
	Alias []string `json:"alias"`
}

type Options struct {
	ApnsProduction bool `json:"apns_production"`
}

type Notification struct {
	Android Android `json:"android"`
	Ios     Ios     `json:"ios"`
}

type Android struct {
	Title  string            `json:"title"`
	Alert  string            `json:"alert"`
	Extras map[string]string `json:"extras"`
}

type Ios struct {
	IAlert IAlert            `json:"alert"`
	Extras map[string]string `json:"extras"`
}

type IAlert struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func NewJPush(title, msg string, ids []string, m ...map[string]string) *JPush {
	j := &JPush{}
	j.Platform = viper.GetStringSlice("j_push.platform")

	j.Audience.Alias = ids

	j.Options.ApnsProduction = viper.GetBool("j_push.is_production_env")

	j.Notification.Android.Title = fmt.Sprintf(titleFormat, title)
	j.Notification.Android.Alert = fmt.Sprintf(msgFormat, msg)
	j.Notification.Ios.IAlert.Title = fmt.Sprintf(titleFormat, title)
	j.Notification.Ios.IAlert.Body = fmt.Sprintf(msgFormat, msg)
	if len(m) < 1 {
		j.Notification.Android.Extras = make(map[string]string)
		j.Notification.Ios.Extras = make(map[string]string)
		return j
	}

	j.Notification.Android.Extras = m[0]
	j.Notification.Ios.Extras = m[0]

	return j
}

func (j *JPush) POST() error {
	body, err := json.Marshal(j)
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", viper.GetString("j_push.url"), bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Basic "+viper.GetString("j_push.token"))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err1 := Body.Close()
		if err1 != nil {
			log.Println(err1.Error())
			return
		}
	}(resp.Body)

	return nil

}
