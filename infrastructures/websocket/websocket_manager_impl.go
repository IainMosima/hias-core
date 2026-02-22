package websocket

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Configure properly in production
	},
}

type webSocketManagerImpl struct {
	hub *Hub
}

// NewWebSocketManager creates a new WebSocket manager.
func NewWebSocketManager() WebSocketManager {
	return &webSocketManagerImpl{
		hub: NewHub(),
	}
}

func (m *webSocketManagerImpl) HandleConnection(w http.ResponseWriter, r *http.Request, userID string, connType ConnectionType, resourceID string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	clientID := uuid.New().String()
	client := NewClient(clientID, userID, connType, resourceID, conn, m.hub)

	m.hub.register <- client

	go client.WritePump()
	go client.ReadPump()
}

func (m *webSocketManagerImpl) SendToResource(connType ConnectionType, resourceID string, msg Message) {
	m.hub.notify <- wsNotification{
		connType:   connType,
		resourceID: resourceID,
		message:    msg,
	}
}

func (m *webSocketManagerImpl) SendToUser(userID string, msg Message) {
	m.hub.notify <- wsNotification{
		userID:  userID,
		message: msg,
	}
}

func (m *webSocketManagerImpl) Broadcast(connType ConnectionType, msg Message) {
	m.hub.broadcast <- wsBroadcast{
		connType: connType,
		message:  msg,
	}
}

func (m *webSocketManagerImpl) Start() {
	go m.hub.Run()
	log.Println("WebSocket manager started")
}

func (m *webSocketManagerImpl) Stop() {
	close(m.hub.stop)
	log.Println("WebSocket manager stopped")
}

func (m *webSocketManagerImpl) GetConnectionCount() int {
	return m.hub.GetConnectionCount()
}
