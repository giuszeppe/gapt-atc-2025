package ws

import (
	"bytes"
	"context"
	"fmt"
	"github.com/coder/websocket"
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

func UpgradeConnectionToLobbyWebsocket(w http.ResponseWriter, r *http.Request, store stores.ScenarioStore) {
	lobbyCode := r.URL.Query().Get("lobby")
	if lobbyCode == "" {
		http.Error(w, "lobby code required", http.StatusBadRequest)
		return
	}

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{OriginPatterns: []string{"*"}})
	if err != nil {
		fmt.Println("WebSocket accept error:", err)
		return
	}

	client := &Client{
		conn: conn,
		send: make(chan []byte, 64),
	}

	lobby := getOrCreateLobby(lobbyCode)
	addClientToLobby(lobby, client)

	// send existing messages to client
	for _, message := range lobby.Messages {
		client.send <- []byte(message.Role + "\n" + message.Text)
	}

	go clientWriter(client)
	clientReader(lobby, client, r, store)
}

func addClientToLobby(lobby *Lobby, client *Client) {
	lobby.mutex.Lock()
	defer lobby.mutex.Unlock()
	lobby.clients[client] = true
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

func broadcastToLobby(lobby *Lobby, data []byte, sender *Client) {
	lobby.mutex.RLock()
	defer lobby.mutex.RUnlock()

	message := bytes.Split(data, []byte("\n"))
	if bytes.Equal(message[0], []byte("text")) {
		role := message[1]
		content := message[2]
		lobby.Messages = append(lobby.Messages, stores.Message{
			Role: string(role),
			Text: string(content),
		})
		fmt.Println("Lobby broadcast and appended messaage:", lobby.Code, lobby.Messages)
	}

	for client := range lobby.clients {
		if client != sender {
			select {
			case client.send <- data:
			default:
				fmt.Println("Dropping message due to full send buffer")
			}
		}
	}
}
