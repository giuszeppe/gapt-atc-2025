package api

import (
	"log/slog"
	"net/http"

	"github.com/giuszeppe/gatp-atc-2025/backend/internal/api/middlewares"
	"github.com/giuszeppe/gatp-atc-2025/backend/internal/services"
	"github.com/giuszeppe/gatp-atc-2025/backend/internal/stores"
)

func addRoutes(
	mux *http.ServeMux,
	logger *slog.Logger,
	tokenStore stores.Store[string],
	userStore stores.UserStore,
	// authProxy           *authProxy,
) {
	mux.Handle("/login", services.HandleLoginService(logger, userStore, tokenStore))
	mux.Handle("/test", middlewares.Auth(services.HandleTestService(logger), tokenStore))
	mux.Handle("/", http.NotFoundHandler())
}
