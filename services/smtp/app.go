package smtp

import (
	"broker/smtp/config"
	delivery "broker/smtp/sender/delivery/http"
	"broker/smtp/sender/services"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi"
	chim "github.com/go-chi/chi/middleware"
	blogger "github.com/sirupsen/logrus"
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
		blogger.Fatal(err)
	}
	appCloser := make(chan os.Signal)
	signal.Notify(appCloser, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-appCloser
		blogger.Info("[os.SIGNAL] close request")
		go web.Stop()
		blogger.Info("[os.SIGNAL] done")
	}()
	web.WaitForDone()
}

