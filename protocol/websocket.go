package protocol

import (
	"MyTransfer/apps/websocket/impl"
	"MyTransfer/conf"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewWebSocketService() *WebSocketService {
	s := &WebSocketService{
		l: zap.L().Named("WebSocket Service"),
	}
	s.server = &http.Server{
		Addr:    conf.C().WEBSOCKET.SocketAddr(),
		Handler: http.HandlerFunc(s.handleConnections),
	}
	return s
}

type WebSocketService struct {
	server *http.Server
	l      logger.Logger
}

func (s *WebSocketService) handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.l.Fatal(err)
	}
	// Make sure we close the connection when the function returns
	defer ws.Close()
	// Initialize a new Connection
	conn, err := impl.InitConnection(ws)
	if err != nil {
		s.l.Fatal("init connection:", err)
	}

	// Loop indefinitely
	for {
		// Read message from browser
		msg, err := conn.ReadMessage()
		if err != nil {
			s.l.Println("read:", err)
			break
		}

		// Print the message to the console
		s.l.Printf("recv: %s", msg)

		// TODO: Process the message
		// For now, we'll just echo the same message back
		processedMsg := msg

		// Write message back to browser
		err = conn.WriteMessage(processedMsg)
		if err != nil {
			s.l.Println("write:", err)
			break
		}
	}
}

func (s *WebSocketService) Start() error {
	s.l.Info("WebSocket Service start")
	if err := s.server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			s.l.Info("Websocket service stopped success")
			return nil
		}
		return fmt.Errorf("start Websocket service error, %s", err.Error())
	}
	return nil
}

func (s *WebSocketService) Stop() {
	s.l.Info("start close websocket service")
	if err := s.server.Close(); err != nil {
		s.l.Warnf("close websocket service error, %s", err)
	}
}
