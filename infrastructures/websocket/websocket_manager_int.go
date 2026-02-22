package websocket

import "net/http"

// WebSocketManager manages WebSocket connections and messaging.
type WebSocketManager interface {
	HandleConnection(w http.ResponseWriter, r *http.Request, userID string, connType ConnectionType, resourceID string)
	SendToResource(connType ConnectionType, resourceID string, msg Message)
	SendToUser(userID string, msg Message)
	Broadcast(connType ConnectionType, msg Message)
	Start()
	Stop()
	GetConnectionCount() int
}
