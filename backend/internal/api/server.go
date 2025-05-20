package api

import (
	"fmt"
	"github.com/rs/cors"
	"log/slog"
	"net/http"
	"os"

	"github.com/giuszeppe/gatp-atc-2025/backend/internal/stores"
)

func NewServer(
	logger *slog.Logger,
	tokenStore *stores.TokenStore,
	userStore *stores.UserStore,
	scenarioStore *stores.ScenarioStore,
) http.Handler {
	mux := http.NewServeMux()
	addRoutes(
		mux,
		logger,
		tokenStore,
		*userStore,
		*scenarioStore,
	)
	var handler http.Handler = mux
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		Debug:            false,
	})
	handler = c.Handler(handler)
	return handler
}

func Run(
	getenv func(string) string,
	tokenStore *stores.TokenStore,
	userStore *stores.UserStore,
	scenarioStore *stores.ScenarioStore,
) error {
	// Open logfile
	logFile, err := os.OpenFile("gatp.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	logger := slog.New(slog.NewJSONHandler(logFile, nil))
	srv := NewServer(logger, tokenStore, userStore, scenarioStore)

	logger.Info("Serving on " + getenv("APP_URL"))
	fmt.Println("Serving on " + getenv("APP_URL"))
	err = http.ListenAndServe(getenv("APP_URL"), srv)
	if err != nil {
		return err
	}

	return nil
}
