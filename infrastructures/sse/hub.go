package sse

import (
	"encoding/json"
	"log"
	"sync"
)

// Hub manages all SSE client connections.
type Hub struct {
	// Connections organized by type and resource ID
	connections map[ConnectionType]map[string]map[string]*Client // type → resourceID → clientID → Client
	// User connections for user-targeted events
	userConnections map[string]map[string]*Client // userID → clientID → Client

	register    chan *Client
	unregister  chan *Client
	notify      chan notification
	stop        chan struct{}
	mu          sync.RWMutex
}

type notification struct {
	connType   ConnectionType
	resourceID string
	userID     string
	event      Event
}

// NewHub creates a new SSE hub.
func NewHub() *Hub {
	return &Hub{
		connections:     make(map[ConnectionType]map[string]map[string]*Client),
		userConnections: make(map[string]map[string]*Client),
		register:        make(chan *Client, 256),
		unregister:      make(chan *Client, 256),
		notify:          make(chan notification, 1024),
		stop:            make(chan struct{}),
	}
}

// Run starts the hub event loop.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.addClient(client)
		case client := <-h.unregister:
			h.removeClient(client)
		case n := <-h.notify:
			h.broadcast(n)
		case <-h.stop:
			h.closeAll()
			return
		}
	}
}

func (h *Hub) addClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Add to resource connections
	if _, ok := h.connections[client.ConnectionType]; !ok {
		h.connections[client.ConnectionType] = make(map[string]map[string]*Client)
	}
	if _, ok := h.connections[client.ConnectionType][client.ResourceID]; !ok {
		h.connections[client.ConnectionType][client.ResourceID] = make(map[string]*Client)
	}
	h.connections[client.ConnectionType][client.ResourceID][client.ID] = client

	// Add to user connections
	if _, ok := h.userConnections[client.UserID]; !ok {
		h.userConnections[client.UserID] = make(map[string]*Client)
	}
	h.userConnections[client.UserID][client.ID] = client

	log.Printf("SSE client %s registered (type: %s, resource: %s, user: %s)",
		client.ID, client.ConnectionType, client.ResourceID, client.UserID)
}

func (h *Hub) removeClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Remove from resource connections
	if typeConns, ok := h.connections[client.ConnectionType]; ok {
		if resourceConns, ok := typeConns[client.ResourceID]; ok {
			delete(resourceConns, client.ID)
			if len(resourceConns) == 0 {
				delete(typeConns, client.ResourceID)
			}
		}
	}

	// Remove from user connections
	if userConns, ok := h.userConnections[client.UserID]; ok {
		delete(userConns, client.ID)
		if len(userConns) == 0 {
			delete(h.userConnections, client.UserID)
		}
	}

	log.Printf("SSE client %s unregistered", client.ID)
}

func (h *Hub) broadcast(n notification) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	data, err := json.Marshal(n.event)
	if err != nil {
		log.Printf("Failed to marshal SSE event: %v", err)
		return
	}

	// Send to user if specified
	if n.userID != "" {
		if userConns, ok := h.userConnections[n.userID]; ok {
			for _, client := range userConns {
				client.Send(data)
			}
		}
		return
	}

	// Send to resource subscribers
	if typeConns, ok := h.connections[n.connType]; ok {
		if resourceConns, ok := typeConns[n.resourceID]; ok {
			for _, client := range resourceConns {
				client.Send(data)
			}
		}
	}
}

func (h *Hub) closeAll() {
	h.mu.Lock()
	defer h.mu.Unlock()

	for _, typeConns := range h.connections {
		for _, resourceConns := range typeConns {
			for _, client := range resourceConns {
				client.Close()
			}
		}
	}
	h.connections = make(map[ConnectionType]map[string]map[string]*Client)
	h.userConnections = make(map[string]map[string]*Client)
}

func (h *Hub) GetConnectionCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	count := 0
	for _, typeConns := range h.connections {
		for _, resourceConns := range typeConns {
			count += len(resourceConns)
		}
	}
	return count
}
