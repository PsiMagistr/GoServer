package game

import (
	"time"

	"GoServer/internal/models"

	"github.com/gorilla/websocket"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

type Client struct {
	Conn      *websocket.Conn
	Character *models.Character
	Send      chan interface{}
}

func (c *Client) WritePump() { // Читаем сообщения от сервера и пишем в браузер.
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				_ = c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			err := c.Conn.WriteJSON(message)
			if err != nil {
				return
			}
		case <-ticker.C:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			err := c.Conn.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				return
			}

		}
	}
}

func (c *Client) ReadPump(h *Hub) { // Читаем сообщения от браузера и отправляем на сервер.
	c.Conn.SetReadLimit(512) // Максимальный размер сообщения от игрока (защита от спама)
	_ = c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		_ = c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	defer func() {
		h.Unregister <- c
		c.Conn.Close()
	}()
	for {
		var incoming map[string]interface{}
		err := c.Conn.ReadJSON(&incoming)
		if err != nil {
			break
		}
		msgType, ok := incoming["type"].(string)
		if !ok {
			continue
		}
		if msgType == "chat_msg" {
			text, _ := incoming["text"].(string)
			chatPacket := map[string]interface{}{
				"type":   msgType,
				"sender": c.Character.Name,
				"text":   text,
			}
			h.RoomBroadcast <- RoomMessage{
				LocationID: c.Character.LocationID,
				Payload:    chatPacket,
			}
		}
	}
}
