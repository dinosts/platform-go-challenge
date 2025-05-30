package server

import (
	"net/http"
	"platform-go-challenge/internal/config"
)

func wireDependencies() (*RouterDependencies, error) {
	cfg := config.LoadConfig()
	routerDependencies := RouterDependencies{
		JWTAuth: NewJWTAuth(cfg.JWTSecretKey),
	}

	return &routerDependencies, nil
}

func StartServer() {
	routerDependencies, _ := wireDependencies()

	router := SetupRouter(*routerDependencies)

	http.ListenAndServe(":8008", router)
}
