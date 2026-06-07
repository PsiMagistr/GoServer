package game

import (
	"fmt"
	"math"
	"sync"
	"time"

	"GoServer/internal/database"
)

// Структура для комнатных сообщений
type RoomMessage struct {
	WorldID    string
	LocationID string
	Payload    interface{}
}

type MoveData struct {
	DestinationID      string
	ArrivalTime        time.Time
	TargetWorldName    string
	TargetLocationName string
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
	regenTicker := time.NewTicker(5 * time.Second)
	defer regenTicker.Stop()
	for {
		select {
		case client := <-h.Register: // Регистрация.
			h.handleRegister(client)
		case client := <-h.Unregister:
			h.handleUnregister(client)
		case globalMessage := <-h.Broadcast:
			h.BroadcastToAll(globalMessage)
		case roomMessage := <-h.RoomBroadcast:
			h.BroadcastToRoom(roomMessage.WorldID, roomMessage.LocationID, roomMessage.Payload)
		case <-regenTicker.C:
			h.handleRegeniration()
		}
	}
}

////////handlers

func (h *Hub) handleRegister(client *Client) {
	h.mu.Lock()
	// 1. Кикаем старую сессию (защита от мульти-вкладок)
	if oldClient, ok := h.Clients[client.Character.ID]; ok {
		oldClient.Conn.Close()
		fmt.Printf("Персонаж %s зашел из другого места, старая сессия закрыта.\n", client.Character.Name)
	}
	// Регистрация в карте онлайна
	h.Clients[client.Character.ID] = client
	// 2. Сбор данных для атомарного пакета
	moveInfo, isMoving := h.movingPlayers[client.Character.ID]
	neighbors := h.getNeighbors(client.Character.WorldID, client.Character.LocationID)
	currentWorld := Universe[client.Character.WorldID]
	currentNode := currentWorld.Points[client.Character.LocationID]
	h.mu.Unlock()

	// 3. ОТПРАВЛЯЕМ ЕДИНЫЙ ПАКЕТ СИНХРОНИЗАЦИИ
	// Теперь фронтенд получит всё: кто он, где он, кто рядом и какие порталы доступны
	var timeLeft int
	worldName := currentWorld.Name
	locationName := currentNode.Name
	if isMoving {
		timeLeft = int(math.Ceil(time.Until(moveInfo.ArrivalTime).Seconds()))
		worldName = moveInfo.TargetWorldName
		locationName = moveInfo.TargetLocationName
	}
	h.Send(client, map[string]interface{}{
		"type":          "world_sync",
		"is_moving":     isMoving,
		"player":        client.Character,   // Данные персонажа (HP, мана, статы)
		"world":         currentWorld,       // Данные мира (точки для канваса)
		"players":       neighbors,          // Список людей в комнате
		"worlds":        currentNode.Worlds, // Доступные переходы (порталы)
		"duration":      timeLeft,
		"world_id":      client.Character.WorldID,
		"location_id":   client.Character.LocationID,
		"world_name":    worldName,
		"location_name": locationName,
	})

	// 4. Оповещаем соседей (это всё еще отдельный пакет для ДРУГИХ игроков)
	h.BroadcastToRoomExcept(client.Character.WorldID, client.Character.LocationID, client.Character.ID, map[string]interface{}{
		"type": "player_joined",
		"player": map[string]interface{}{
			"id":        client.Character.ID,
			"name":      client.Character.Name,
			"avatar_id": client.Character.AvatarID,
			"level":     client.Character.Level,
			"gender":    client.Character.Gender,
		},
	})
	fmt.Printf("Персонаж %s онлайн.\n", client.Character.Name)
}

func (h *Hub) handleUnregister(client *Client) {
	h.mu.Lock()
	currentInMap, ok := h.Clients[client.Character.ID]
	if ok && currentInMap == client {
		type Params = struct {
			charID      int64
			charName    string
			charWorldID string
			charLocID   string
			charHP      float64
			charMana    float64
		}
		param := Params{
			charID:      client.Character.ID,
			charWorldID: client.Character.WorldID,
			charLocID:   client.Character.LocationID,
			charName:    client.Character.Name,
			charHP:      client.Character.HP,
			charMana:    client.Character.Mana,
		}
		fmt.Printf("Персонаж %s не в сети. \n", param.charName)
		go func(p Params) {
			err := database.UpdateCharacterHpMana(p.charID, p.charHP, p.charMana)
			if err != nil {
				fmt.Printf("ОШИБКА сохранения персонажа %s (ID %d): %v", p.charName, p.charID, err)
			}
		}(param)
		delete(h.Clients, client.Character.ID)
		close(client.Send)
		h.mu.Unlock()
		h.BroadcastToRoom(param.charWorldID, param.charLocID, map[string]interface{}{
			"type": "player_left",
			"player": map[string]interface{}{
				"id":   param.charID,
				"name": param.charName,
			},
		})
	} else {
		// Если это "призрак" старой сессии, просто отпускаем замок и ничего не делаем
		h.mu.Unlock()
		// log.Println("Игнорируем попытку удаления устаревшей сессии")
	}
}

func (h *Hub) BroadcastToRoom(worldID string, locationID string, message interface{}) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, client := range h.Clients {
		if client.Character.WorldID == worldID && client.Character.LocationID == locationID {
			select {
			case client.Send <- message:
			default:
				client.Conn.Close()
			}
		}
	}
}

func (h *Hub) BroadcastToRoomExcept(worldID string, locationID string, exeptID int64, message interface{}) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, client := range h.Clients {
		if client.Character.WorldID == worldID && client.Character.LocationID == locationID && client.Character.ID != exeptID {
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

func (h *Hub) getNeighbors(worldID string, locationID string) []map[string]interface{} {
	neighbors := make([]map[string]interface{}, 0)
	for _, other := range h.Clients {
		if other.Character.WorldID == worldID && other.Character.LocationID == locationID {
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

/*func (h *Hub) ResyncRoomPresence(c *Client) {
	h.mu.RLock()
	neighbors := h.getNeighbors(c.Character.WorldID, c.Character.LocationID)
	h.mu.RUnlock()
	currentWorld := Universe[c.Character.WorldID]
	h.Send(c, map[string]interface{}{
		"type":    "room_presence",
		"players": neighbors,
		"worlds":  currentWorld.Points[c.Character.LocationID].Worlds,
	})
}*/

func (h *Hub) IsPlayerMoving(charID int64) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, moving := h.movingPlayers[charID]
	return moving
}

// Для чата.
func (h *Hub) GetClientByName(name string) *Client {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, client := range h.Clients {
		if client.Character.Name == name {
			return client
		}
	}
	return nil
}

// Регенерация жизни и маны.
func (h *Hub) handleRegeniration() {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, client := range h.Clients {
		if h.IsPlayerMoving(client.Character.ID) {
			continue
		}
		hpChanged := client.AddHP(2)
		mpChanged := client.AddMana(5)
		if hpChanged || mpChanged {
			h.Send(client, map[string]interface{}{
				"type":   "player_update",
				"player": client.Character,
			})
		}
	}
}
