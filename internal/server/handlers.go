package server

import (
	"net/http"
	"platform-go-challenge/internal/utils"
)

func getHealth(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithMessage(w, http.StatusOK, "Healthy!")
}
