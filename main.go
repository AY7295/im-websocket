package main

import (
	"fmt"
	"github.com/spf13/viper"
	"webSocket-be/config"
	"webSocket-be/model"
	"webSocket-be/route"
	"webSocket-be/service"
)

func main() {

	err := Init()
	if err != nil {
		config.Logfile.Println(err)
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
		return fmt.Errorf("viper read config err: %w", err)
	}

	config.InitLogger("./config/err.txt")

	err = model.InitRedis()
	if err != nil {
		return fmt.Errorf("init redis err: %w", err)
	}

	return nil

}
