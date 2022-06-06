package user


import "github.com/go-chi/chi"

func RegisterRouter(r chi.Router, controller *Controller) {
	r.Group(func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Post("/signUp", controller.SignUp)
		})
	})
}