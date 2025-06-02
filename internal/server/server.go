package server

import (
	"net/http"
	"platform-go-challenge/internal/audience"
	"platform-go-challenge/internal/chart"
	"platform-go-challenge/internal/config"
	"platform-go-challenge/internal/database"
	"platform-go-challenge/internal/favourite"
	"platform-go-challenge/internal/insight"
	"platform-go-challenge/internal/user"
	"platform-go-challenge/internal/utils"
)

func wireDependencies() (*RouterDependencies, error) {
	cfg := config.NewConfig()
	jwtAuth := utils.NewJWTAuth(cfg.JWTSecretKey)
	db := database.NewInMemoryDatabase(cfg.Environment)

	// Users
	userRepository := user.NewInMemoryDBUserRepository(db)
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
	chartRepository := chart.NewInMemoryDBChartRepository(db)
	insightRepository := insight.NewInMemoryDBInsightRepository(db)
	audienceRepository := audience.NewInMemoryDBAudienceRepository(db)
	favouriteRepository := favourite.NewInMemoryDBFavouriteRepository(db)

	favouriteService := favourite.NewFavouriteService(favourite.FavouriteServiceDependencies{
		ChartRepository:     chartRepository,
		InsightRepository:   insightRepository,
		AudienceRepository:  audienceRepository,
		FavouriteRepository: favouriteRepository,
	})

	getFavouritesHandler := favourite.GetFavouritesHandler(
		favourite.GetFavouritesHandlerDependencies{
			FavouriteService: &favouriteService,
		},
	)

	createFavouriteHandler := favourite.CreateFavouriteHandler(
		favourite.CreateFavouriteHandlerDependencies{
			FavouriteService: &favouriteService,
		},
	)

	updateFavouriteHandler := favourite.UpdateFavouriteHandler(
		favourite.UpdateFavouriteHandlerDependencies{
			FavouriteService: &favouriteService,
		},
	)

	// Routing
	routerDependencies := RouterDependencies{
		JWTAuth:                jwtAuth,
		UserLoginHandler:       userLoginHandler,
		GetFavouritesHandler:   getFavouritesHandler,
		CreateFavouriteHandler: createFavouriteHandler,
		UpdateFavouriteHandler: updateFavouriteHandler,
	}

	return &routerDependencies, nil
}

func StartServer() {
	dependencies, _ := wireDependencies()

	router := SetupRouter(*dependencies)

	http.ListenAndServe(":3008", router)
}
