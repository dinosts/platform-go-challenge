package user

import (
	"errors"
	"net/http"
	"platform-go-challenge/internal/utils"
)

type UserLoginDependencies struct {
	UserService UserService
}

func UserLoginHandler(dependencies UserLoginDependencies) http.HandlerFunc {
	return utils.BodyValidator[UserLoginRequestBody](
		func(w http.ResponseWriter, r *http.Request) {
			body, ok := utils.GetParsedBody[UserLoginRequestBody](r)
			if !ok {
				utils.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
				return
			}

			token, expiresAt, err := dependencies.UserService.LoginUser(body.Email, body.Password)
			if err != nil {
				if errors.Is(err, ErrLoginFailed) {
					utils.RespondWithError(w, http.StatusUnauthorized, "Invalid email or password")
					return
				}

				utils.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
				return
			}

			response := UserLoginResponseBody{
				Token:     token,
				ExpiresAt: expiresAt,
			}

			utils.RespondWithData(w, http.StatusOK, response)
		},
	)
}
