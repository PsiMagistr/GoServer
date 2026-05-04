package game

import (
	"log"
	"sync"
)

type Hub struct {
	mu      sync.RWMutex
	Clients map[int]*Client
}

func NewHub() *Hub {
	return &Hub{
		Clients: make(map[int]*Client),
	}
}

func (h *Hub) Register(c *Client) {
	h.mu.Lock()
	h.Clients[c.Character.ID] = c
	h.mu.Unlock()
	log.Printf("Персонаж %s вошел в онлайн. Всего %d персонажей.", c.Character.Name, len(h.Clients))
}

func (h *Hub) Unregister(charID int) {
	h.mu.Lock()
	delete(h.Clients, charID)
	h.mu.Unlock()
}

func (h *Hub) Broadcast(message interface{}) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, client := range h.Clients {
		client.Send <- message
	}
}
