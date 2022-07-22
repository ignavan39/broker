package smtp

import (
	"broker/pkg/logger"
	"broker/smtp/config"
	delivery "broker/smtp/sender/delivery/http"
	"broker/smtp/sender/services"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi"
	chim "github.com/go-chi/chi/middleware"
)

func Start() {
	config.Init()

	web := NewAPIServer(":80").WithCors()

	router := delivery.NewRouter(*delivery.NewController(*services.NewSender()))

	web.Router().Route("/api/v1", func(v1 chi.Router) {
		v1.Use(
			chim.Logger,
			chim.RequestID,
		)
		router.InitRoutes(v1)
	})

	if err := web.Start(); err != nil {
		logger.Logger.Fatal(err)
	}
	appCloser := make(chan os.Signal)
	signal.Notify(appCloser, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-appCloser
		logger.Logger.Info("[os.SIGNAL] close request")
		go web.Stop()
		logger.Logger.Info("[os.SIGNAL] done")
	}()
	web.WaitForDone()
}
