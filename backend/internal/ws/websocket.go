package ws

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/coder/websocket"
)

type Client struct {
	conn *websocket.Conn
	send chan []byte
}

type Lobby struct {
	clients map[*Client]bool
	mutex   sync.RWMutex
}

var lobbies = make(map[string]*Lobby)
var globalMutex sync.RWMutex

func getOrCreateLobby(code string) *Lobby {
	globalMutex.Lock()
	defer globalMutex.Unlock()

	lobby, exists := lobbies[code]
	if !exists {
		lobby = &Lobby{
			clients: make(map[*Client]bool),
		}
		lobbies[code] = lobby
	}
	return lobby
}

func UpgradeConnectionToLobbyWebsocket(w http.ResponseWriter, r *http.Request) {
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

	go clientWriter(client)
	clientReader(lobby, client, r)
}

func addClientToLobby(lobby *Lobby, client *Client) {
	lobby.mutex.Lock()
	defer lobby.mutex.Unlock()
	lobby.clients[client] = true
	fmt.Println("Client joined lobby")
}

func removeClientFromLobby(lobby *Lobby, client *Client) {
	lobby.mutex.Lock()
	defer lobby.mutex.Unlock()
	delete(lobby.clients, client)
	close(client.send)
	fmt.Println("Client left lobby")
}

func clientReader(lobby *Lobby, client *Client, r *http.Request) {
	defer func() {
		removeClientFromLobby(lobby, client)
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
