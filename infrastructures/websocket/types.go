package websocket

// ConnectionType represents WebSocket connection categories.
type ConnectionType string

const (
	ConnectionTypeDashboard ConnectionType = "dashboard"
	ConnectionTypeClaimQueue ConnectionType = "claim_queue"
	ConnectionTypeUser      ConnectionType = "user"
)

// Message represents a WebSocket message.
type Message struct {
	Type       string      `json:"type"`
	ResourceID string      `json:"resource_id,omitempty"`
	Data       interface{} `json:"data"`
}

// ClientMessage represents a message from a WebSocket client.
type ClientMessage struct {
	Type       string `json:"type"`       // ping, subscribe, unsubscribe
	ResourceID string `json:"resource_id,omitempty"`
}
