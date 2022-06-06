package main

import (
	"broker/app/api"
	"broker/app/api/user"
	"broker/app/config"
	"broker/app/repository"
	"broker/app/services"
	"broker/pkg/pg"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi"
)

func main() {
	ctx := context.Background()
	if err := config.Init(); err != nil {
		log.Fatalln(err)
	}

	conf := config.GetConfig()
	log.Printf("read config %v", conf)

	pgConn, err := pg.NewReadAndWriteConnection(ctx, conf.Database, conf.Database, nil)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Database connection established")
	web := api.NewAPIServer(":80")

	authService := services.NewAuthService([]byte(conf.JWT.SigningKey), conf.JWT.ExpireDuration)
	userRepo := repository.NewUserRepository(pgConn)
	userController := user.NewController(authService, userRepo)
	web.Router().Route("/api/v1", func(v1 chi.Router) {
		user.RegisterRouter(v1, userController)
	})

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
