package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

const parsedBodyKey string = "parsedBody"

func BodyValidator[T any](next http.HandlerFunc) http.HandlerFunc {
	validate := validator.New()
	return func(w http.ResponseWriter, r *http.Request) {
		var parsedBody T

		if err := json.NewDecoder(r.Body).Decode(&parsedBody); err != nil {
			RespondWithError(w, http.StatusBadRequest, "Invalid JSON body")
			return
		}

		if err := validate.Struct(parsedBody); err != nil {
			errs := err.(validator.ValidationErrors)
			message := fmt.Sprintf("Body Validation Failed, %s", errs)
			RespondWithError(w, http.StatusBadRequest, message)
			return
		}

		ctx := context.WithValue(r.Context(), parsedBodyKey, &parsedBody)
		next(w, r.WithContext(ctx))
	}
}

func GetParsedBody[T any](r *http.Request) (T, bool) {
	body, ok := r.Context().Value(parsedBodyKey).(*T)

	if !ok || body == nil {
		var empty T
		return empty, false
	}

	return *body, true
}
