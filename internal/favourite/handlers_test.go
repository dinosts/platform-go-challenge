package favourite_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"platform-go-challenge/internal/favourite"
	"platform-go-challenge/internal/utils"
	"strings"
	"testing"

	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type StubFavouriteService struct {
	GetPaginatedForUserFunc func(userId uuid.UUID, pageSize, pageNumber int) (*favourite.AssetFavourites, *utils.Pagination, error)
	CreateForUserFunc       func(userId, assetId uuid.UUID, description string) (*favourite.Favourite, error)
	UpdateFunc              func(userId, favouriteId uuid.UUID, description string) (*favourite.Favourite, error)
	DeleteFunc              func(userId, favouriteId uuid.UUID) error
}

func (s *StubFavouriteService) GetPaginatedForUser(userId uuid.UUID, pageSize, pageNumber int) (*favourite.AssetFavourites, *utils.Pagination, error) {
	if s.GetPaginatedForUserFunc != nil {
		return s.GetPaginatedForUserFunc(userId, pageSize, pageNumber)
	}
	return nil, nil, errors.New("not implemented")
}

func (s *StubFavouriteService) CreateForUser(userId, assetId uuid.UUID, description string) (*favourite.Favourite, error) {
	if s.CreateForUserFunc != nil {
		return s.CreateForUserFunc(userId, assetId, description)
	}
	return nil, errors.New("not implemented")
}

func (s *StubFavouriteService) Update(userId, favouriteId uuid.UUID, description string) (*favourite.Favourite, error) {
	if s.UpdateFunc != nil {
		return s.UpdateFunc(userId, favouriteId, description)
	}
	return nil, errors.New("not implemented")
}

func (s *StubFavouriteService) Delete(userId, favouriteId uuid.UUID) error {
	if s.DeleteFunc != nil {
		return s.DeleteFunc(userId, favouriteId)
	}
	return errors.New("not implemented")
}

func injectJWT(ctx context.Context, userID string) context.Context {
	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)
	token, _, _ := tokenAuth.Encode(map[string]interface{}{"sub": userID})
	return jwtauth.NewContext(ctx, token, nil)
}

func injectParsedBody[T any](ctx context.Context, body T) context.Context {
	return context.WithValue(ctx, "parsedBody", body)
}

