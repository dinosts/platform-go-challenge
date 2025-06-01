package utils_test

import (
	"net/http"
	"net/url"
	"platform-go-challenge/internal/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPaginationQuery(t *testing.T) {
	t.Run("Should return defaults when no query params are given", func(t *testing.T) {
		// Arrange
		r := &http.Request{URL: &url.URL{RawQuery: ""}}

		// Act
		pageSize, pageNumber, err := utils.GetPaginationQuery(r, 20, 2)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 20, pageSize)
		assert.Equal(t, 2, pageNumber)
	})

	t.Run("Should parse valid pageSize and pageNumber from query", func(t *testing.T) {
		// Arrange
		r := &http.Request{URL: &url.URL{RawQuery: "pageSize=50&pageNumber=3"}}

		// Act
		pageSize, pageNumber, err := utils.GetPaginationQuery(r, 20, 2)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 50, pageSize)
		assert.Equal(t, 3, pageNumber)
	})

	t.Run("Should return error when pageSize exceeds max", func(t *testing.T) {
		// Arrange
		r := &http.Request{URL: &url.URL{RawQuery: "pageSize=500"}}

		// Act
		_, _, err := utils.GetPaginationQuery(r, 20, 2)

		// Assert
		assert.ErrorIs(t, err, utils.ErrInvalidPageSize)
	})

	t.Run("Should return error when pageSize is not a number", func(t *testing.T) {
		// Arrange
		r := &http.Request{URL: &url.URL{RawQuery: "pageSize=abc"}}

		// Act
		_, _, err := utils.GetPaginationQuery(r, 20, 2)

		// Assert
		assert.ErrorIs(t, err, utils.ErrCouldNotParsePageSize)
	})

	t.Run("Should return error when pageNumber is not a number", func(t *testing.T) {
		// Arrange
		r := &http.Request{URL: &url.URL{RawQuery: "pageNumber=abc"}}

		// Act
		_, _, err := utils.GetPaginationQuery(r, 20, 2)

		// Assert
		assert.ErrorIs(t, err, utils.ErrCouldNotParsePageNumber)
	})

	t.Run("Should return error when pageNumber is negative", func(t *testing.T) {
		// Arrange
		r := &http.Request{URL: &url.URL{RawQuery: "pageNumber=-1"}}

		// Act
		_, _, err := utils.GetPaginationQuery(r, 20, 2)

		// Assert
		assert.ErrorIs(t, err, utils.ErrInvalidPageNumber)
	})

	t.Run("Should return error when pageSize is too small", func(t *testing.T) {
		// Arrange
		r := &http.Request{URL: &url.URL{RawQuery: "pageSize=0"}}

		// Act
		_, _, err := utils.GetPaginationQuery(r, 20, 2)

		// Assert
		assert.ErrorIs(t, err, utils.ErrInvalidPageSize)
	})
}
