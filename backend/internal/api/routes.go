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
	// config              Config,
	// tenantsStore        *TenantsStore,
	// commentsStore       *CommentsStore,
	// conversationService *ConversationService,
	// chatGPTService      *ChatGPTService,
	// authProxy           *authProxy,
) {
	// mux.Handle("/api/v1/", handleTenantsGet(logger, tenantsStore))
	// mux.Handle("/oauth2/", handleOAuth2Proxy(logger, authProxy))
	// mux.HandleFunc("/healthz", handleHealthzPlease(logger))
	mux.Handle("/login", services.HandleLoginService(logger, userStore, tokenStore))
	mux.Handle("/test", middlewares.Auth(services.HandleTestService(logger), tokenStore))
	mux.Handle("/", http.NotFoundHandler())
}
