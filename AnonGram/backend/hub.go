package main

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"sync"
)

type Client struct {
	conn   *websocket.Conn
	chatID string
	send   chan []byte
}
type Hub struct {
	mu      sync.RWMutex
	clients map[*Client]bool
}

func NewHub() *Hub { return &Hub{clients: map[*Client]bool{}} }
func (h *Hub) Broadcast(id string, b []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for c := range h.clients {
		if id == "" || c.chatID == id {
			select {
			case c.send <- b:
			default:
			}
		}
	}
}
func (h *Hub) Run(ctx context.Context, r *redis.Client) { <-ctx.Done() }
