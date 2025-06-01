package utils

import (
	"errors"
	"net/http"
	"strconv"
)

type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
	MaxPage  int `json:"maxPage"`
}

type PaginatedDataResponse[T any] struct {
	Data       T          `json:"data"`
	Pagination Pagination `json:"pagination"`
}

var (
	ErrCouldNotParsePageSize   = errors.New("could not parse page size")
	ErrCouldNotParsePageNumber = errors.New("could not parse page number")
	ErrInvalidPageSize         = errors.New("page size out of bounds")
	ErrInvalidPageNumber       = errors.New("page number out of bounds")
)

func GetPaginationQuery(r *http.Request, defaultPageSize int, defaultPageNumber int) (int, int, error) {
	maxPageSize := 100
	minPageSize := 1
	minPage := 0

	// Apply limits to defaults
	if defaultPageSize > maxPageSize {
		defaultPageSize = maxPageSize
	}
	if defaultPageSize < minPageSize {
		defaultPageSize = minPageSize
	}
	if defaultPageNumber < minPage {
		defaultPageNumber = minPage
	}

	pageSize := defaultPageSize
	pageNumber := defaultPageNumber

	query := r.URL.Query()

	pageSizeInQuery := query.Get("pageSize")
	pageNumberInQuery := query.Get("pageNumber")

	if pageSizeInQuery != "" {
		parsedPS, err := strconv.Atoi(pageSizeInQuery)
		if err != nil {
			return 0, 0, ErrCouldNotParsePageSize
		}
		if parsedPS < minPageSize || parsedPS > maxPageSize {
			return 0, 0, ErrInvalidPageSize
		}
		pageSize = parsedPS
	}

	if pageNumberInQuery != "" {
		parsedPN, err := strconv.Atoi(pageNumberInQuery)
		if err != nil {
			return 0, 0, ErrCouldNotParsePageNumber
		}
		if parsedPN < minPage {
			return 0, 0, ErrInvalidPageNumber
		}
		pageNumber = parsedPN
	}

	return pageSize, pageNumber, nil
}

func RespondWithPaginatedData[T any](w http.ResponseWriter, status int, data T, pagination Pagination) {
	respondWithJSON(w, status, PaginatedDataResponse[T]{Data: data, Pagination: pagination})
}
