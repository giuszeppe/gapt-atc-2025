package api

import (
	"log/slog"
	"net/http"

	"github.com/giuszeppe/gatp-atc-2025/backend/internal/services"
	"github.com/giuszeppe/gatp-atc-2025/backend/internal/stores"
)

func addRoutes(
	mux *http.ServeMux,
	logger *slog.Logger,
	tokenStore stores.Store[string],
	userStore stores.UserStore,
	scenarioStore stores.ScenarioStore,
	// authProxy           *authProxy,
) {
	mux.Handle("/login", services.HandleLoginService(logger, userStore, tokenStore))
	mux.Handle("/get-scenarios", services.HandleGetScenario(logger, scenarioStore))
	mux.Handle("/post-simulation", services.HandlePostSimulation(logger, scenarioStore))
	mux.Handle("/end-simulation", services.HandleEndSimulation(logger, scenarioStore))
	mux.Handle("/get-transcripts", services.HandleGetTranscripts(logger, scenarioStore))
	mux.Handle("/get-transcripts/{id}", services.HandleGetTranscript(logger, scenarioStore))
	mux.Handle("/simulation-lobby", services.HandleMultiplayerLobbyWebsocket(logger, scenarioStore))

	mux.Handle("/", http.NotFoundHandler())
}
