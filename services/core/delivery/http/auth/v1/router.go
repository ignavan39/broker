package auth

import "github.com/go-chi/chi"

type Router struct {
	controller *Controller
}

func NewRouter(controller *Controller) *Router {
	return &Router{
		controller: controller,
	}
}

func (rt Router) InitRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/signUp", rt.controller.SignUp)
			r.Post("/signIn", rt.controller.SignIn)
			r.Post("/sendVerifyCode", rt.controller.SendVerifyCode)
		})
	})
}
