package server

import (
	"net/http"
	"platform-go-challenge/internal/config"
	"platform-go-challenge/internal/database"
	"platform-go-challenge/internal/user"
)

func wireDependencies() (*RouterDependencies, error) {
	cfg := config.LoadConfig()
	jwtAuth := NewJWTAuth(cfg.JWTSecretKey)
	db := database.NewInMemoryDatabase()

	// Users
	userRepository := user.InMemoryDBUserRepository{DB: db}
	userService := user.UserService{Dependencies: user.ServiceDependencies{
		GenerateToken:  NewJWTokenIssuer(jwtAuth),
		UserRepository: &userRepository,
	}}

	userLoginHandler := user.UserLoginHandler(
		user.UserLoginDependencies{
			UserService: userService,
		},
	)

	// Favourites

	routerDependencies := RouterDependencies{
		JWTAuth:          jwtAuth,
		UserLoginHandler: userLoginHandler,
	}

	return &routerDependencies, nil
}

func StartServer() {
	routerDependencies, _ := wireDependencies()

	router := SetupRouter(*routerDependencies)

	http.ListenAndServe(":8008", router)
}
