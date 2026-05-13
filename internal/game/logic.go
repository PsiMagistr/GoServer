package game

import (
	"strings"
)

type CommandHandler func(c *Client, h *Hub, data map[string]interface{})

var commands = map[string]CommandHandler{
	"chat_msg": handleChat,
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
