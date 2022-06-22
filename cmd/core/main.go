package main

import (
	"broker/app"
	"broker/app/config"

	blogger "github.com/sirupsen/logrus"
)

func main() {
	if err := config.Init(); err != nil {
		blogger.Fatalln(err)
	}

	conf := config.GetConfig()
	blogger.Printf("read config %v", conf)

	app := app.NewApp(conf)
	app.Run()
}
