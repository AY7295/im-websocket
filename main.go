package main

import (
	"github.com/spf13/viper"
	"log"
	"webSocket-be/model"
	"webSocket-be/route"
	"webSocket-be/service"
)

func main() {

	err := Init()
	if err != nil {
		log.Println(err.Error())
		panic(err)
	}

	m := service.NewManager()

	go m.Start()

	err = route.NewRouter(m).Run(viper.GetString("utils.httpPort"))
	if err != nil {
		return
	}

}

func Init() error {
	viper.SetConfigType("json")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	err = model.InitRedis()
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil

}

//test for JPush

//broadcast := &service.Broadcast{
//	Message: model.DialogMessage{
//		Text: "test for content",
//		User: model.User{
//			Id:   "8008121377",
//			Name: "yzx",
//		},
//	},
//}
//
//err = model.NewJPush(broadcast.Message.User.Name, broadcast.Message.Text, []string{broadcast.Message.User.Id}).POST()
//if err != nil {
//	log.Println("JPush error:", err)
//}
