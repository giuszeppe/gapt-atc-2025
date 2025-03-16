package services

import (
	"log/slog"
	"net/http"
)

func HandleTestService(logger *slog.Logger) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// use thing to handle request
			logger.Info("handleSomething")
		},
	)
}
