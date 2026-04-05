package hub

import (
	"context"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn *websocket.Conn
	Send chan interface{}
}

type InMemoryHub struct {
	clients map[string][]*Client
	mu      sync.RWMutex
}

func NewInMemoryHub() *InMemoryHub {
	return &InMemoryHub{
		clients: make(map[string][]*Client),
	}
}

func (h *InMemoryHub) Send(ctx context.Context, userID string, payload interface{}) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if clients, ok := h.clients[userID]; ok {
		for _, client := range clients {
			select {
			case client.Send <- payload:
			default:
			}
		}
	}

	return nil
}

func (h *InMemoryHub) Register(userID string, conn *websocket.Conn) {
	client := &Client{
		Conn: conn,
		Send: make(chan interface{}, 256),
	}

	h.mu.Lock()
	h.clients[userID] = append(h.clients[userID], client)
	h.mu.Unlock()

	go h.writePump(client)
	go h.readPump(client)
}

func (h *InMemoryHub) Close() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	for _, clients := range h.clients {
		for _, c := range clients {
			close(c.Send)
			if err := c.Conn.Close(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (h *InMemoryHub) unregister(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for userID, clients := range h.clients {
		for i, c := range clients {
			if c == client {
				h.clients[userID] = append(clients[:i], clients[i+1:]...)
				break
			}
		}
	}
}

func (h *InMemoryHub) readPump(c *Client) {
	defer func() {
		h.unregister(c)
		c.Conn.Close()
	}()

	for {
		if _, _, err := c.Conn.ReadMessage(); err != nil {
			break
		}
	}
}

func (h *InMemoryHub) writePump(c *Client) {
	defer func() {
		h.unregister(c)
		c.Conn.Close()
	}()

	for msg := range c.Send {
		if err := c.Conn.WriteJSON(msg); err != nil {
			break
		}
	}
}
