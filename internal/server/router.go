package server

import (
	"net/http"
	"platform-go-challenge/internal/utils"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/go-chi/jwtauth/v5"
)

type RouterDependencies struct {
	JWTAuth                *jwtauth.JWTAuth
	UserLoginHandler       http.HandlerFunc
	GetFavouritesHandler   http.HandlerFunc
	CreateFavouriteHandler http.HandlerFunc
	UpdateFavouriteHandler http.HandlerFunc
	DeleteFavouriteHandler http.HandlerFunc
}

func SetupRouter(dependencies RouterDependencies) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// ratelimit: 100req per 1min
	r.Use(httprate.LimitByIP(100, 1*time.Minute))

	r.Get("/", GetHealth)
	r.Route("/v1", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Route("/user", func(r chi.Router) {
				// Public
				r.Group(func(r chi.Router) {
					r.Post("/login", dependencies.UserLoginHandler)
				})

				// Private
				r.Group(func(r chi.Router) {
					r.Use(utils.VerifierMiddleware(dependencies.JWTAuth))
					r.Use(utils.AuthenticatorMiddleware())

					r.Get("/favourites", dependencies.GetFavouritesHandler)
					r.Post("/favourites", dependencies.CreateFavouriteHandler)
					r.Patch("/favourites", dependencies.UpdateFavouriteHandler)
					r.Delete("/favourites", dependencies.DeleteFavouriteHandler)
				})
			})
		})
	})

	return r
}
