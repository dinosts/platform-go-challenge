package favourite_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"platform-go-challenge/internal/favourite"
	"platform-go-challenge/internal/utils"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// Stub implementation of favouriteService interface without using any mock package
type StubFavouriteService struct {
	GetPaginatedForUserFunc func(userId uuid.UUID, pageSize, pageNumber int) (*favourite.AssetFavourites, *utils.Pagination, error)
}

func (s *StubFavouriteService) GetPaginatedForUser(userId uuid.UUID, pageSize, pageNumber int) (*favourite.AssetFavourites, *utils.Pagination, error) {
	return s.GetPaginatedForUserFunc(userId, pageSize, pageNumber)
}

func TestGetFavouritesHandler(t *testing.T) {
	t.Run("Should return 200 and data when valid request", func(t *testing.T) {
		// Arrange
		validUUID := uuid.New()
		stubService := &StubFavouriteService{
			GetPaginatedForUserFunc: func(userId uuid.UUID, pageSize, pageNumber int) (*favourite.AssetFavourites, *utils.Pagination, error) {
				assert.Equal(t, validUUID, userId)
				assert.Equal(t, 10, pageSize)
				assert.Equal(t, 0, pageNumber)
				return &favourite.AssetFavourites{}, &utils.Pagination{}, nil
			},
		}
		handler := favourite.GetFavouritesHandler(favourite.GetFavouritesHandlerDependencies{
			FavouriteService: stubService,
		})

		req := httptest.NewRequest(http.MethodGet, "/favourites/"+validUUID.String(), nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", validUUID.String())
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		w := httptest.NewRecorder()

		// Act
		handler(w, req)

		// Assert
		resp := w.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Should return 400 when invalid UUID is provided", func(t *testing.T) {
		// Arrange
		invalidUUID := "not-a-uuid"
		stubService := &StubFavouriteService{}

		handler := favourite.GetFavouritesHandler(favourite.GetFavouritesHandlerDependencies{
			FavouriteService: stubService,
		})

		req := httptest.NewRequest(http.MethodGet, "/favourites/"+invalidUUID, nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", invalidUUID)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		w := httptest.NewRecorder()

		// Act
		handler(w, req)

		// Assert
		resp := w.Result()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Should return 400 when invalid pagination parameters are provided", func(t *testing.T) {
		// Arrange
		validUUID := uuid.New().String()
		stubService := &StubFavouriteService{}

		handler := favourite.GetFavouritesHandler(favourite.GetFavouritesHandlerDependencies{
			FavouriteService: stubService,
		})

		req := httptest.NewRequest(http.MethodGet, "/favourites/"+validUUID+"?pageSize=abc", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", validUUID)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		w := httptest.NewRecorder()

		// Act
		handler(w, req)

		// Assert
		resp := w.Result()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Should return 500 when service returns error", func(t *testing.T) {
		// Arrange
		validUUID := uuid.New()
		stubService := &StubFavouriteService{
			GetPaginatedForUserFunc: func(userId uuid.UUID, pageSize, pageNumber int) (*favourite.AssetFavourites, *utils.Pagination, error) {
				return nil, nil, errors.New("internal error")
			},
		}

		handler := favourite.GetFavouritesHandler(favourite.GetFavouritesHandlerDependencies{
			FavouriteService: stubService,
		})

		req := httptest.NewRequest(http.MethodGet, "/favourites/"+validUUID.String(), nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", validUUID.String())
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		w := httptest.NewRecorder()

		// Act
		handler(w, req)

		// Assert
		resp := w.Result()
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}
