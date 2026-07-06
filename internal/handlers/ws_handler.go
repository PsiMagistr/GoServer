package handlers

import (
	"log"
	"net/http"

	"GoServer/internal/auth"
	"GoServer/internal/database"
	"GoServer/internal/game"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WSHandler(hub *game.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		val := r.Context().Value(auth.UserContextKey)
		if val == nil {
			log.Println("WS: Попытка входа без авторизации")
			return
		}
		claims := val.(*auth.Claims)
		char, err := database.GetCharacterByUserID(claims.UserID)
		if err != nil || char == nil {
			log.Printf("WS: У пользователя %d нет персонажа", claims.UserID)
			return
		}

		spells, err := database.GetCharacterSpells(char.ID)
		if err != nil {
			log.Println("Ошибка получения заклинаний БД. ", err)
			return
		}
		char.Spells = spells

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("WS Upgrade error:", err)
			return
		}

		client := &game.Client{
			Conn:      conn,
			Character: char,
			Send:      make(chan interface{}, 256),
		}
		hub.Register <- client
		go client.WritePump()
		go client.ReadPump(hub)
	}
}
