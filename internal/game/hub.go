package game

import (
	"fmt"
	"sync"
)

type Hub struct {
	mu         sync.RWMutex
	Clients    map[int64]*Client
	Register   chan *Client
	Unregister chan int64
	Broadcast  chan interface{}
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[int64]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan int64),
		Broadcast:  make(chan interface{}),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register: // Регистрация.
			h.mu.Lock()
			// var neighbors []map[string]interface{}
			h.Clients[client.Character.ID] = client
			neighbors := make([]map[string]interface{}, 0)
			for _, other := range h.Clients {
				if other.Character.LocationID == client.Character.LocationID {
					neighbors = append(neighbors, map[string]interface{}{
						"id":        other.Character.ID,
						"name":      other.Character.Name,
						"avatar_id": other.Character.AvatarID,
						"level":     other.Character.Level,
						"gender":    other.Character.Gender,
					})
				}
			}
			h.mu.Unlock()
			client.Send <- map[string]interface{}{ // Отправляем список тех кто уже был в комнате.
				"type":    "room_presence",
				"players": neighbors,
			}
			exeptID := client.Character.ID
			lockID := client.Character.LocationID
			h.BroadcastToRoomExcept(lockID, exeptID, map[string]interface{}{
				"type": "player_joined",
				"player": map[string]interface{}{
					"id":        client.Character.ID,
					"name":      client.Character.Name,
					"avatar_id": client.Character.AvatarID,
					"gender":    client.Character.Gender,
					"level":     client.Character.Level,
				},
			})
			fmt.Printf("Персонаж %s онлайн. \n", client.Character.Name)
		case id := <-h.Unregister:
			h.mu.Lock()
			if client, ok := h.Clients[id]; ok {
				fmt.Printf("Персонаж %s не в сети. \n", h.Clients[id].Character.Name)
				locID := client.Character.LocationID
				name := client.Character.Name
				delete(h.Clients, id)
				close(client.Send)
				h.mu.Unlock()
				h.BroadcastToRoom(locID, map[string]interface{}{
					"type": "player_left",
					"player": map[string]interface{}{
						"name": name,
					},
				})

			} else {
				h.mu.Unlock()
			}
		case message := <-h.Broadcast:
			h.BroadcastToAll(message)
		}
	}
}

func (h *Hub) BroadcastToRoom(locationID string, message interface{}) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, client := range h.Clients {
		if client.Character.LocationID == locationID {
			select {
			case client.Send <- message:
			default:
			}
		}
	}
}

func (h *Hub) BroadcastToRoomExcept(locationID string, exeptID int64, message interface{}) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, client := range h.Clients {
		if client.Character.LocationID == locationID && client.Character.ID != exeptID {
			select {
			case client.Send <- message:
			default:
			}
		}
	}
}

func (h *Hub) BroadcastToAll(message interface{}) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, client := range h.Clients {
		select {
		case client.Send <- message:
		default:
		}
	}
}
