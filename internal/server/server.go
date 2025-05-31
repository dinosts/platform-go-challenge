package server

import (
	"net/http"
	"platform-go-challenge/internal/config"
	"platform-go-challenge/internal/database"
	"platform-go-challenge/internal/user"
	"platform-go-challenge/internal/utils"
)

func wireDependencies() (*RouterDependencies, error) {
	cfg := config.NewConfig()
	jwtAuth := utils.NewJWTAuth(cfg.JWTSecretKey)
	db := database.NewInMemoryDatabase(cfg.Environment)

	// Users
	userRepository := user.InMemoryDBUserRepository{DB: db}
	userService := user.NewUserService(user.ServiceDependencies{
		GenerateToken:  utils.NewJWTokenIssuer(jwtAuth),
		UserRepository: &userRepository,
	})

	userLoginHandler := user.UserLoginHandler(
		user.UserLoginDependencies{
			UserService: &userService,
		},
	)

	// Favourites

	// Routing
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
