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
	tokenStore stores.Store[string],
	userStore stores.UserStore,
	scenarioStore stores.ScenarioStore,
	// authProxy           *authProxy,
) {
	mux.Handle("/login", middlewares.UseCORS(services.HandleLoginService(logger, userStore, tokenStore)))
	mux.Handle("/get-scenarios", middlewares.UseCORS(services.HandleGetScenario(logger, scenarioStore)))
	mux.Handle("/post-simulation", middlewares.UseCORS(services.HandlePostSimulation(logger, scenarioStore)))
	mux.Handle("/", http.NotFoundHandler())
}
