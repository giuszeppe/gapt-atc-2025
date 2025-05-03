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

type PostScenarioRequest struct {
	Id int
}

func HandlePostSimulation(logger *slog.Logger, scenarioStore stores.ScenarioStore) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// get simulation params
			data, err := encoder.Decode[PostScenarioRequest](r)
			if err != nil {
				encoder.EncodeError(w, http.StatusBadRequest, nil, err.Error())
				return
			}

			// fetch scenario steps
			steps, err := scenarioStore.GetScenarioStepsForId(data.Id)
			if err != nil {
				encoder.EncodeError(w, 500, err, err.Error())
			}

			// return steps
			encoder.Encode(w, r, 200, steps)

		},
	)
}
