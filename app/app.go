package app

import (
	"broker/app/config"
	delivery "broker/app/delivery/http"
	"broker/app/delivery/http/auth/v1"
	"broker/app/delivery/http/middleware"
	"broker/app/delivery/http/peer/v1"
	"broker/app/delivery/http/workspace/v1"
	userRepo "broker/app/repository/user"
	workspaceRepo "broker/app/repository/workspace"
	authSrv "broker/app/service/auth"
	workspaceSrv "broker/app/service/workspace"

	peerSrv "broker/app/service/peer"

	peerConsumerAmqp "broker/app/service/peer/consumer/amqp"
	"fmt"
	"time"

	"broker/pkg/pg"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi"

	chim "github.com/go-chi/chi/middleware"
	"github.com/streadway/amqp"
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

	connStr := fmt.Sprintf("amqp://%s:%s@%s:%d", config.AMQP.User, config.AMQP.Pass, config.AMQP.Host, config.AMQP.Port)
	fmt.Println(connStr)

	amqpConn, err := amqp.Dial(connStr)
	if err != nil {
		// hack for await rabbitmq connection
		time.Sleep(10 * time.Second)
		amqpConn, err = amqp.Dial(connStr)
		if err != nil {
			log.Fatalln(err)
		}
	}

	log.Println("AMQP connection established")

	a.web = delivery.NewAPIServer(":80").WithCors()

	userRepo := userRepo.NewRepository(pgConn)
	authService := authSrv.NewAuthService([]byte(a.config.JWT.SigningKey), a.config.JWT.ExpireDuration, userRepo)
	authController := auth.NewController(authService)
	authRouter := auth.NewRouter(authController)

	peerConsumer := peerConsumerAmqp.NewAmqpWorkspaceConsumer(amqpConn)
	if err := peerConsumer.Init(); err != nil {
		log.Fatalln(err)
	}

	workspaceRepo := workspaceRepo.NewRepository(pgConn)
	workspaceService := workspaceSrv.NewWorkspaceService(workspaceRepo, userRepo)
	workspaceController := workspace.NewController(workspaceService)

	authGuard := middleware.NewAuthGuard(authService)

	workspaceRouter := workspace.NewRouter(workspaceController, *authGuard)

	peerService := peerSrv.NewPeerService(peerConsumer)
	peerController := peer.NewController(peerService)
	peerRouter := peer.NewRouter(peerController, authGuard)

	a.web.Router().Route("/api/v1", func(v1 chi.Router) {
		v1.Use(
			chim.Logger,
			chim.RequestID,
		)
		authRouter.InitRoutes(v1)
		workspaceRouter.InitRoutes(v1)
		peerRouter.InitRoutes(v1)
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
