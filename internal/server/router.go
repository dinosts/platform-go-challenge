package server

import (
	"net/http"
	"platform-go-challenge/internal/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
)

type RouterDependencies struct {
	JWTAuth          *jwtauth.JWTAuth
	UserLoginHandler func(w http.ResponseWriter, r *http.Request)
}

func SetupRouter(dependencies RouterDependencies) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", GetHealth)

	r.Route("/v1", func(r chi.Router) {
		// Public
		r.Group(func(r chi.Router) {
			r.Route("/user", func(r chi.Router) {
				r.Post("/login", dependencies.UserLoginHandler)
			})
		})

		// Private
		r.Group(func(r chi.Router) {
			r.Use(utils.VerifierMiddleware(dependencies.JWTAuth))
			r.Use(utils.AuthenticatorMiddleware())

			// Admin
			r.Group(func(r chi.Router) {})
		})
	})

	return r
}
