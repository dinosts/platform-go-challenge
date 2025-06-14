package favourite

import (
	"errors"
	"net/http"
	"platform-go-challenge/internal/utils"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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

type UpdateFavouriteHandlerDependencies struct {
	FavouriteService FavouriteService
}

func UpdateFavouriteHandler(dependencies UpdateFavouriteHandlerDependencies) http.HandlerFunc {
	validation := utils.BodyValidator[UpdateFavouriteRequestBody]
	handler := func(w http.ResponseWriter, r *http.Request) {
		favouriteId, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Favourite Id param is not a UUID")
			return
		}

		userId, err := utils.GetUserIdFromAuthToken(r)
		if err != nil {
			// Should not happen since we have auth middlewares before this route
			utils.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		body, ok := utils.GetParsedBody[UpdateFavouriteRequestBody](r)
		if !ok {
			// Should not happen since we validate body before getting in to handler
			utils.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		favourite, err := dependencies.FavouriteService.Update(userId, favouriteId, body.Description)
		if err != nil {
			if errors.Is(err, ErrFavouriteNotFound) {
				utils.RespondWithError(w, http.StatusNotFound, "Could not find Favourite with this Id")
				return
			}
			if errors.Is(err, ErrFavouriteNotUnderGivenUser) {
				utils.RespondWithError(w, http.StatusUnauthorized, "Favourite is not under given user")
				return
			}

			utils.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		utils.RespondWithData(w, http.StatusOK, favourite)
	}

	return validation(handler)
}

type DeleteFavouriteHandlerDependencies struct {
	FavouriteService FavouriteService
}

func DeleteFavouriteHandler(dependencies DeleteFavouriteHandlerDependencies) http.HandlerFunc {
	handler := func(w http.ResponseWriter, r *http.Request) {
		favouriteId, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Favourite Id param is not a UUID")
			return
		}

		userId, err := utils.GetUserIdFromAuthToken(r)
		if err != nil {
			// Should not happen since we have auth middlewares before this route
			utils.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		err = dependencies.FavouriteService.Delete(userId, favouriteId)
		if err != nil {
			if errors.Is(err, ErrFavouriteNotFound) {
				utils.RespondWithError(w, http.StatusNotFound, "Could not find Asset with this Id")
				return
			}
			if errors.Is(err, ErrFavouriteNotUnderGivenUser) {
				utils.RespondWithError(w, http.StatusUnauthorized, "Favourite is not under given user")
				return
			}

			utils.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		utils.RespondWithMessage(w, http.StatusOK, "Favourite deleted")
	}

	return handler
}
