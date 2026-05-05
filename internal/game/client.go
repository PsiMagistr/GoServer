package game

import (
	"GoServer/internal/models"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn      *websocket.Conn
	Character *models.Character
	Send      chan interface{}
}

func (c *Client) WritePump() {
	defer c.Conn.Close()
	for {
		message, ok := <-c.Send
		if !ok {
			c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}
		c.Conn.WriteJSON(message)
	}
}

func (c *Client) ReadPump(h *Hub) {
	defer func() {
		h.Unregister(c.Character.ID)
		c.Conn.Close()
	}()
	for {
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
	}
}
