package invitation

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
		r.Route("/invitations", func(r chi.Router) {
			r.Use(rt.authGuard.Next())
			r.Post("/create/{workspaceID}", rt.controller.SendInvitation)
			r.Get("/{workspaceID}", rt.controller.GetInvitations)
		})
	})
}
