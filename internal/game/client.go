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
}

func (c *Client) ReadPump() {
}
