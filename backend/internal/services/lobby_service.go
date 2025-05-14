package services

import (
	"github.com/giuszeppe/gatp-atc-2025/backend/internal/stores"
	"github.com/giuszeppe/gatp-atc-2025/backend/internal/ws"
	"log/slog"
	"net/http"
)

func HandleMultiplayerLobbyWebsocket(logger *slog.Logger, scenarioStore stores.ScenarioStore, tokenStore *stores.TokenStore) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			ws.UpgradeConnectionToLobbyWebsocket(logger, w, r, scenarioStore, tokenStore)
		},
	)
}
