package sse

import (
	"fmt"
	"net/http"
	"time"
)

// Client represents an SSE client connection.
type Client struct {
	ID             string
	UserID         string
	ConnectionType ConnectionType
	ResourceID     string
	send           chan []byte
	done           chan struct{}
	writer         http.ResponseWriter
	flusher        http.Flusher
	createdAt      time.Time
}

// NewClient creates a new SSE client.
func NewClient(id, userID string, connType ConnectionType, resourceID string, w http.ResponseWriter) (*Client, error) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		return nil, fmt.Errorf("streaming not supported")
	}

	return &Client{
		ID:             id,
		UserID:         userID,
		ConnectionType: connType,
		ResourceID:     resourceID,
		send:           make(chan []byte, 256),
		done:           make(chan struct{}),
		writer:         w,
		flusher:        flusher,
		createdAt:      time.Now(),
	}, nil
}

// Listen starts the SSE event loop for this client.
func (c *Client) Listen() {
	// Set SSE headers
	c.writer.Header().Set("Content-Type", "text/event-stream")
	c.writer.Header().Set("Cache-Control", "no-cache")
	c.writer.Header().Set("Connection", "keep-alive")
	c.writer.Header().Set("X-Accel-Buffering", "no")

	// Send initial connection event
	fmt.Fprintf(c.writer, "data: {\"type\":\"connection_established\",\"client_id\":\"%s\"}\n\n", c.ID)
	c.flusher.Flush()

	heartbeat := time.NewTicker(30 * time.Second)
	defer heartbeat.Stop()

	for {
		select {
		case msg, ok := <-c.send:
			if !ok {
				return
			}
			fmt.Fprintf(c.writer, "data: %s\n\n", msg)
			c.flusher.Flush()
		case <-heartbeat.C:
			fmt.Fprint(c.writer, ": heartbeat\n\n")
			c.flusher.Flush()
		case <-c.done:
			return
		}
	}
}

// Send sends data to the client's channel.
func (c *Client) Send(data []byte) {
	select {
	case c.send <- data:
	default:
		// Channel full, drop message
	}
}

// Close closes the client connection.
func (c *Client) Close() {
	close(c.done)
	close(c.send)
}
