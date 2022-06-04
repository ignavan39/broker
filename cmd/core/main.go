package main

import (
	"broker/app/api"
	"broker/app/config"
	"broker/pkg/pg"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)


func main(){
	ctx := context.Background()
	if err := config.Init(); err != nil {
		log.Fatalln(err)
	}

	conf := config.GetConfig()
	log.Printf("read config %v",conf)
	
	pgConn,err := pg.NewReadAndWriteConnection(ctx,conf.Database,conf.Database,nil)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Database connection established")
	
	// TODO remove after using
	pgConn.Write()

	web := api.NewAPIServer(fmt.Sprintf("%d",conf.Port))
	if err := web.Start(); err != nil {
		log.Fatal(err)
	}

	appCloser := make(chan os.Signal)
	signal.Notify(appCloser, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-appCloser
		log.Println("[os.SIGNAL] close request")
		go web.Stop()
		log.Println("[os.SIGNAL] done")
	}()
	web.WaitForDone()
}