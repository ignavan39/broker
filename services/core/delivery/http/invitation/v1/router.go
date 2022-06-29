package invitation

import (
	"broker/core/delivery/http/middleware"

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
			r.Post("/accept", rt.controller.AcceptInvitation)
			r.Get("/{workspaceID}", rt.controller.GetInvitations)
			r.Delete("/cancel/{invitationID}", rt.controller.CancelInvitation)
		})
	})
}
