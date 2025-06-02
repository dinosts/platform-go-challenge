package favourite

import (
	"errors"
	"net/http"
	"platform-go-challenge/internal/utils"
)

type GetFavouritesHandlerDependencies struct {
	FavouriteService FavouriteService
}

func GetFavouritesHandler(dependencies GetFavouritesHandlerDependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, err := utils.GetUserIdFromAuthToken(r)
		if err != nil {
			// Should not happen since we have auth middlewares before this route
			utils.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		pageSize, pageNumber, err := utils.GetPaginationQuery(r, 10, 0)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		assetFavourites, pagination, err := dependencies.FavouriteService.GetPaginatedForUser(userId, pageSize, pageNumber)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		utils.RespondWithPaginatedData(w, http.StatusOK, *assetFavourites, *pagination)
	}
}

type CreateFavouriteHandlerDependencies struct {
	FavouriteService FavouriteService
}

func CreateFavouriteHandler(dependencies CreateFavouriteHandlerDependencies) http.HandlerFunc {
	validation := utils.BodyValidator[CreateFavouriteRequestBody]
	handler := func(w http.ResponseWriter, r *http.Request) {
		userId, err := utils.GetUserIdFromAuthToken(r)
		if err != nil {
			// Should not happen since we have auth middlewares before this route
			utils.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		body, ok := utils.GetParsedBody[CreateFavouriteRequestBody](r)
		if !ok {
			// Should not happen since we validate body before getting in to handler
			utils.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		favourite, err := dependencies.FavouriteService.CreateForUser(userId, body.AssetId, body.Description)
		if err != nil {
			if errors.Is(err, ErrAssetNotFound) {
				utils.RespondWithError(w, http.StatusNotFound, "Could not find Asset with this Id")
				return
			}

			utils.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		utils.RespondWithData(w, http.StatusCreated, favourite)
	}

	return validation(handler)
}
