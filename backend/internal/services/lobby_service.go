package services

import (
	"errors"
	"github.com/giuszeppe/gatp-atc-2025/backend/internal/stores"
	"github.com/giuszeppe/gatp-atc-2025/backend/internal/ws"
	"log/slog"
	"math/rand/v2"
	"net/http"
)

func HandleMultiplayerLobbyWebsocket(logger *slog.Logger, scenarioStore stores.ScenarioStore, tokenStore *stores.TokenStore) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			ws.UpgradeConnectionToLobbyWebsocket(logger, w, r, scenarioStore, tokenStore)
		},
	)
}
func GenerateLobbyCode(scenarioStore stores.ScenarioStore) (string, error) {
	tryCount := 0
	for tryCount < 30 {
		code := getRandomCode()
		exist, err := doesCodeExists(code, scenarioStore)
		if err != nil {
			return "", err
		}
		if !exist {
			return code, nil
		}
		tryCount++
	}
	return "", errors.New("could not generate lobby code")
}

func getRandomCode() string {
	code := ""
	letters := "QWERTYUIOPASDFGHJKLZXCVBNM"
	numbers := "0123456789"
	for i := 0; i < 6; i++ {
		isLetter := rand.IntN(1)
		if isLetter == 0 {
			code += string(letters[rand.IntN(len(letters))])
		} else {
			code += string(numbers[rand.IntN(len(numbers))])
		}
	}
	return code
}

func doesCodeExists(code string, store stores.ScenarioStore) (bool, error) {
	exist, err := store.DoesLobbyCodeExist(code)
	if err != nil {
		return false, err
	}
	return exist, nil
}
