package game

import (
	"fmt"
	"strings"
	"time"

	"GoServer/internal/database"
)

type CommandHandler func(c *Client, h *Hub, data map[string]interface{})

var commands = map[string]CommandHandler{
	"chat_msg": handleChat,
	"move":     handleMoveRequest,
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
		LocationID: c.Character.LocationID,
		Payload:    chatPacket,
	}
}

func handleMoveRequest(c *Client, h *Hub, data map[string]interface{}) {
	if c.IsMoving {
		return
	}
	targetID, ok := data["target_id"].(string)
	if !ok {
		return
	}

	world, exists := Universe[c.Character.WorldID]

	targetNode, exists := world.Points[targetID]

	if !exists || targetID == c.Character.LocationID {
		return
	}

	c.IsMoving = true
	moveDuration := 5

	h.Send(c, map[string]interface{}{
		"type":        "move_starting",
		"target_name": targetNode.Name,
		"duration":    moveDuration,
	})

	go func() {
		time.Sleep(time.Duration(moveDuration) * time.Second)
		h.mu.RLock()
		_, online := h.Clients[c.Character.ID]
		h.mu.RUnlock()
		if !online {
			return
		}
		oldLockID := c.Character.LocationID
		c.Character.LocationID = targetID
		c.IsMoving = false
		err := database.UpdateCharacterLocation(c.Character.ID, targetID)
		if err != nil {
			fmt.Printf("Ошибка сохранения локации: %v", err)
		}
		h.Send(c, map[string]interface{}{
			"type":        "move_complete",
			"location_id": targetID,
		})
		h.BroadcastToRoomExcept(oldLockID, c.Character.ID, map[string]interface{}{
			"type": "player_left",
			"player": map[string]interface{}{
				"id":   c.Character.ID,
				"name": c.Character.Name,
			},
		})
		h.BroadcastToRoomExcept(targetID, c.Character.ID, map[string]interface{}{
			"type":   "player_joined",
			"player": c.Character,
		})
		h.ResyncRoomPresents(c)
	}()
}
