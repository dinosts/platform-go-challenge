package server

import (
	"net/http"
)

func getHealth(w http.ResponseWriter, r *http.Request) {
	RespondWithMessage(w, http.StatusOK, "Healthy!")
}
