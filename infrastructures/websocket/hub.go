package websocket

import (
	"encoding/json"
	"log"
	"sync"
)

// Hub manages all WebSocket client connections.
type Hub struct {
	connections     map[ConnectionType]map[string]map[string]*Client // type → resourceID → clientID → Client
	userConnections map[string]map[string]*Client                    // userID → clientID → Client

	register   chan *Client
	unregister chan *Client
	notify     chan wsNotification
	broadcast  chan wsBroadcast
	stop       chan struct{}
	mu         sync.RWMutex
}

type wsNotification struct {
	connType   ConnectionType
	resourceID string
	userID     string
	message    Message
}

type wsBroadcast struct {
	connType ConnectionType
	message  Message
}

// NewHub creates a new WebSocket hub.
func NewHub() *Hub {
	return &Hub{
		connections:     make(map[ConnectionType]map[string]map[string]*Client),
		userConnections: make(map[string]map[string]*Client),
		register:        make(chan *Client, 256),
		unregister:      make(chan *Client, 256),
		notify:          make(chan wsNotification, 1024),
		broadcast:       make(chan wsBroadcast, 256),
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
			h.sendNotification(n)
		case b := <-h.broadcast:
			h.broadcastMessage(b)
		case <-h.stop:
			h.closeAll()
			return
		}
	}
}

func (h *Hub) addClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.connections[client.ConnectionType]; !ok {
		h.connections[client.ConnectionType] = make(map[string]map[string]*Client)
	}
	if _, ok := h.connections[client.ConnectionType][client.ResourceID]; !ok {
		h.connections[client.ConnectionType][client.ResourceID] = make(map[string]*Client)
	}
	h.connections[client.ConnectionType][client.ResourceID][client.ID] = client

	if _, ok := h.userConnections[client.UserID]; !ok {
		h.userConnections[client.UserID] = make(map[string]*Client)
	}
	h.userConnections[client.UserID][client.ID] = client

	log.Printf("WS client %s registered (type: %s, resource: %s)", client.ID, client.ConnectionType, client.ResourceID)
}

func (h *Hub) removeClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if typeConns, ok := h.connections[client.ConnectionType]; ok {
		if resourceConns, ok := typeConns[client.ResourceID]; ok {
			if _, ok := resourceConns[client.ID]; ok {
				delete(resourceConns, client.ID)
				close(client.send)
			}
			if len(resourceConns) == 0 {
				delete(typeConns, client.ResourceID)
			}
		}
	}

	if userConns, ok := h.userConnections[client.UserID]; ok {
		delete(userConns, client.ID)
		if len(userConns) == 0 {
			delete(h.userConnections, client.UserID)
		}
	}
}

func (h *Hub) sendNotification(n wsNotification) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	data, err := json.Marshal(n.message)
	if err != nil {
		return
	}

	if n.userID != "" {
		if userConns, ok := h.userConnections[n.userID]; ok {
			for _, client := range userConns {
				select {
				case client.send <- data:
				default:
				}
			}
		}
		return
	}

	if typeConns, ok := h.connections[n.connType]; ok {
		if resourceConns, ok := typeConns[n.resourceID]; ok {
			for _, client := range resourceConns {
				select {
				case client.send <- data:
				default:
				}
			}
		}
	}
}

func (h *Hub) broadcastMessage(b wsBroadcast) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	data, err := json.Marshal(b.message)
	if err != nil {
		return
	}

	if typeConns, ok := h.connections[b.connType]; ok {
		for _, resourceConns := range typeConns {
			for _, client := range resourceConns {
				select {
				case client.send <- data:
				default:
				}
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
				close(client.send)
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
