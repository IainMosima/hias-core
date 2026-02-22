package sse

import "log"

type sseManagerImpl struct {
	hub *Hub
}

// NewSSEManager creates a new SSE manager.
func NewSSEManager() SSEManager {
	return &sseManagerImpl{
		hub: NewHub(),
	}
}

func (m *sseManagerImpl) Register(client *Client) {
	m.hub.register <- client
}

func (m *sseManagerImpl) Unregister(client *Client) {
	m.hub.unregister <- client
}

func (m *sseManagerImpl) SendToResource(connType ConnectionType, resourceID string, event Event) {
	m.hub.notify <- notification{
		connType:   connType,
		resourceID: resourceID,
		event:      event,
	}
}

func (m *sseManagerImpl) SendToUser(userID string, event Event) {
	m.hub.notify <- notification{
		userID: userID,
		event:  event,
	}
}

func (m *sseManagerImpl) Start() {
	go m.hub.Run()
	log.Println("SSE manager started")
}

func (m *sseManagerImpl) Stop() {
	close(m.hub.stop)
	log.Println("SSE manager stopped")
}

func (m *sseManagerImpl) GetConnectionCount() int {
	return m.hub.GetConnectionCount()
}
