package api

import (
	"context"
	"log/slog"
	"net/http"
	"os"
)

func NewServer(
	logger *slog.Logger,
	// config *Config,
	//
	//	anotherStore *anotherStore,
) http.Handler {
	mux := http.NewServeMux()
	addRoutes(
		mux,
		logger,
		// Config,
		// commentStore,
		// anotherStore,
	)
	var handler http.Handler = mux
	return handler
}

func Run(
	ctx context.Context,
	getenv func(string) string,
) error {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	srv := NewServer(logger)

	logger.Info("Serving on" + getenv("APP_URL"))
	http.ListenAndServe(getenv("APP_URL"), srv)

	return nil
}
