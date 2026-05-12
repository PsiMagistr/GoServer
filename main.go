package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"GoServer/internal/auth"
	"GoServer/internal/database"
	"GoServer/internal/game"
	"GoServer/internal/handlers"
)

func main() {
	if err := database.InitDB(); err != nil {
		log.Fatal("Произошла ошибка ", err)
	}
	gameHub := game.NewHub()
	go gameHub.Run()
	defer func() {
		err := database.DB.Close()
		if err != nil {
			log.Println("Ошибка при закрытии базы данных:", err)
		}
	}()
	fs := http.FileServer(http.Dir("static"))
	mux := http.NewServeMux()
	mux.Handle("/", fs)
	mux.Handle("/api/me", auth.AuthMiddleware(http.HandlerFunc(handlers.MeHandler)))
	mux.Handle("/api/character/me", auth.AuthMiddleware(http.HandlerFunc(handlers.CheckCharacterHandler)))
	mux.Handle("/api/character/create", auth.AuthMiddleware(http.HandlerFunc(handlers.CreateCharacterHandler)))
	mux.HandleFunc("/api/register", handlers.RegisterHandler)
	mux.HandleFunc("/api/login", handlers.LoginHandler)
	mux.HandleFunc("/api/refresh", handlers.RefreshHandler)
	mux.Handle("/ws", auth.AuthMiddleware(handlers.WSHandler(gameHub)))
	var MyServer Server = NewHTTPServer("0.0.0.0:8080", mux)
	go func() {
		err := MyServer.Run()
		if err != nil {
			fmt.Println("Ошибка при запуске сервера:", err)
		}
	}()
	fmt.Println("К полету готов. ")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Эта строка заблокирует main до тех пор, пока в канал не придет сигнал
	<-stop
	// 6. Graceful Shutdown (Плавная остановка)
	// Даем серверу 5 секунд на то, чтобы завершить текущие запросы
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := MyServer.Shutdown(ctx); err != nil {
		log.Printf("Ошибка при плавной остановке: %v", err)
	}
	// Здесь сработает твой defer database.DB.Close()
	fmt.Println("Программа завершена успешно.")
}
