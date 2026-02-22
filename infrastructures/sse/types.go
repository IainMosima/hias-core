package sse

// ConnectionType represents SSE connection categories.
type ConnectionType string

const (
	ConnectionTypeClaim    ConnectionType = "claim"
	ConnectionTypePayment  ConnectionType = "payment"
	ConnectionTypeDashboard ConnectionType = "dashboard"
	ConnectionTypeUser     ConnectionType = "user"
)

// Event represents an SSE event to be sent to clients.
type Event struct {
	Type       string      `json:"type"`
	ResourceID string      `json:"resource_id"`
	Data       interface{} `json:"data"`
}

// ClientRegistration holds info for registering an SSE client.
type ClientRegistration struct {
	Client         *Client
	ConnectionType ConnectionType
	ResourceID     string
}
