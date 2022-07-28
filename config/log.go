package config

import (
	"log"
	"os"
)

var Logfile *log.Logger

func InitLogger(path string) {

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 666)
	if err != nil {
		panic(err)
	}

	Logfile = log.New(file, "", log.LstdFlags|log.Llongfile) //ErrLog会收集err,
	return
}
