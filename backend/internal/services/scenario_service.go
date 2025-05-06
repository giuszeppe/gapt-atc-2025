package services

import (
	"github.com/giuszeppe/gatp-atc-2025/backend/internal/encoder"
	"github.com/giuszeppe/gatp-atc-2025/backend/internal/stores"
	"log/slog"
	"net/http"
	"strconv"
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
	Mode                      string `json:"mode"`                        // singleplayer, multiplayer
}

type PostScenarioResponse struct {
	Steps      [][]stores.Step   `json:"steps"`
	Simulation stores.Simulation `json:"simulation"`
	LobbyCode  string            `json:"lobby_code,omitempty"`
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
				return
			}
			simulation, err := scenarioStore.StoreSimulation(
				data.Id,
				data.Role,
				data.InputType,
				data.ScenarioType,
				data.SimulationAdvancementType,
				data.Mode,
			)
			if err != nil {
				encoder.EncodeError(w, 500, err, err.Error())
				return
			}

			if data.Mode == "singleplayer" {
				encoder.Encode(w, r, 200, PostScenarioResponse{Steps: steps, Simulation: simulation})
			} else if data.Mode == "multiplayer" {
				lobbyCode, err := GenerateLobbyCode(scenarioStore)
				if err != nil {
					encoder.EncodeError(w, 500, err, err.Error())
					return
				}

				encoder.Encode(w, r, 200, PostScenarioResponse{Steps: steps, Simulation: simulation, LobbyCode: lobbyCode})
				return
			}

		},
	)
}

type EndSimulationRequest struct {
	SimulationId int              `json:"simulation_id"`
	Messages     []stores.Message `json:"messages,omitempty"`
}

func HandleEndSimulation(logger *slog.Logger, scenarioStore stores.ScenarioStore) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// get simulation params
			data, err := encoder.Decode[EndSimulationRequest](r)
			if err != nil {
				encoder.EncodeError(w, http.StatusBadRequest, nil, err.Error())
				return
			}

			// fetch scenario steps
			err = scenarioStore.EndSimulation(data.SimulationId, data.Messages)
			if err != nil {
				encoder.EncodeError(w, http.StatusInternalServerError, err, err.Error())
				return
			}
			encoder.Encode(w, r, http.StatusNoContent, "")
		},
	)
}

func HandleGetTranscripts(logger *slog.Logger, scenarioStore stores.ScenarioStore) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

			// fetch scenario steps
			transcripts, err := scenarioStore.GetGroupedTranscripts()
			if err != nil {
				encoder.EncodeError(w, http.StatusInternalServerError, err, err.Error())
				return
			}
			encoder.Encode(w, r, http.StatusOK, transcripts)
		},
	)
}
func HandleGetTranscript(logger *slog.Logger, scenarioStore stores.ScenarioStore) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			simulationId, err := strconv.Atoi(r.PathValue("id"))
			if err != nil {
				logger.Error(err.Error())
				encoder.EncodeError(w, http.StatusBadRequest, nil, err.Error())
				return
			}

			// fetch scenario steps
			transcripts, err := scenarioStore.GetTranscriptBySimulationId(simulationId)
			if err != nil {
				logger.Error(err.Error())
				encoder.EncodeError(w, http.StatusInternalServerError, err, err.Error())
				return
			}
			encoder.Encode(w, r, http.StatusOK, transcripts)
		},
	)
}
