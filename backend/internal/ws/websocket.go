package ws

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/coder/websocket"
)

// echoServer is the WebSocket echo server implementation.
type echoServer struct {
	// logf controls where logs are sent.
	logf func(f string, v ...interface{})
}

func (s echoServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true, // Disable origin verification
	})
	if err != nil {
		log.Println("Failed to accept connection:", err)
		return
	}
	defer conn.Close(websocket.StatusNormalClosure, "")

	for {
		msgType, data, err := conn.Read(r.Context())
		if websocket.CloseStatus(err) != -1 {
			log.Println("Connection closed:", err)
			return
		}
		if err != nil {
			log.Println("Read error:", err)
			return
		}

		log.Printf("Received: %s\n", data)

		err = conn.Write(r.Context(), msgType, data)
		if err != nil {
			log.Println("Write error:", err)
			return
		}
	}
}

func Run(
	logger *slog.Logger,
	getenv func(string) string,
) error {
	ws_url := getenv("WS_URL")
	l, err := net.Listen("tcp", ws_url)
	if err != nil {
		return err
	}
	logger.Info(fmt.Sprintf("listening on ws://%v", l.Addr().String()))

	s := &http.Server{
		Handler: echoServer{
			logf: log.Printf,
		},
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}
	errc := make(chan error, 1)
	go func() {
		errc <- s.Serve(l)
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	select {
	case err := <-errc:
		logger.Info(fmt.Sprintf("failed to serve: %v", err))
	case sig := <-sigs:
		logger.Info(fmt.Sprintf("terminating: %v", sig))
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	return s.Shutdown(ctx)
}
