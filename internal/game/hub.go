package game

import (
	"fmt"
	"math"
	"sync"
	"time"

	"GoServer/internal/database"
	"GoServer/internal/models"
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

type BattleChallenge struct { // Заявка на бой
	SenderID   int64     `json:"sender_id"`
	SenderName string    `json:"sender_name"`
	TargetID   int64     `json:"target_id"`
	ExpiresAt  time.Time `json:"expires_at"`
	TimeLeft   int       `json:"time_left"`
}

type Hub struct {
	mu             sync.RWMutex
	Clients        map[int64]*Client
	movingPlayers  map[int64]*MoveData
	challenges     map[int64]map[int64]*BattleChallenge
	activeBattles  map[int64]*Battle
	playerToBattle map[int64]int64
	Register       chan *Client
	Unregister     chan *Client
	Broadcast      chan interface{}
	RoomBroadcast  chan RoomMessage
}

func NewHub() *Hub {
	return &Hub{
		Clients:        make(map[int64]*Client),
		movingPlayers:  make(map[int64]*MoveData),
		challenges:     make(map[int64]map[int64]*BattleChallenge),
		activeBattles:  make(map[int64]*Battle),
		playerToBattle: make(map[int64]int64),
		Register:       make(chan *Client, 64),
		Unregister:     make(chan *Client, 64),
		Broadcast:      make(chan interface{}, 256),
		RoomBroadcast:  make(chan RoomMessage, 256),
	}
}

func (h *Hub) Run() {
	regenTicker := time.NewTicker(5 * time.Second)
	cleanupTicker := time.NewTicker(1 * time.Minute)
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
		case <-cleanupTicker.C:
			h.cleanupChallenges()
		}
	}
}

/*Подготовка пакета syncWorld*/

func (h *Hub) prepareSyncState(client *Client) map[string]interface{} {
	neighbors := h.getNeighbors(client.Character.WorldID, client.Character.LocationID)
	currentWorld := Universe[client.Character.WorldID]
	currentNode := currentWorld.Points[client.Character.LocationID]
	challenges := h.GetChallenges(client.Character.ID)
	var timeLeft int
	worldName := currentWorld.Name
	locationName := currentNode.Name
	moveInfo, isMoving := h.movingPlayers[client.Character.ID]
	if isMoving {
		client.Character.State = models.StatusMoving // Синхронизируем стейт
		secondsLeft := time.Until(moveInfo.ArrivalTime).Seconds()
		timeLeft = int(math.Ceil(secondsLeft))
		worldName = moveInfo.TargetWorldName
		locationName = moveInfo.TargetLocationName
	}
	return map[string]interface{}{
		"type":          "world_sync",
		"player":        client.Character,
		"world":         currentWorld,
		"players":       neighbors,
		"challenges":    challenges,
		"worlds":        currentNode.Worlds,
		"is_moving":     isMoving,
		"duration":      timeLeft,
		"world_id":      client.Character.WorldID, // Например "main_city"
		"location_id":   client.Character.LocationID,
		"world_name":    worldName,
		"location_name": locationName,
		//"battle_info":   battleInfo,
	}
}

func (h *Hub) handleRegister(client *Client) {
	h.mu.Lock()
	// 1. Безопасность
	if oldClient, ok := h.Clients[client.Character.ID]; ok {
		oldClient.Conn.Close()
	}
	// 2. Регистрация
	h.Clients[client.Character.ID] = client
	// 3. Подготовка данных (Snapshot)
	syncData := h.prepareSyncState(client)
	h.mu.Unlock()
	// 4. Отправка (уже без блокировки всего сервера!)
	h.Send(client, syncData)
	// 5. Уведомление окружающих
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
		if h.GetFullStatus(client.Character.ID) != models.StatusFree {
			// fmt.Println("Вы не можете восстанавливаться")
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

func (h *Hub) SystemMsg(c *Client, text string) {
	h.Send(c, map[string]interface{}{
		"type": "sys_msg",
		"text": text,
	})
}

func (h *Hub) GetFullStatus(charID int64) models.PlayerStatus {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, moving := h.movingPlayers[charID]
	if moving {
		return models.StatusMoving
	}
	client, ok := h.Clients[charID]
	if ok {
		return client.Character.State
	}
	return models.StatusFree
}

func (h *Hub) GetActiveClient(charID int64) (*Client, bool) { // Не использовать в местах где мьютекс уже взят.
	h.mu.RLock()
	defer h.mu.RUnlock()
	client, ok := h.Clients[charID]
	return client, ok
}

func (h *Hub) GetChallenges(RecipientID int64) []*BattleChallenge {
	var myChallenges []*BattleChallenge
	pending, exists := h.challenges[RecipientID]
	if exists {
		for _, challenge := range pending {
			duration := time.Until(challenge.ExpiresAt)
			if duration > 0 {
				challenge.TimeLeft = int(math.Ceil(duration.Seconds()))
				myChallenges = append(myChallenges, challenge)
			}
		}
	}
	return myChallenges
}

func (h *Hub) cleanupChallenges() {
	h.mu.Lock()
	defer h.mu.Unlock()
	now := time.Now()
	count := 0
	for recipientID, invites := range h.challenges {
		for senderID, challenge := range invites {
			if now.After(challenge.ExpiresAt) {
				delete(invites, senderID)
				count++
			}
			if len(invites) == 0 {
				delete(h.challenges, recipientID)
			}
		}
	}
	if count > 0 {
		fmt.Printf("[CLEANUP] Удалено просроченных заявок: %d\n", count)
	}
}
