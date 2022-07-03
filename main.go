package main

import (
	"log"
	"webSocket-be/route"
	"webSocket-be/service"
)

func main() {

	h, err := route.NewHandler("config/config.json")
	if err != nil {
		log.Println(err.Error())
		panic(err)
	}
	err = service.InitGRPC(h.Config)
	if err != nil {
		log.Println(err.Error())
		panic(err)
	}

	go h.Manager.Start()

	e := route.NewRouter(h)

	err = e.Run(h.Config.HttpPort)
	if err != nil {
		return
	}

}
