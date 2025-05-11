package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/coder/websocket"
	"github.com/giuszeppe/gatp-atc-2025/backend/internal/encoder"
	"github.com/giuszeppe/gatp-atc-2025/backend/internal/stores"
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

func UpgradeConnectionToLobbyWebsocket(w http.ResponseWriter, r *http.Request, store stores.ScenarioStore, tokenStore *stores.TokenStore) {
	lobbyCode := r.URL.Query().Get("lobby")
	if lobbyCode == "" {
		http.Error(w, "lobby code required", http.StatusBadRequest)
		return
	}

	lobby := getOrCreateLobby(lobbyCode)

	token := r.Header.Get("Authorization")

	user, err := tokenStore.GetUserByToken(token)
	role := ""
	if err != nil {
		encoder.EncodeError(w, 401, err, err.Error())
	}
	userId := user.ID
	simulation, err := store.GetSimulationByLobbyCode(lobbyCode)
	if err != nil {
		encoder.EncodeError(w, 500, err, err.Error())
		return
	}

	if simulation.TowerUserId != userId && simulation.AircraftUserId != userId { //user is not in simulation
		if simulation.TowerUserId == -1 {
			role = "tower"
		} else {
			role = "aircraft"
		}
		err = store.UpdateSimulationRoleIds(simulation.Id, userId, role)
		if err != nil {
			encoder.EncodeError(w, 500, err, err.Error())
			return
		}
	} else if simulation.TowerUserId == userId {
		role = "tower"
	} else if simulation.AircraftUserId == userId {
		role = "aircraft"
	}

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

	client.send <- []byte(fmt.Sprintf(`{"type":"role","content":"%s"}`, role))

	addClientToLobby(lobby, client)

	// send existing messages to client
	for _, message := range lobby.Messages {
		msgJson, err := json.Marshal(message)
		if err != nil {
			fmt.Println(err)
			return
		}
		client.send <- msgJson
	}

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
