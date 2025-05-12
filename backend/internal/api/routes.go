package api

import (
	"github.com/giuszeppe/gatp-atc-2025/backend/internal/api/middlewares"
	"log/slog"
	"net/http"

	"github.com/giuszeppe/gatp-atc-2025/backend/internal/services"
	"github.com/giuszeppe/gatp-atc-2025/backend/internal/stores"
)

func addRoutes(
	mux *http.ServeMux,
	logger *slog.Logger,
	tokenStore *stores.TokenStore,
	userStore stores.UserStore,
	scenarioStore stores.ScenarioStore,
	// authProxy           *authProxy,
) {
	mux.Handle("/login", services.HandleLoginService(logger, userStore, tokenStore))
	mux.Handle("/get-scenarios", middlewares.Auth(services.HandleGetScenario(logger, scenarioStore), tokenStore))
	mux.Handle("/post-simulation", middlewares.Auth(services.HandlePostSimulation(logger, scenarioStore, tokenStore), tokenStore))
	mux.Handle("/end-simulation", middlewares.Auth(services.HandleEndSimulation(logger, scenarioStore), tokenStore))
	mux.Handle("/get-transcripts", middlewares.Auth(services.HandleGetTranscripts(logger, scenarioStore), tokenStore))
	mux.Handle("/get-transcripts/{id}", middlewares.Auth(services.HandleGetTranscript(logger, scenarioStore), tokenStore))
	mux.Handle("/simulation-lobby", services.HandleMultiplayerLobbyWebsocket(logger, scenarioStore, tokenStore))
	mux.Handle("/", http.NotFoundHandler())
}
