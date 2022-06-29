package peer

import (
	"broker/core/delivery/http/middleware"

	"github.com/go-chi/chi"
)

type Router struct {
	controller *Controller
	authGuard  *middleware.AuthGuard
}

func NewRouter(
	controller *Controller,
	authGuard *middleware.AuthGuard,
) *Router {
	return &Router{
		controller: controller,
		authGuard:  authGuard,
	}
}

func (rt Router) InitRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Route("/peers", func(r chi.Router) {
			r.Use(rt.authGuard.Next())
			r.Post("/connect", rt.controller.CreateConnection)
		})
	})
}
