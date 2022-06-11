package main

import (
	"broker/app"
	"broker/app/config"
	"log"
)

func main() {
	if err := config.Init(); err != nil {
		log.Fatalln(err)
	}

	conf := config.GetConfig()
	log.Printf("read config %v", conf)

	app := app.NewApp(conf)
	app.Run()
}
