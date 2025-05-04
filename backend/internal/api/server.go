package api

import (
	"context"
	"github.com/rs/cors"
	"log/slog"
	"net/http"
	"os"

	"github.com/giuszeppe/gatp-atc-2025/backend/internal/stores"
)

func NewServer(
	logger *slog.Logger,
	tokenStore stores.Store[string],
	userStore *stores.UserStore,
	scenarioStore *stores.ScenarioStore,
	// config *Config,
) http.Handler {
	mux := http.NewServeMux()
	addRoutes(
		mux,
		logger,
		tokenStore,
		*userStore,
		*scenarioStore,
		// Config,
	)
	var handler http.Handler = mux
	handler = cors.Default().Handler(handler)
	return handler
}

func Run(
	ctx context.Context,
	getenv func(string) string,
	tokenStore stores.Store[string],
	userStore *stores.UserStore,
	scenarioStore *stores.ScenarioStore,
) error {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	srv := NewServer(logger, tokenStore, userStore, scenarioStore)

	logger.Info("Serving on" + getenv("APP_URL"))
	http.ListenAndServe(getenv("APP_URL"), srv)

	return nil
}
