package delivery

import "github.com/go-chi/chi"

type Router struct {
	controller Controller
}

func NewRouter(controller Controller) *Router {
	return &Router{
		controller: controller,
	}
}

func (rt Router) InitRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Route("/mail", func(r chi.Router) {
			r.Post("/send", rt.controller.Send)
		})
	})
}
