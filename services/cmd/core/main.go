package main

import (
	"broker/core"
	"broker/core/config"
	"broker/pkg/logger"
)

func main() {
	logger.Init()
	
	if err := config.Init(); err != nil {
		logger.Logger.Fatalln(err)
	}

	conf := config.GetConfig()
	logger.Logger.Info("read config %v", conf)

	app := app.NewApp(conf)
	app.Run()
}
