package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/go-chi/chi/v5/middleware"
)

func SetupRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", getHealth)
	r.Route("/v1/", func(r chi.Router) {
		// Public
		r.Group(func(r chi.Router) {
		})

		// Private
		r.Group(func(r chi.Router) {
			// r.Use(AuthMiddleware)
			// Admin
			r.Group(func(r chi.Router) {})
		})
	})

	return r
}

func StartServer() {
	router := SetupRouter()
	http.ListenAndServe(":3000", router)
}
