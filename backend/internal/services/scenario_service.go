package services

import (
	"github.com/giuszeppe/gatp-atc-2025/backend/internal/encoder"
	"github.com/giuszeppe/gatp-atc-2025/backend/internal/stores"
	"log/slog"
	"net/http"
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
	Id                        int    `json:"scenario_id"`
	InputType                 string `json:"input_type"`                  // block, text, speech
	ScenarioType              string `json:"scenario_type"`               // takeoff, enroute, landing
	Role                      string `json:"role"`                        // tower, aircraft
	SimulationAdvancementType string `json:"simulation_advancement_type"` // continuous, steps
	Mode                      string `json:"mode"`                        // single, multi
}

type PostScenarioResponse struct {
	Steps      [][]stores.Step
	Simulation stores.Simulation
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

			if data.Mode == "single" {
				simulation, err := scenarioStore.StoreSimulation(
					1, // logged-in user id
					data.Id,
					data.Role,
					data.InputType,
					data.ScenarioType,
					data.SimulationAdvancementType,
					data.Mode,
				)
				if err != nil {
					encoder.EncodeError(w, 500, err, err.Error())
				}

				encoder.Encode(w, r, 200, PostScenarioResponse{Steps: steps, Simulation: simulation})

			} else {
				return
			}

		},
	)
}

type EndSimulationRequest struct {
	SimulationId int `json:"simulation_id"`
	Transcripts  []stores.Transcript
}

func HandleEndSimulation(logger *slog.Logger, scenarioStore stores.ScenarioStore) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
		},
	)
}
