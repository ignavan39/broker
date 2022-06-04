package api

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

type Server struct {
	srv    *http.Server
	done   chan error
	router chi.Router
}

func NewAPIServer(listenOn string) *Server {
	router := chi.NewRouter()

	return &Server{
		srv:    &http.Server{Addr: listenOn, Handler: router},
		done:   make(chan error),
		router: router,
	}
}
func (a *Server) Router() chi.Router {
	return a.router
}

func (a *Server) Start() error {
	go func() {
		defer close(a.done)
		if err := a.srv.ListenAndServe(); err != nil {
			a.done <- err
		}
	}()
	return nil
}
func (a *Server) Stop() {
	a.srv.Shutdown(context.Background())
}
func (a *Server) WaitForDone() error {
	log.Println("Server has been started")
	return <-a.done
}
