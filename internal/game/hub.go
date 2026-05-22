package game

import (
	"fmt"
	"math"
	"sync"
	"time"
)

// Структура для комнатных сообщений
type RoomMessage struct {
	LocationID string
	Payload    interface{}
}

type MoveData struct {
	DestinationID string
	ArrivalTime   time.Time
	TargetName    string
}

type Hub struct {
	mu            sync.RWMutex
	Clients       map[int64]*Client
	movingPlayers map[int64]*MoveData
	Register      chan *Client
	Unregister    chan *Client
	Broadcast     chan interface{}
	RoomBroadcast chan RoomMessage
}

func NewHub() *Hub {
	return &Hub{
		Clients:       make(map[int64]*Client),
		movingPlayers: make(map[int64]*MoveData),
		Register:      make(chan *Client, 64),
		Unregister:    make(chan *Client, 64),
		Broadcast:     make(chan interface{}, 256),
		RoomBroadcast: make(chan RoomMessage, 256),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register: // Регистрация.
			h.handleRegister(client)
		case client := <-h.Unregister:
			h.handleUnregister(client)
		case globalMessage := <-h.Broadcast:
			h.BroadcastToAll(globalMessage)
		case roomMessage := <-h.RoomBroadcast:
			h.BroadcastToRoom(roomMessage.LocationID, roomMessage.Payload)
		}
	}
}

////////handlers

func (h *Hub) handleRegister(client *Client) {
	h.mu.Lock()
	if oldClient, ok := h.Clients[client.Character.ID]; ok {
		oldClient.Conn.Close()
		fmt.Printf("Персонаж %s зашел из другого места, старая сессия закрыта.\n", client.Character.Name)
	}
	h.Clients[client.Character.ID] = client
	moveInfo, isMoving := h.movingPlayers[client.Character.ID]
	neighbors := h.getNeighbors(client.Character.LocationID)
	h.mu.Unlock()
	currentWorld := Universe[client.Character.WorldID]
	h.Send(client, map[string]interface{}{
		"type":   "self_load",
		"player": client.Character,
		"world":  currentWorld,
	})
	h.Send(client, map[string]interface{}{
		"type":    "room_presence",
		"players": neighbors,
	})
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
	if isMoving {
		secondsLeft := time.Until(moveInfo.ArrivalTime).Seconds()
		timeLeft := int(math.Ceil(secondsLeft))
		if timeLeft > 0 {
			h.Send(client, map[string]interface{}{
				"type":        "move_starting",
				"target_name": moveInfo.TargetName,
				"duration":    timeLeft,
			})
		}
	}
}

func (h *Hub) handleUnregister(client *Client) {
	h.mu.Lock()
	currentInMap, ok := h.Clients[client.Character.ID]
	if ok && currentInMap == client {
		locID := client.Character.LocationID
		name := client.Character.Name
		fmt.Printf("Персонаж %s не в сети. \n", name)
		delete(h.Clients, client.Character.ID)
		close(client.Send)
		h.mu.Unlock()
		h.BroadcastToRoom(locID, map[string]interface{}{
			"type": "player_left",
			"player": map[string]interface{}{
				"id":   client.Character.ID,
				"name": name,
			},
		})
	} else {
		// Если это "призрак" старой сессии, просто отпускаем замок и ничего не делаем
		h.mu.Unlock()
		// log.Println("Игнорируем попытку удаления устаревшей сессии")
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
				client.Conn.Close()
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
				client.Conn.Close()
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
			client.Conn.Close()
		}
	}
}

func (h *Hub) BroadcastPrivateMessage(charID int64, message interface{}) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	client, ok := h.Clients[charID]
	if !ok {
		return
	}
	select {
	case client.Send <- message:
	default:
		client.Conn.Close()
	}
}

func (h *Hub) Send(client *Client, message interface{}) {
	select {
	case client.Send <- message:
	default:
		client.Conn.Close()
	}
}

func (h *Hub) getNeighbors(locationID string) []map[string]interface{} {
	neighbors := make([]map[string]interface{}, 0)
	for _, other := range h.Clients {
		if other.Character.LocationID == locationID {
			neighbors = append(neighbors, map[string]interface{}{
				"id":        other.Character.ID,
				"name":      other.Character.Name,
				"avatar_id": other.Character.AvatarID,
				"level":     other.Character.Level,
				"gender":    other.Character.Gender,
			})
		}
	}
	return neighbors
}

func (h *Hub) ResyncRoomPresence(c *Client) {
	h.mu.RLock()
	neighbors := h.getNeighbors(c.Character.LocationID)
	h.mu.RUnlock()
	h.Send(c, map[string]interface{}{
		"type":    "room_presence",
		"players": neighbors,
	})
}

func (h *Hub) isMoving(charID int64) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, moving := h.movingPlayers[charID]
	return moving
}