func TestGetFavouritesHandler(t *testing.T) {
	t.Run("Should return 200 and data when valid JWT", func(t *testing.T) {
		// Arrange
		validUUID := uuid.New()
		stubService := &StubFavouriteService{
			GetPaginatedForUserFunc: func(userId uuid.UUID, pageSize, pageNumber int) (*favourite.AssetFavourites, *utils.Pagination, error) {
				assert.Equal(t, validUUID, userId)
				return &favourite.AssetFavourites{}, &utils.Pagination{}, nil
			},
		}
		handler := favourite.GetFavouritesHandler(favourite.GetFavouritesHandlerDependencies{
			FavouriteService: stubService,
		})
		req := httptest.NewRequest(http.MethodGet, "/favourites", nil)
		req = req.WithContext(injectJWT(req.Context(), validUUID.String()))
		w := httptest.NewRecorder()

		// Act
		handler(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	})

	t.Run("Should return 500 when JWT sub is invalid UUID", func(t *testing.T) {
		// Arrange
		handler := favourite.GetFavouritesHandler(favourite.GetFavouritesHandlerDependencies{
			FavouriteService: &StubFavouriteService{},
		})
		req := httptest.NewRequest(http.MethodGet, "/favourites", nil)
		req = req.WithContext(injectJWT(req.Context(), "not-a-uuid"))
		w := httptest.NewRecorder()

		// Act
		handler(w, req)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	})

	t.Run("Should return 400 when pagination query is invalid", func(t *testing.T) {
		// Arrange
		validUUID := uuid.New()
		handler := favourite.GetFavouritesHandler(favourite.GetFavouritesHandlerDependencies{
			FavouriteService: &StubFavouriteService{},
		})
		req := httptest.NewRequest(http.MethodGet, "/favourites?pageSize=invalid", nil)
		req = req.WithContext(injectJWT(req.Context(), validUUID.String()))
		w := httptest.NewRecorder()

		// Act
		handler(w, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("Should return 500 when service returns error", func(t *testing.T) {
		// Arrange
		validUUID := uuid.New()
		handler := favourite.GetFavouritesHandler(favourite.GetFavouritesHandlerDependencies{
			FavouriteService: &StubFavouriteService{
				GetPaginatedForUserFunc: func(userId uuid.UUID, pageSize, pageNumber int) (*favourite.AssetFavourites, *utils.Pagination, error) {
					return nil, nil, errors.New("fail")
				},
			},
		})
		req := httptest.NewRequest(http.MethodGet, "/favourites", nil)
		req = req.WithContext(injectJWT(req.Context(), validUUID.String()))
		w := httptest.NewRecorder()

		// Act
		handler(w, req)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	})

	t.Run("Should return 500 when JWT is missing", func(t *testing.T) {
		// Arrange
		handler := favourite.GetFavouritesHandler(favourite.GetFavouritesHandlerDependencies{
			FavouriteService: &StubFavouriteService{},
		})
		req := httptest.NewRequest(http.MethodGet, "/favourites", nil)
		w := httptest.NewRecorder()

		// Act
		handler(w, req)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	})
}

func TestCreateFavouriteHandler(t *testing.T) {
	t.Run("Should return 201 when favourite is created successfully", func(t *testing.T) {
		// Arrange
		userId := uuid.New()
		assetId := uuid.New()

		requestBody := map[string]interface{}{
			"assetId":     assetId.String(),
			"description": "test",
		}
		expected := &favourite.Favourite{
			Id:          uuid.New(),
			UserId:      userId,
			AssetId:     assetId,
			AssetType:   "chart",
			Description: "test",
		}
		stubService := &StubFavouriteService{
			CreateForUserFunc: func(uId, aId uuid.UUID, desc string) (*favourite.Favourite, error) {
				assert.Equal(t, userId, uId)
				assert.Equal(t, assetId, aId)
				assert.Equal(t, "test", desc)
				return expected, nil
			},
		}
		handler := favourite.CreateFavouriteHandler(favourite.CreateFavouriteHandlerDependencies{
			FavouriteService: stubService,
		})

		bodyBytes, _ := json.Marshal(requestBody)
		req := httptest.NewRequest(http.MethodPost, "/favourites", bytes.NewReader(bodyBytes))
		req = req.WithContext(injectJWT(req.Context(), userId.String()))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler(w, req)

		// Assert
		assert.Equal(t, http.StatusCreated, w.Result().StatusCode)
	})

	t.Run("Should return 404 when asset is not found", func(t *testing.T) {
		// Arrange
		userId := uuid.New()
		assetId := uuid.New()

		requestBody := map[string]interface{}{
			"assetId":     assetId.String(),
			"description": "test",
		}
		stubService := &StubFavouriteService{
			CreateForUserFunc: func(_, _ uuid.UUID, _ string) (*favourite.Favourite, error) {
				return nil, favourite.ErrAssetNotFound
			},
		}
		handler := favourite.CreateFavouriteHandler(favourite.CreateFavouriteHandlerDependencies{
			FavouriteService: stubService,
		})

		bodyBytes, _ := json.Marshal(requestBody)
		req := httptest.NewRequest(http.MethodPost, "/favourites", bytes.NewReader(bodyBytes))
		req = req.WithContext(injectJWT(req.Context(), userId.String()))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler(w, req)

		// Assert
		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
	})

	t.Run("Should return 400 if body is invalid JSON", func(t *testing.T) {
		// Arrange
		userId := uuid.New()

		handler := favourite.CreateFavouriteHandler(favourite.CreateFavouriteHandlerDependencies{
			FavouriteService: &StubFavouriteService{},
		})

		req := httptest.NewRequest(http.MethodPost, "/favourites", strings.NewReader("{invalid json}"))
		req = req.WithContext(injectJWT(req.Context(), userId.String()))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler(w, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("Should return 400 if required fields are missing", func(t *testing.T) {
		// Arrange
		userId := uuid.New()
		requestBody := map[string]interface{}{}

		bodyBytes, _ := json.Marshal(requestBody)

		handler := favourite.CreateFavouriteHandler(favourite.CreateFavouriteHandlerDependencies{
			FavouriteService: &StubFavouriteService{},
		})

		req := httptest.NewRequest(http.MethodPost, "/favourites", bytes.NewReader(bodyBytes))
		req = req.WithContext(injectJWT(req.Context(), userId.String()))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler(w, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("Should return 500 if JWT is missing", func(t *testing.T) {
		// Arrange
		assetId := uuid.New()
		requestBody := map[string]interface{}{
			"assetId":     assetId.String(),
			"description": "testaki",
		}

		bodyBytes, _ := json.Marshal(requestBody)

		handler := favourite.CreateFavouriteHandler(favourite.CreateFavouriteHandlerDependencies{
			FavouriteService: &StubFavouriteService{},
		})

		req := httptest.NewRequest(http.MethodPost, "/favourites", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler(w, req)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	})
}

func TestUpdateFavouriteHandler(t *testing.T) {
	t.Run("Should return 200 when update is successful", func(t *testing.T) {
		// Arrange
		userId := uuid.New()
		favouriteId := uuid.New()

		requestBody := map[string]interface{}{
			"id":          favouriteId.String(),
			"description": "updated description",
		}

		expected := &favourite.Favourite{
			Id:          favouriteId,
			UserId:      userId,
			AssetId:     uuid.New(),
			AssetType:   "chart",
			Description: "test",
		}
		stubService := &StubFavouriteService{
			UpdateFunc: func(uId, fId uuid.UUID, desc string) (*favourite.Favourite, error) {
				assert.Equal(t, userId, uId)
				assert.Equal(t, favouriteId, fId)
				assert.Equal(t, "updated description", desc)
				return expected, nil
			},
		}
		handler := favourite.UpdateFavouriteHandler(favourite.UpdateFavouriteHandlerDependencies{
			FavouriteService: stubService,
		})

		bodyBytes, _ := json.Marshal(requestBody)
		req := httptest.NewRequest(http.MethodPatch, "/favourites", bytes.NewReader(bodyBytes))
		req = req.WithContext(injectJWT(req.Context(), userId.String()))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	})

	t.Run("Should return 404 when favourite is not found", func(t *testing.T) {
		// Arrange
		userId := uuid.New()
		favouriteId := uuid.New()

		requestBody := map[string]interface{}{
			"id":          favouriteId.String(),
			"description": "does not matter",
		}

		stubService := &StubFavouriteService{
			UpdateFunc: func(_, _ uuid.UUID, _ string) (*favourite.Favourite, error) {
				return nil, favourite.ErrFavouriteNotFound
			},
		}
		handler := favourite.UpdateFavouriteHandler(favourite.UpdateFavouriteHandlerDependencies{
			FavouriteService: stubService,
		})

		bodyBytes, _ := json.Marshal(requestBody)
		req := httptest.NewRequest(http.MethodPatch, "/favourites", bytes.NewReader(bodyBytes))
		req = req.WithContext(injectJWT(req.Context(), userId.String()))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler(w, req)

		// Assert
		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
	})

	t.Run("Should return 401 when favourite is not under user", func(t *testing.T) {
		// Arrange
		userId := uuid.New()
		favouriteId := uuid.New()
		requestBody := map[string]interface{}{
			"id":          favouriteId.String(),
			"description": "test",
		}
		stubService := &StubFavouriteService{
			UpdateFunc: func(_, _ uuid.UUID, _ string) (*favourite.Favourite, error) {
				return nil, favourite.ErrFavouriteNotUnderGivenUser
			},
		}
		handler := favourite.UpdateFavouriteHandler(favourite.UpdateFavouriteHandlerDependencies{
			FavouriteService: stubService,
		})

		bodyBytes, _ := json.Marshal(requestBody)
		req := httptest.NewRequest(http.MethodPatch, "/favourites", bytes.NewReader(bodyBytes))
		req = req.WithContext(injectJWT(req.Context(), userId.String()))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler(w, req)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
	})

	t.Run("Should return 500 when JWT is missing", func(t *testing.T) {
		// Arrange
		favouriteId := uuid.New()
		requestBody := map[string]interface{}{
			"id":          favouriteId.String(),
			"description": "test",
		}
		handler := favourite.UpdateFavouriteHandler(favourite.UpdateFavouriteHandlerDependencies{
			FavouriteService: &StubFavouriteService{},
		})
		bodyBytes, _ := json.Marshal(requestBody)
		req := httptest.NewRequest(http.MethodPatch, "/favourites", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler(w, req)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	})

	t.Run("Should return 400 on unexpected service error", func(t *testing.T) {
		// Arrange
		userId := uuid.New()
		stubService := &StubFavouriteService{
			UpdateFunc: func(_, _ uuid.UUID, _ string) (*favourite.Favourite, error) {
				return nil, errors.New("random error")
			},
		}
		handler := favourite.UpdateFavouriteHandler(favourite.UpdateFavouriteHandlerDependencies{
			FavouriteService: stubService,
		})
		req := httptest.NewRequest(http.MethodPut, "/favourites", strings.NewReader(`{}`))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(injectJWT(req.Context(), userId.String()))
		w := httptest.NewRecorder()

		// Act
		handler(w, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})
}

func TestDeleteFavouriteHandler(t *testing.T) {
	t.Run("Should return 200 when delete is successful", func(t *testing.T) {
		// Arrange
		userId := uuid.New()
		favouriteId := uuid.New()

		requestBody := map[string]any{"id": favouriteId.String()}

		stubService := &StubFavouriteService{
			DeleteFunc: func(uId, fId uuid.UUID) error {
				assert.Equal(t, userId, uId)
				assert.Equal(t, favouriteId, fId)

				return nil
			},
		}
		handler := favourite.DeleteFavouriteHandler(
			favourite.DeleteFavouriteHandlerDependencies{
				FavouriteService: stubService,
			},
		)

		bodyBytes, _ := json.Marshal(requestBody)
		req := httptest.NewRequest(http.MethodDelete, "/favourites", bytes.NewReader(bodyBytes))
		req = req.WithContext(injectJWT(req.Context(), userId.String()))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	})

	t.Run("Should return 404 when favourite is not found", func(t *testing.T) {
		// Arrange
		userId := uuid.New()

		requestBody := map[string]interface{}{
			"id": uuid.New().String(),
		}

		stubService := &StubFavouriteService{
			DeleteFunc: func(_, _ uuid.UUID) error {
				return favourite.ErrFavouriteNotFound
			},
		}
		handler := favourite.DeleteFavouriteHandler(favourite.DeleteFavouriteHandlerDependencies{
			FavouriteService: stubService,
		})

		bodyBytes, _ := json.Marshal(requestBody)
		req := httptest.NewRequest(http.MethodDelete, "/favourites", bytes.NewReader(bodyBytes))
		req = req.WithContext(injectJWT(req.Context(), userId.String()))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler(w, req)

		// Assert
		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
	})

	t.Run("Should return 401 when favourite is not under user", func(t *testing.T) {
		// Arrange
		userId := uuid.New()
		requestBody := map[string]interface{}{
			"id": uuid.New().String(),
		}

		stubService := &StubFavouriteService{
			DeleteFunc: func(_, _ uuid.UUID) error {
				return favourite.ErrFavouriteNotUnderGivenUser
			},
		}

		handler := favourite.DeleteFavouriteHandler(favourite.DeleteFavouriteHandlerDependencies{FavouriteService: stubService})

		bodyBytes, _ := json.Marshal(requestBody)
		req := httptest.NewRequest(http.MethodPatch, "/favourites", bytes.NewReader(bodyBytes))
		req = req.WithContext(injectJWT(req.Context(), userId.String()))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler(w, req)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
	})

	t.Run("Should return 500 when JWT is missing", func(t *testing.T) {
		// Arrange
		requestBody := map[string]interface{}{
			"id": uuid.New().String(),
		}
		handler := favourite.DeleteFavouriteHandler(favourite.DeleteFavouriteHandlerDependencies{
			FavouriteService: &StubFavouriteService{},
		})
		bodyBytes, _ := json.Marshal(requestBody)
		req := httptest.NewRequest(http.MethodPatch, "/favourites", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Act
		handler(w, req)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	})

	t.Run("Should return 400 on unexpected service error", func(t *testing.T) {
		// Arrange
		userId := uuid.New()
		stubService := &StubFavouriteService{
			DeleteFunc: func(_, _ uuid.UUID) error {
				return errors.New("random error")
			},
		}

		handler := favourite.DeleteFavouriteHandler(favourite.DeleteFavouriteHandlerDependencies{FavouriteService: stubService})
		req := httptest.NewRequest(http.MethodPut, "/favourites", strings.NewReader(`{}`))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(injectJWT(req.Context(), userId.String()))
		w := httptest.NewRecorder()

		// Act
		handler(w, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})
}
