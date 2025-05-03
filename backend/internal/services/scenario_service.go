package services

import (
	"log/slog"
	"net/http"

	"github.com/giuszeppe/gatp-atc-2025/backend/internal/encoder"
	"github.com/giuszeppe/gatp-atc-2025/backend/internal/stores"
)

func HandleGetScenario(logger *slog.Logger, scenarioStore stores.ScenarioStore) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			scenarioType := r.URL.Query().Get("type")
			scenarios, err := scenarioStore.View(scenarioType)
			if err != nil {
				encoder.EncodeError(w, 500, err, err.Error())
			}
			encoder.Encode(w, r, 200, scenarios)

		},
	)
}
func HandlePostSimulation(logger *slog.Logger, scenarioStore stores.ScenarioStore) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

		},
	)
}
