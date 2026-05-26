package game

import (
	"fmt"
	"math"
	"strings"
	"time"

	"GoServer/internal/database"
)

type CommandHandler func(c *Client, h *Hub, data map[string]interface{})

var commands = map[string]CommandHandler{
	"chat_msg":       handleChat,
	"move":           handleMoveRequest,
	"portal_request": handlePortalMoveRequest,
}

func handleChat(c *Client, h *Hub, data map[string]interface{}) {
	rawText, _ := data["text"].(string)
	text := strings.TrimSpace(rawText)
	if text == "" {
		return
	}
	if len([]rune(text)) > 150 {
		text = string([]rune(text)[:150])
	}
	chatPacket := map[string]interface{}{
		"type":   "chat_msg",
		"sender": c.Character.Name,
		"text":   text,
	}
	h.RoomBroadcast <- RoomMessage{
		WorldID:    c.Character.WorldID,
		LocationID: c.Character.LocationID,
		Payload:    chatPacket,
	}
}

func handleMoveRequest(c *Client, h *Hub, data map[string]interface{}) {
	if h.IsPlayerMoving(c.Character.ID) {
		return
	}
	targetID, ok := data["target_id"].(string)
	if !ok {
		return
	}

	world, exists := Universe[c.Character.WorldID]

	targetNode, exists := world.Points[targetID]
	sourceNode := world.Points[c.Character.LocationID]

	if !exists || targetID == c.Character.LocationID {
		return
	}
	dx := float64(targetNode.X - sourceNode.X)
	dy := float64(targetNode.Y - sourceNode.Y)
	result := math.Ceil(math.Sqrt(dx*dx+dy*dy) / 10)

	duration := time.Duration(result) * time.Second
	charID := c.Character.ID
	worldID := c.Character.WorldID
	h.mu.Lock()
	h.movingPlayers[c.Character.ID] = &MoveData{
		DestinationID:      targetID,
		ArrivalTime:        time.Now().Add(duration),
		TargetWorldName:    world.Name,
		TargetLocationName: targetNode.Name,
	}
	h.mu.Unlock()

	h.Send(c, map[string]interface{}{
		"type": "move_starting",
		/*"target_name":   targetNode.Name,*/
		"world_name":    world.Name,
		"location_name": targetNode.Name,
		"duration":      duration.Seconds(),
	})

	go func() {
		time.Sleep(duration)

		h.mu.Lock()
		delete(h.movingPlayers, charID)
		h.mu.Unlock()

		// 1. Обновляем базу (это можно делать без мьютекса)
		_ = database.UpdateCharacterLocation(charID, targetID)

		// 2. Берем замок, чтобы безопасно обновить данные и собрать список
		h.mu.Lock() // Берем Lock, так как мы будем ИЗМЕНЯТЬ данные персонажа
		activeClient, online := h.Clients[charID]

		if !online {
			h.mu.Unlock()
			return
		}

		// Сначала официально "переставляем" игрока в новую комнату в памяти
		oldLockID := activeClient.Character.LocationID
		activeClient.Character.LocationID = targetID
		// ТЕПЕРЬ собираем соседей. Теперь игрок сам попадет в этот список!
		newNeighbors := h.getNeighbors(activeClient.Character.WorldID, targetID)

		// Получаем данные комнаты для списка порталов
		currentWorld := Universe[activeClient.Character.WorldID]
		currentNode := currentWorld.Points[targetID]
		h.mu.Unlock() // Все операции с данными закончены, отпускаем

		// 3. Отправляем пакет прибытия
		h.Send(activeClient, map[string]interface{}{
			"type":          "move_complete",
			"location_id":   targetID,
			"location_name": currentNode.Name,
			"players":       newNeighbors,
			"worlds":        currentNode.Worlds,
		})

		// 4. Оповещаем остальных
		h.BroadcastToRoomExcept(worldID, oldLockID, charID, map[string]interface{}{
			"type": "player_left",
			"player": map[string]interface{}{
				"id":   charID,
				"name": activeClient.Character.Name,
			},
		})
		h.BroadcastToRoomExcept(worldID, targetID, charID, map[string]interface{}{
			"type":   "player_joined",
			"player": activeClient.Character,
		})
	}()
}

func handlePortalMoveRequest(c *Client, h *Hub, data map[string]interface{}) {
	if h.IsPlayerMoving(c.Character.ID) {
		return
	}
	targetWorldID, ok := data["world_id"].(string)
	if !ok {
		return
	}

	targetWorld, exists := Universe[targetWorldID]
	if !exists {
		return
	}

	currentWorld := Universe[c.Character.WorldID]
	currentNode := currentWorld.Points[c.Character.LocationID]

	// Проверяем, есть ли портал в этот мир
	canTeleport := false
	for _, el := range currentNode.Worlds {
		if el.ID == targetWorldID {
			canTeleport = true
			break
		}
	}

	if !canTeleport {
		fmt.Printf("Игрок %s: попытка незаконной телепортации\n", c.Character.Name)
		return
	}

	const portalDuration = 10 * time.Second // Твои 200 секунд
	charID := c.Character.ID
	oldWorldID := c.Character.WorldID
	oldLocID := c.Character.LocationID

	h.mu.Lock()
	h.movingPlayers[charID] = &MoveData{
		DestinationID:      "portal",
		TargetWorldName:    targetWorld.Name,
		TargetLocationName: Universe[targetWorldID].Points["portal"].Name,
		ArrivalTime:        time.Now().Add(portalDuration),
	}
	h.mu.Unlock()

	// Сообщаем о начале долгого перехода
	h.Send(c, map[string]interface{}{
		"type":          "move_starting",
		"world_name":    targetWorld.Name,
		"location_name": Universe[targetWorldID].Points["portal"].Name,
		"duration":      int(portalDuration.Seconds()),
	})

	go func() {
		time.Sleep(portalDuration)

		h.mu.Lock()
		delete(h.movingPlayers, charID)
		h.mu.Unlock()

		// Обновляем БД
		_ = database.UpdateCharacterWorld(charID, targetWorldID, "portal")

		h.mu.RLock()
		activeClient, online := h.Clients[charID]
		h.mu.RUnlock()

		if online {
			activeClient.Character.WorldID = targetWorldID
			activeClient.Character.LocationID = "portal"

			// Полная синхронизация для прыгнувшего
			h.Send(activeClient, map[string]interface{}{
				"type":        "world_sync",
				"location_id": "portal",
				"world_id":    targetWorldID,
				"player":      activeClient.Character,
				"world":       Universe[targetWorldID],
				"players":     h.getNeighbors(targetWorldID, "portal"),
				"worlds":      Universe[targetWorldID].Points["portal"].Worlds,
			})

			// Оповещаем старый мир
			h.BroadcastToRoomExcept(oldWorldID, oldLocID, charID, map[string]interface{}{
				"type":   "player_left",
				"player": map[string]interface{}{"id": charID, "name": activeClient.Character.Name},
			})
			// Оповещаем новый мир
			h.BroadcastToRoomExcept(targetWorldID, "portal", charID, map[string]interface{}{
				"type":   "player_joined",
				"player": activeClient.Character,
			})
		}
	}()
}
