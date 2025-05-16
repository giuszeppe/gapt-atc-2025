package ws

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/coder/websocket"
	"github.com/giuszeppe/gatp-atc-2025/backend/internal/stores"
	"log/slog"
	"net/http"
	"sync"
)

type Client struct {
	conn *websocket.Conn
	send chan []byte
}

type Lobby struct {
	clients  map[*Client]bool
	mutex    sync.RWMutex
	Messages []stores.Message
	Code     string
}

type InitializationMessage struct {
	InputType     string             `json:"input_type"`
	Role          string             `json:"role"`
	Steps         []stores.Step      `json:"steps"`
	ExtendedSteps []stores.Step      `json:"extended_steps"`
	SimulationId  int                `json:"simulation_id"`
	Messages      []WebsocketMessage `json:"messages"`
}

var lobbies = make(map[string]*Lobby)
var globalMutex sync.RWMutex

func getOrCreateLobby(code string) *Lobby {
	globalMutex.Lock()
	defer globalMutex.Unlock()

	lobby, exists := lobbies[code]
	if !exists {
		lobby = &Lobby{
			clients:  make(map[*Client]bool),
			Messages: []stores.Message{},
			Code:     code,
		}
		lobbies[code] = lobby
	}
	return lobby
}

func getUserRole(simulation stores.Simulation, userId int, store stores.ScenarioStore) (string, error) {
	var role string
	if simulation.TowerUserId != userId && simulation.AircraftUserId != userId { //user is not in simulation
		if simulation.TowerUserId == -1 {
			role = "tower"
		} else {
			role = "aircraft"
		}
	} else if simulation.TowerUserId == userId {
		role = "tower"
	} else if simulation.AircraftUserId == userId {
		role = "aircraft"
	}
	err := store.UpdateSimulationRoleIds(simulation.Id, userId, role)
	return role, err
}

func UpgradeConnectionToLobbyWebsocket(logger *slog.Logger, w http.ResponseWriter, r *http.Request, store stores.ScenarioStore, tokenStore *stores.TokenStore) {
	lobbyCode := r.URL.Query().Get("lobby")
	if lobbyCode == "" {
		http.Error(w, "lobby code required", http.StatusBadRequest)
		logger.Error("Lobby code required", errors.New("lobby code required"))
		return
	}

	lobby := getOrCreateLobby(lobbyCode)

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{OriginPatterns: []string{"*"}})
	conn.SetReadLimit(-1)
	if err != nil {
		fmt.Println("WebSocket accept error:", err)
		return
	}

	client := &Client{
		conn: conn,
		send: make(chan []byte, 64),
	}

	initMsg := InitializationMessage{}

	// Get Token from Connection
	_, data, err := client.conn.Read(r.Context())
	authToken := string(data)
	logger.Info("Auth token", "token", authToken)

	user, err := tokenStore.GetUserByToken(authToken)
	if err != nil {
		logger.Error("Error getting user by token", "error", err)
		client.conn.Close(websocket.StatusInternalError, "User not authorized")
		return
	}
	userId := user.ID

	// Get Simulation by Lobby Code
	simulation, err := store.GetSimulationByLobbyCode(lobbyCode)
	if err != nil {
		// close ws connection and return error
		logger.Error("Error getting simulation by lobby code", "error", err)
		client.conn.Close(websocket.StatusInternalError, "Simulation not found")
		return
	}
	initMsg.SimulationId = simulation.Id
	initMsg.InputType = simulation.InputType

	// Get User Role
	role, err := getUserRole(simulation, userId, store)
	if err != nil {
		logger.Error("Error getting user role", "error", err)
		return
	}

	logger.Info("User role", "userId", userId, "role", role)
	initMsg.Role = role

	// Get Simulation Steps
	steps, err := store.GetScenarioStepsForId(simulation.ScenarioId)
	if err != nil {
		logger.Error("Error getting simulation steps", "error", err)
		return
	}
	initMsg.Steps = steps[0]
	initMsg.ExtendedSteps = steps[1]

	addClientToLobby(lobby, client)

	// send existing messages to client
	initMsg.Messages = []WebsocketMessage{}
	for _, message := range lobby.Messages {
		initMsg.Messages = append(initMsg.Messages, WebsocketMessage{
			Type:    "text",
			Content: json.RawMessage(message.Text),
			Role:    message.Role,
		})
	}

	// send initialization message to client
	initJson, err := json.Marshal(initMsg)
	if err != nil {
		logger.Error("Error marshalling init message", "error", err)
		return
	}
	client.send <- []byte(fmt.Sprintf(`{"type":"init","content":%s}`, initJson))

	go clientWriter(client)
	clientReader(lobby, client, r, store)
}

type NewClientMsg struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

func addClientToLobby(lobby *Lobby, newClient *Client) {
	lobby.mutex.Lock()
	defer lobby.mutex.Unlock()
	lobby.clients[newClient] = true

	for client := range lobby.clients {
		if client != newClient {
			msg := NewClientMsg{
				Type:    "newClient",
				Content: "Added new client to lobby",
			}
			msgJson, err := json.Marshal(msg)
			if err != nil {
				fmt.Println(err)
			}
			select {
			case client.send <- msgJson:
			}
		}
	}

	fmt.Println("Client joined lobby")
}

func removeClientFromLobby(lobby *Lobby, client *Client, store stores.ScenarioStore) {
	lobby.mutex.Lock()
	defer lobby.mutex.Unlock()
	delete(lobby.clients, client)
	close(client.send)
	fmt.Println("Client left lobby")
	if len(lobby.clients) == 0 {
		err := store.AddTranscriptToSimulationUsingLobbyCode(lobby.Code, lobby.Messages)
		if err != nil {
			return
		}
	}
}

func clientReader(lobby *Lobby, client *Client, r *http.Request, store stores.ScenarioStore) {
	defer func() {
		removeClientFromLobby(lobby, client, store)
		client.conn.Close(websocket.StatusNormalClosure, "closing")
	}()

	for {
		_, data, err := client.conn.Read(r.Context())
		if err != nil {
			fmt.Println("Read error:", err)
			return
		}
		broadcastToLobby(lobby, data, client)
	}
}

func clientWriter(client *Client) {
	for data := range client.send {
		err := client.conn.Write(context.Background(), websocket.MessageBinary, data)
		if err != nil {
			fmt.Println("Write error:", err)
			return
		}
	}
}

type WebsocketMessage struct {
	Type    string          `json:"type"`
	Content json.RawMessage `json:"content"`
	Role    string          `json:"role"`
}

func broadcastToLobby(lobby *Lobby, data []byte, sender *Client) {
	lobby.mutex.RLock()
	defer lobby.mutex.RUnlock()
	wsMsg := WebsocketMessage{}
	err := json.Unmarshal(data, &wsMsg)
	if err != nil {
		fmt.Println("Unmarshal error:", err)
	}

	if wsMsg.Type == "text" {
		lobby.Messages = append(lobby.Messages, stores.Message{
			Role: wsMsg.Role,
			Text: string(wsMsg.Content),
		})
		fmt.Println("Lobby broadcast and appended messaage:", lobby.Code, lobby.Messages)
	}
	for client := range lobby.clients {
		fmt.Println("Lobby broadcast:", lobby.Code)
		if client != sender {
			select {
			case client.send <- data:
			default:
				fmt.Println("Dropping message due to full send buffer")
			}
		}
	}
}
