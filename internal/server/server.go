package server

import (
	"net/http"
	"platform-go-challenge/internal/config"
	"platform-go-challenge/internal/database"
	"platform-go-challenge/internal/domain/audience"
	"platform-go-challenge/internal/domain/chart"
	"platform-go-challenge/internal/domain/favourite"
	"platform-go-challenge/internal/domain/insight"
	"platform-go-challenge/internal/domain/user"
	"platform-go-challenge/internal/utils"
)

func wireDependencies(cfg config.Config) (*RouterDependencies, error) {
	jwtAuth := utils.NewJWTAuth(cfg.JWTSecretKey)
	passwordHasher := utils.NewHasher(cfg.HashingSalt)
	db := database.NewIMDatabase()

	if cfg.Environment == "dev" {
		database.IMpopulateStorageForDevEnv(db, passwordHasher)
	}

	// Users
	userRepository := user.NewInMemoryDBUserRepository(db)
	userService := user.NewUserService(user.ServiceDependencies{
		GenerateToken:  utils.NewJWTokenIssuer(jwtAuth),
		UserRepository: &userRepository,
		PasswordHasher: passwordHasher,
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

	deleteFavouriteHandler := favourite.DeleteFavouriteHandler(
		favourite.DeleteFavouriteHandlerDependencies{
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
		DeleteFavouriteHandler: deleteFavouriteHandler,
	}

	return &routerDependencies, nil
}

func StartServer() {
	cfg := config.NewConfig()

	dependencies, _ := wireDependencies(*cfg)

	router := SetupRouter(*dependencies)

	http.ListenAndServe(":3008", router)
}
