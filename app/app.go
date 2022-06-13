package app

import (
	"broker/app/config"
	delivery "broker/app/delivery/http"
	"broker/app/delivery/http/auth/v1"
	userRepo "broker/app/repository/user"
	authSrv "broker/app/service/auth"
	"broker/pkg/pg"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi"

	chim "github.com/go-chi/chi/middleware"
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

	userRepo := userRepo.NewRepository(pgConn)
	authService := authSrv.NewAuthService([]byte(a.config.JWT.SigningKey), a.config.JWT.ExpireDuration, userRepo)
	authController := auth.NewController(authService)
	authRouter := auth.NewRouter(authController)

	a.web.Router().Route("/api/v1", func(v1 chi.Router) {
		v1.Use(
			chim.Logger,
			chim.RequestID,
		)
		authRouter.InitRoutes(v1)
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
	log.Println("kek")
}
