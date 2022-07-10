package app

import (
	"broker/core/config"
	delivery "broker/core/delivery/http"
	"broker/core/delivery/http/auth/v1"
	"broker/core/delivery/http/invitation/v1"
	"broker/core/delivery/http/middleware"
	"broker/core/delivery/http/peer/v1"
	"broker/core/delivery/http/workspace/v1"
	invitationRepo "broker/core/repository/invitation"
	peerRepo "broker/core/repository/peer"
	userRepo "broker/core/repository/user"
	workspaceRepo "broker/core/repository/workspace"
	authSrv "broker/core/service/auth"
	connectionSrv "broker/core/service/connection"
	invitationSrv "broker/core/service/invitation"
	invitationPublisher "broker/core/service/invitation/publisher"
	workspaceSrv "broker/core/service/workspace"

	peerSrv "broker/core/service/peer"

	peerConsumerAmqp "broker/core/service/peer/consumer"
	peerPublisherAmqp "broker/core/service/peer/publisher"
	"fmt"
	"time"

	cache "broker/pkg/cache/redis"
	"broker/pkg/logger"
	mailer "broker/pkg/mailer/smtp"
	"broker/pkg/pg"
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi"
	"github.com/go-redis/redis/v8"

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
		time.Sleep(10 * time.Second)
		pgConn, err = pg.NewReadAndWriteConnection(ctx, a.config.Database, a.config.Database, nil)
		if err != nil {
			logger.Logger.Fatalln(err)
		}
	}

	logger.Logger.Info("Database connection established")

	connStr := fmt.Sprintf("amqp://%s:%s@%s:%d", config.AMQP.User, config.AMQP.Pass, config.AMQP.Host, config.AMQP.Port)
	fmt.Println(connStr)

	amqpConn, err := amqp.Dial(connStr)
	if err != nil {
		// hack for await rabbitmq connection
		time.Sleep(10 * time.Second)
		amqpConn, err = amqp.Dial(connStr)
		if err != nil {
			logger.Logger.Fatalln(err)
		}
	}

	logger.Logger.Info("AMQP connection established")

	redis := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port),
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
	})

	authCache := cache.NewRedisCache[int](redis,
		time.Duration(time.Minute*5),
		fmt.Sprintf("%s_%s", getAppPrefix(), "auth"),
		10000).
		WithExpirationTime(time.Duration(time.Minute * 5))

	a.web = delivery.NewAPIServer(":80").WithCors()

	mailer := mailer.NewSmtpMailer()
	userRepo := userRepo.NewRepository(pgConn)
	invitationRepo := invitationRepo.NewRepository(pgConn)
	authService := authSrv.NewAuthService([]byte(a.config.JWT.SigningKey), a.config.JWT.AccessExpireDuration, a.config.JWT.RefreshExpireDuration, userRepo, invitationRepo, authCache, mailer)
	authController := auth.NewController(authService)
	authRouter := auth.NewRouter(authController)

	peerConsumer := peerConsumerAmqp.NewConsumer(amqpConn)
	if err := peerConsumer.Init(); err != nil {
		logger.Logger.Fatalln(err)
	}

	peerPublisher := peerPublisherAmqp.NewPublisher(amqpConn)

	workspaceRepo := workspaceRepo.NewRepository(pgConn)
	peerRepo := peerRepo.NewRepository(pgConn)
	workspaceService := workspaceSrv.NewWorkspaceService(workspaceRepo, userRepo, peerRepo)
	workspaceController := workspace.NewController(workspaceService)

	authGuard := middleware.NewAuthGuard(authService)

	workspaceRouter := workspace.NewRouter(workspaceController, *authGuard)

	peerService := peerSrv.NewPeerService(peerConsumer, peerPublisher)
	peerController := peer.NewController(peerService)
	peerRouter := peer.NewRouter(peerController, authGuard)

	connectionService := connectionSrv.NewConnectionService()

	invitationPublisher := invitationPublisher.NewPublisher(amqpConn)

	invitationService := invitationSrv.NewInvitationService(workspaceRepo, invitationRepo, userRepo, connectionService, mailer, invitationPublisher)

	invitationService.StartScheduler(ctx)

	go func() {
		if err := invitationService.ReadError(); err != nil {
			logger.Logger.Fatalln(err)
		}
	}()

	invitationController := invitation.NewController(invitationService)
	invitationRouter := invitation.NewRouter(invitationController, *authGuard)

	a.web.Router().Route("/api/v1", func(v1 chi.Router) {
		v1.Use(
			chim.Logger,
			chim.RequestID,
		)
		authRouter.InitRoutes(v1)
		workspaceRouter.InitRoutes(v1)
		peerRouter.InitRoutes(v1)
		invitationRouter.InitRoutes(v1)
	})
	return a
}

func (a *App) Run() {
	if err := a.web.Start(); err != nil {
		logger.Logger.Fatal(err)
	}
	appCloser := make(chan os.Signal)
	signal.Notify(appCloser, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-appCloser
		logger.Logger.Info("[os.SIGNAL] close request")
		go a.web.Stop()
		logger.Logger.Info("[os.SIGNAL] done")
	}()
	a.web.WaitForDone()
}

func getAppPrefix() string {
	return "broker_app"
}
