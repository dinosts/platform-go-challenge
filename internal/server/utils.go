package server

import (
	"encoding/json"
	"net/http"
)

func respondWithJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func RespondWithError(w http.ResponseWriter, status int, message string) {
	respondWithJSON(w, status, ErrorResponse{Error: message})
}

func RespondWithMessage(w http.ResponseWriter, status int, message string) {
	respondWithJSON(w, status, MessageResponse{Message: message})
}

func RespondWithData(w http.ResponseWriter, status int, data any) {
	respondWithJSON(w, status, DataResponse{Data: data})
}
