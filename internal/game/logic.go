package game

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"GoServer/internal/database"
	"GoServer/internal/models"
)

type CommandHandler func(c *Client, h *Hub, data map[string]interface{})

type CommitStatsRequest struct {
	Strength  int `json:"strength"`
	Agility   int `json:"agility"`
	Intuition int `json:"intuition"`
	Wisdom    int `json:"wisdom"`
	Charm     int `json:"charm"`
	Vitality  int `json:"vitality"`
}

var commands = map[string]CommandHandler{
	"chat_msg":       handleChat,
	"move":           handleMoveRequest,
	"portal_request": handlePortalMoveRequest,
	"private_chat":   handleWhisperRequest,
	"commit_stats":   handleStatsCommitRequest,
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
		"type":          "move_starting",
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

// Шепот для приват-чата.
func handleWhisperRequest(c *Client, h *Hub, data map[string]interface{}) {
	targetName, _ := data["target_name"].(string)
	text := data["text"].(string)

	if text == "" || targetName == "" {
		return
	}
	if targetName == c.Character.Name {
		h.Send(c, map[string]interface{}{
			"type": "sys_msg",
			"text": "Вы пытаетесь отправить сообщение самому себе!",
		})
		return
	}
	if len([]rune(text)) > 150 {
		text = string([]rune(text)[:150])
	}
	targetClient := h.GetClientByName(targetName)
	if targetClient == nil {
		h.Send(c, map[string]interface{}{
			"type": "sys_msg",
			"text": "Персонаж " + targetName + " не в сети.",
		})
		return
	}
	h.Send(targetClient, map[string]interface{}{
		"type": "whisper_received",
		"from": c.Character.Name,
		"text": text,
	})
	h.Send(c, map[string]interface{}{
		"type": "whisper_sent",
		"to":   targetName,
		"text": text,
	})
}

func handleStatsCommitRequest(c *Client, h *Hub, data map[string]interface{}) {
	statsData, ok := data["stats"]
	if !ok {
		log.Println("Ошибка нет ключа stats")
		return
	}
	var req CommitStatsRequest
	dataBytes, _ := json.Marshal(statsData)
	err := json.Unmarshal(dataBytes, &req)
	if err != nil {
		log.Println("Ошибка демаршалинга статов")
		return
	}
	h.mu.Lock()
	diffStr := req.Strength - c.Character.Strength
	diffAgi := req.Agility - c.Character.Agility
	diffInt := req.Intuition - c.Character.Intuition
	diffVit := req.Vitality - c.Character.Vitality
	diffCharm := req.Charm - c.Character.Charm
	diffWisdom := req.Wisdom - c.Character.Wisdom
	totalSpent := diffStr + diffAgi + diffInt + diffVit + diffCharm + diffWisdom
	if diffStr < 0 || diffAgi < 0 || diffInt < 0 || diffVit < 0 || diffCharm < 0 || diffWisdom < 0 {
		log.Printf("Игрок %s пытался уменьшить характеристики!", c.Character.Name)
		h.mu.Unlock()
		return
	}
	if totalSpent > c.Character.FreePoints {
		log.Printf("Игрок %s пытался потратить %d очков, имея %d",
			c.Character.Name, totalSpent, c.Character.FreePoints)
		h.mu.Unlock()
		return
	}
	updatedChar := *c.Character
	updatedChar.Strength = req.Strength
	updatedChar.Agility = req.Agility
	updatedChar.Intuition = req.Intuition
	updatedChar.Vitality = req.Vitality
	updatedChar.Charm = req.Charm
	updatedChar.Wisdom = req.Wisdom
	updatedChar.MaxHP = 150 + float64(updatedChar.Vitality*2)
	updatedChar.MaxMana = 100 + float64(updatedChar.Wisdom*2)
	updatedChar.FreePoints -= totalSpent
	h.mu.Unlock()
	go func(char models.Character) {
		err := database.UpdateCharacter(&char)
		h.mu.Lock()
		defer h.mu.Unlock()
		activeClient, online := h.Clients[char.ID]
		if online {
			if err != nil {
				log.Printf("Ошибка записи статов в БД: %v", err)
				h.Send(activeClient, map[string]interface{}{
					"type":  "error_msg",
					"error": err,
				})
				return
			}
			// Сверяем, что клиент тот же самый (реконнект)
			*activeClient.Character = char
			h.Send(activeClient, map[string]interface{}{
				"type":   "player_update",
				"player": activeClient.Character,
			})
		}
	}(updatedChar)
}
