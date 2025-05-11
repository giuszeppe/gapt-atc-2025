package services

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/giuszeppe/gatp-atc-2025/backend/internal/auth"
	"github.com/giuszeppe/gatp-atc-2025/backend/internal/encoder"
	"github.com/giuszeppe/gatp-atc-2025/backend/internal/stores"
)

type RequestData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TokenResponse struct {
	Token string
}

func HandleLoginService(logger *slog.Logger, userStore stores.UserStore, tokenStore *stores.TokenStore) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// Check if the request is a POST request
			if r.Method != http.MethodPost {
				encoder.EncodeError(w, http.StatusMethodNotAllowed, nil, "Method not allowed")
				return
			}

			data, err := encoder.Decode[stores.User](r)
			if err != nil {
				encoder.EncodeError(w, http.StatusBadRequest, nil, err.Error())
				return
			}

			// Get the username and password from the form
			username := strings.TrimSpace(data.Username)
			password := data.Password

			logger.Info(username + " " + password)

			// Verify login credentials
			if user, ok := auth.Login(userStore, username, password); ok {
				var token string
				t, _ := auth.RandomHex(20)
				token += "Bearer " + t
				user.Token = token
				user.Password = ""
				err := tokenStore.Store(user)
				if err != nil {
					logger.Error(err.Error())
					return
				}
				response := TokenResponse{Token: token}

				encoder.Encode(w, r, http.StatusOK, response)
				return
			} else {
				encoder.EncodeError(w, http.StatusUnauthorized, nil, "Wrong Credentials")
			}
			return

		},
	)
}
