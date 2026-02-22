package sse

// SSEManager manages Server-Sent Events connections and notifications.
type SSEManager interface {
	Register(client *Client)
	Unregister(client *Client)
	SendToResource(connType ConnectionType, resourceID string, event Event)
	SendToUser(userID string, event Event)
	Start()
	Stop()
	GetConnectionCount() int
}
