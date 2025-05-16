package encoder

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func Encode[T any](w http.ResponseWriter, r *http.Request, status int, v T, logger *slog.Logger) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	response := APIResponse{
		Success: true,
		Data:    v,
		Error:   "",
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		EncodeError(w, http.StatusInternalServerError, nil, err.Error(), logger)
	}
}

func Decode[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("decode json: %w", err)
	}
	return v, nil
}

func EncodeError(w http.ResponseWriter, status int, data interface{}, errMsg string, logger *slog.Logger) error {
	w.Header().Set("Content-Type", "application/json")
	logger.Error(errMsg)
	if status != http.StatusOK {
		w.WriteHeader(status) // statusOK is implicit when one calls write
	}

	response := APIResponse{
		Success: errMsg == "",
		Data:    data,
		Error:   errMsg,
	}

	return json.NewEncoder(w).Encode(response)
}
