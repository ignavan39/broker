package app

import (
	"broker/app/config"
	delivery "broker/app/delivery/http"
	"broker/app/delivery/http/auth/v1"
	"broker/app/repository"
	"broker/pkg/pg"
	authSrv "broker/app/service/auth"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi"
)

type App struct {
	config config.Config
	web    *delivery.Server
}

func NewApp(config config.Config) *App {
	a := &App{
		config: config,
	}
	ctx := context.Background()

	pgConn, err := pg.NewReadAndWriteConnection(ctx, a.config.Database, a.config.Database, nil)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Database connection established")
	a.web = delivery.NewAPIServer(":80").WithCors()

	authService := authSrv.NewAuthService([]byte(a.config.JWT.SigningKey), a.config.JWT.ExpireDuration)
	userRepo := repository.NewUserRepository(pgConn)
	userController := auth.NewController(authService, userRepo)

	userRouter := auth.NewRouter(userController)
	a.web.Router().Route("/api/v1", func(v1 chi.Router) {
		userRouter.InitRoutes(v1)
	})
	return a
}

func (a *App) Run() {
	if err := a.web.Start(); err != nil {
		log.Fatal(err)
	}
	appCloser := make(chan os.Signal)
	signal.Notify(appCloser, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-appCloser
		log.Println("[os.SIGNAL] close request")
		go a.web.Stop()
		log.Println("[os.SIGNAL] done")
	}()
	a.web.WaitForDone()
}
