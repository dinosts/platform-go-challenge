package favourite

import (
	"fmt"
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
		userId, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "User Id param is not a UUID")
			return
		}

		pageSize, pageNumber, err := utils.GetPaginationQuery(r, 10, 0)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		fmt.Printf("%s, %d, %d", userId, pageNumber, pageSize)
		assetFavourites, pagination, err := dependencies.FavouriteService.GetPaginatedForUser(userId, pageSize, pageNumber)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		utils.RespondWithPaginatedData(w, http.StatusOK, *assetFavourites, *pagination)
	}
}
