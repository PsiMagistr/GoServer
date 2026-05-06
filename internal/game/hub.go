package game

import (
	"fmt"
	"sync"
)

type Hub struct {
	mu         sync.RWMutex
	Clients    map[int]*Client
	Register   chan *Client
	Unregister chan int
	Broadcast  chan interface{}
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[int]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan int),
		Broadcast:  make(chan interface{}),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.Clients[client.Character.ID] = client
			h.mu.Unlock()
			fmt.Printf("Персонаж %s онлайн. \n", client.Character.Name)
		case id := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Clients[id]; ok {
				fmt.Printf("Персонаж %s не в сети. \n", h.Clients[id].Character.Name)
				close(h.Clients[id].Send)
				delete(h.Clients, id)
			}
			h.mu.Unlock()
		case message := <-h.Broadcast:
			h.mu.RLock()
			for _, client := range h.Clients {
				select {
				case client.Send <- message:
				default:
				}
			}
			h.mu.RUnlock()
		}
	}
}
