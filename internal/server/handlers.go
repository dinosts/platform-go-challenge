package server

import (
	"net/http"
	"platform-go-challenge/internal/utils"
)

func GetHealth(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithMessage(w, http.StatusOK, "Healthy!")
}
