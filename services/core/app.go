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
	invitationSrv "broker/core/service/invitation"
	workspaceSrv "broker/core/service/workspace"

	peerSrv "broker/core/service/peer"

	peerConsumerAmqp "broker/core/service/peer/consumer/amqp"
	peerPublisherAmqp "broker/core/service/peer/publisher/amqp"
	"fmt"
	"time"

	cache "broker/pkg/cache/redis"
	mailer "broker/pkg/mailer/mock"
	"broker/pkg/pg"
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi"
	"github.com/go-redis/redis/v8"

	chim "github.com/go-chi/chi/middleware"
	"github.com/streadway/amqp"

	blogger "github.com/sirupsen/logrus"
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
		blogger.Fatalln(err)
	}

	blogger.Info("Database connection established")

	connStr := fmt.Sprintf("amqp://%s:%s@%s:%d", config.AMQP.User, config.AMQP.Pass, config.AMQP.Host, config.AMQP.Port)
	fmt.Println(connStr)

	amqpConn, err := amqp.Dial(connStr)
	if err != nil {
		// hack for await rabbitmq connection
		time.Sleep(10 * time.Second)
		amqpConn, err = amqp.Dial(connStr)
		if err != nil {
			blogger.Fatalln(err)
		}
	}

	blogger.Info("AMQP connection established")

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

	mailGun := mailer.NewMockMailer()
	userRepo := userRepo.NewRepository(pgConn)
	invitationRepo := invitationRepo.NewRepository(pgConn)
	authService := authSrv.NewAuthService([]byte(a.config.JWT.SigningKey), a.config.JWT.AccessExpireDuration, a.config.JWT.RefreshExpireDuration, userRepo, invitationRepo, authCache, mailGun)
	authController := auth.NewController(authService)
	authRouter := auth.NewRouter(authController)

	peerConsumer := peerConsumerAmqp.NewConsumer(amqpConn)
	if err := peerConsumer.Init(); err != nil {
		blogger.Fatalln(err)
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

	invitationService := invitationSrv.NewInvitationService(workspaceRepo, invitationRepo)
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
		blogger.Fatal(err)
	}
	appCloser := make(chan os.Signal)
	signal.Notify(appCloser, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-appCloser
		blogger.Info("[os.SIGNAL] close request")
		go a.web.Stop()
		blogger.Info("[os.SIGNAL] done")
	}()
	a.web.WaitForDone()
}

func getAppPrefix() string {
	return "broker_app"
}
