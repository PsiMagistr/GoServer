package main

import (
	"GoServer/internal/auth"
	"GoServer/internal/database"
	"GoServer/internal/handlers"
	"fmt"
	"log"
	"net/http"
)

func main() {
	if err := database.InitDB(); err != nil {
		log.Fatal("Произошла ошибка ", err)
	}
	defer func() {
		err := database.DB.Close()
		if err != nil {
			log.Println("Ошибка при закрытии базы данных:", err)
		}
	}()
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)
	//http.Handle("/api/me/", auth.AuthMiddleware(http.HandlerFunc(handlers.MeHandler)))
	http.Handle("/api/me", auth.AuthMiddleware(http.HandlerFunc(handlers.MeHandler)))
	http.HandleFunc("/api/register", handlers.RegisterHandler)
	http.HandleFunc("/api/login", handlers.LoginHandler)
	port := ":8080"
	fmt.Println("Listening on " + port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Println("Ошибка при запуске сервера:", err)
	}
	fmt.Println("Сервер RPG запущен через IDE!")
}
