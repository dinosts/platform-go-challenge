package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
)

type RouterDependencies struct {
	JWTAuth *jwtauth.JWTAuth
}

func SetupRouter(dependencies RouterDependencies) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Get("/", getHealth)
	r.Route("/v1", func(r chi.Router) {
		// Public
		r.Group(func(r chi.Router) {
		})

		// Private
		r.Group(func(r chi.Router) {
			r.Use(VerifierMiddleware(dependencies.JWTAuth))
			r.Use(AuthenticatorMiddleware())

			// Admin
			r.Group(func(r chi.Router) {})
		})
	})

	return r
}
