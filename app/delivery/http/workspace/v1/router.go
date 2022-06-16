package workspace

import (
	"broker/app/delivery/http/middleware"

	"github.com/go-chi/chi"
)

type Router struct {
	controller *Controller
	authGuard  *middleware.AuthGuard
}

func NewRouter(controller *Controller, authGuard middleware.AuthGuard) *Router {
	return &Router{
		controller: controller,
		authGuard:  &authGuard,
	}
}

func (rt Router) InitRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Route("/workspaces", func(r chi.Router) {
			r.Use(rt.authGuard.Next())
			r.Get("/", rt.controller.GetManyByUser)
			r.Post("/create", rt.controller.Create)
		})
	})
}
