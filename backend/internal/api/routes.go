package api

import (
	"log/slog"
	"net/http"

	"github.com/giuszeppe/gatp-atc-2025/backend/internal/services"
)

func addRoutes(
	mux *http.ServeMux,
	logger *slog.Logger,
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
    mux.Handle("/test",services.HandleTestService(logger))
	mux.Handle("/", http.NotFoundHandler())
}
