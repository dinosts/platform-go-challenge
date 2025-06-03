package test

import (
	"net/http/httptest"
	"platform-go-challenge/internal/database"
	"platform-go-challenge/internal/domain/audience"
	"platform-go-challenge/internal/domain/chart"
	"platform-go-challenge/internal/domain/favourite"
	"platform-go-challenge/internal/domain/insight"
	"platform-go-challenge/internal/domain/user"
	"platform-go-challenge/internal/server"
	"platform-go-challenge/internal/utils"
)

func StartServer() (*httptest.Server, string) {
	jwtAuth := utils.NewJWTAuth("test-secret")
	passwordHasher := utils.NewHasher("test-secret")
	db := database.NewIMDatabase()
	tokenIssuer := utils.NewJWTokenIssuer(jwtAuth)

	database.IMpopulateStorageForDevEnv(db, passwordHasher)

	// Users
	userRepository := user.NewInMemoryDBUserRepository(db)
	userService := user.NewUserService(user.ServiceDependencies{
		GenerateToken:  tokenIssuer,
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
	routerDependencies := server.RouterDependencies{
		JWTAuth:                jwtAuth,
		UserLoginHandler:       userLoginHandler,
		GetFavouritesHandler:   getFavouritesHandler,
		CreateFavouriteHandler: createFavouriteHandler,
		UpdateFavouriteHandler: updateFavouriteHandler,
		DeleteFavouriteHandler: deleteFavouriteHandler,
	}

	router := server.SetupRouter(routerDependencies)

	server := httptest.NewServer(router)

	token, _, _ := tokenIssuer(
		map[string]any{
			"sub": "a3973a1c-a77b-4a04-a296-ddec19034419",
		},
	)
	return server, token
}
