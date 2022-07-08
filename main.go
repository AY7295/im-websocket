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

	m := model.NewManager()

	go m.Start()

	err = route.NewRouter(m).Run(viper.GetString("utils.httpPort"))
	if err != nil {
		return
	}

}

func Init() error {
	viper.SetConfigName("config.json")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	err = service.InitGRPC()
	if err != nil {
		log.Println(err.Error())
		return err
	}

	err = model.InitRedis()
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil

}
