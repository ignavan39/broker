package delivery

import (
	"broker/core/delivery/http/middleware"

	"github.com/go-chi/chi"
)

type Router struct {
	controller *Controller
	authGuard  *middleware.AuthGuard
}

func NewRouter(authGuard *middleware.AuthGuard, controller *Controller) *Router {
	return &Router{
		controller: controller,
		authGuard:  authGuard,
	}
}

func (rt Router) InitRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Route("/connect", func(r chi.Router) {
			r.Use(rt.authGuard.Next())
			r.Post("/ping", rt.controller.Ping)
		})
	})
}
