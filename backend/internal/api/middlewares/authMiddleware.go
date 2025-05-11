package middlewares

import (
	"net/http"

	"github.com/giuszeppe/gatp-atc-2025/backend/internal/encoder"
	"github.com/giuszeppe/gatp-atc-2025/backend/internal/stores"
)

func Auth(h http.Handler, tokenStore stores.TokenStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v := r.Header.Get("Authorization")
		_, err := tokenStore.GetUserByToken(v)
		if err != nil {
			encoder.EncodeError(w, http.StatusUnauthorized, nil, "User not authorized")
			return
		}
		h.ServeHTTP(w, r)
	})
}
