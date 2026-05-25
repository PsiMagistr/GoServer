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
	if h.isMoving(c.Character.ID) {
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
		DestinationID: targetID,
		ArrivalTime:   time.Now().Add(duration),
		TargetName:    targetNode.Name,
	}
	h.mu.Unlock()

	h.Send(c, map[string]interface{}{
		"type":        "move_starting",
		"target_name": targetNode.Name,
		"duration":    duration.Seconds(),
	})

	go func() {
		time.Sleep(duration)
		h.mu.Lock()
		delete(h.movingPlayers, charID)
		h.mu.Unlock()
		h.mu.RLock()
		activeClient, online := h.Clients[charID]
		h.mu.RUnlock()
		err := database.UpdateCharacterLocation(charID, targetID)
		if err != nil {
			fmt.Printf("Ошибка сохранения локации: %v", err)
		}
		if !online {
			return
		}
		oldLockID := activeClient.Character.LocationID
		activeClient.Character.LocationID = targetID
		h.Send(activeClient, map[string]interface{}{
			"type":        "move_complete",
			"location_id": targetID,
		})
		/*h.Send(activeClient, map[string]interface{}{
			"type":        "world_sync",
			"location_id": activeClient.Character.LocationID,
			"world_id":    activeClient.Character.WorldID,
			"player":      activeClient.Character,                   // Данные персонажа
			"world":       Universe[activeClient.Character.WorldID], // Данные карты
			"players":     h.getNeighbors(activeClient.Character.WorldID, activeClient.Character.LocationID),
		})*/

		h.BroadcastToRoomExcept(worldID, oldLockID, c.Character.ID, map[string]interface{}{
			"type": "player_left",
			"player": map[string]interface{}{
				"id":   activeClient.Character.ID,
				"name": activeClient.Character.Name,
			},
		})
		h.BroadcastToRoomExcept(worldID, targetID, activeClient.Character.ID, map[string]interface{}{
			"type":   "player_joined",
			"player": activeClient.Character,
		})
		h.ResyncRoomPresence(activeClient)
	}()
}

// Путешествуем по мирам.
func handlePortalMoveRequest(c *Client, h *Hub, data map[string]interface{}) {
	if h.isMoving(c.Character.ID) {
		fmt.Println("Игрок в движении")
		return
	}
	targetWorldID, ok := data["world_id"].(string)
	if !ok {
		fmt.Println("Переданные данные неверны")
		return
	}
	targetWorld, exists := Universe[targetWorldID]

	if !exists {
		fmt.Println("Данные о мире неверны")
		return
	}
	currentWorld := Universe[c.Character.WorldID]
	currentNode := currentWorld.Points[c.Character.LocationID]
	canTeleport := false
	for _, el := range currentNode.Worlds {
		if el.ID == targetWorldID {
			canTeleport = true
		}
	}
	if !canTeleport {
		fmt.Printf("Игрок %s пытался войти в портал незаконно\n", c.Character.Name)
		return
	}
	const portalDuration = 5 * time.Second
	charID := c.Character.ID
	h.mu.Lock()
	h.movingPlayers[charID] = &MoveData{
		DestinationID: "portal",
		TargetName:    targetWorld.Name,
		ArrivalTime:   time.Now().Add(portalDuration),
	}
	h.mu.Unlock()
	h.Send(c, map[string]interface{}{
		"type":        "move_starting",
		"target_name": targetWorld.Name,
		"duration":    int(portalDuration.Seconds()),
	})
	go func() {
		time.Sleep(portalDuration)
		h.mu.Lock()
		delete(h.movingPlayers, charID)
		h.mu.Unlock()
		_ = database.UpdateCharacterWorld(charID, targetWorldID, "portal")
		h.mu.RLock()
		activeClient, online := h.Clients[charID]
		h.mu.RUnlock()
		if online {
			oldLocID := activeClient.Character.LocationID
			activeClient.Character.WorldID = targetWorldID
			activeClient.Character.LocationID = "portal"
			fmt.Println("dsdsdsdsds ", targetWorld)
			newWorld := Universe[targetWorldID]
			arrivalNode := newWorld.Points["portal"]
			h.Send(activeClient, map[string]interface{}{
				"type":        "world_sync",
				"location_id": "portal",
				"world_id":    targetWorldID,
				"player":      activeClient.Character,  // Данные персонажа
				"world":       Universe[targetWorldID], // Данные карты
				"players":     h.getNeighbors(targetWorldID, activeClient.Character.LocationID),
				"worlds":      arrivalNode.Worlds,
			})
			h.BroadcastToRoomExcept(currentWorld.ID, oldLocID, charID, map[string]interface{}{
				"type":   "player_left",
				"player": map[string]interface{}{"id": charID, "name": activeClient.Character.Name},
			})
			h.BroadcastToRoomExcept(targetWorldID, "portal", charID, map[string]interface{}{
				"type":   "player_joined",
				"player": activeClient.Character,
			})

		}
	}()
}
