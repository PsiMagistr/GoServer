package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"GoServer/internal/auth"
	"GoServer/internal/config"
	"GoServer/internal/database"
	"GoServer/internal/game"
	"GoServer/internal/handlers"
	"GoServer/internal/middleware"

	"golang.org/x/time/rate"
)

func main() {
	err := config.LoadFile("config.json")
	if err != nil {
		log.Fatal("Критическая ошибка: конфиг не найден!", err)
	}
	err = game.LoadUniverse("data/universe.json")
	if err != nil {
		log.Fatal("Критическая ошибка загрузки мира:", err)
	}
	if err = database.InitDB(); err != nil {
		log.Fatal("Произошла ошибка ", err)
	}
	err = database.InitSpells(database.DB)
	if err != nil {
		log.Fatal("Ошибка загрузки заклинаний:", err)
	}
	addr := fmt.Sprintf("%s:%s", config.Get().Server.IP, config.Get().Server.PORT)
	gameHub := game.NewHub()
	go gameHub.Run()
	defer func() {
		err := database.DB.Close()
		if err != nil {
			log.Println("Ошибка при закрытии базы данных:", err)
		}
	}()
	fs := http.FileServer(http.Dir("static"))
	apiMux := http.NewServeMux()
	mainMux := http.NewServeMux()
	mainMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Если это корень сайта, генерируем HTML с base_path
		if r.URL.Path == "/" || r.URL.Path == "/index.html" {
			tmpl, err := template.ParseFiles("static/index.html")
			if err != nil {
				http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
				return
			}
			// Передаем BasePath из конфига в наш HTML
			data := struct{ BasePath string }{
				BasePath: config.Get().Server.BASE_PATH,
			}
			fmt.Println(data.BasePath)
			err = tmpl.Execute(w, data)
			if err != nil {
				// Логируем ошибку на сервере, чтобы вы знали причину
				log.Printf("Ошибка рендеринга шаблона index.html: %v", err)
				// Если заголовки еще не были отправлены, можно вернуть 500 ошибку
				// (Но учтите: если часть HTML уже ушла в браузер, http.Error может сработать некорректно)
				return
			}
			return
		}
		// Если это другие файлы (например, css, js, картинки), отдаем их как обычно
		fs.ServeHTTP(w, r)
	})
	apiMux.Handle("/api/me", auth.AuthMiddleware(http.HandlerFunc(handlers.MeHandler)))
	apiMux.Handle("/api/character/me", auth.AuthMiddleware(http.HandlerFunc(handlers.CheckCharacterHandler)))
	apiMux.Handle("/api/character/create", auth.AuthMiddleware(http.HandlerFunc(handlers.CreateCharacterHandler)))
	apiMux.HandleFunc("/api/register", handlers.RegisterHandler)
	apiMux.HandleFunc("/api/login", handlers.LoginHandler)
	apiMux.HandleFunc("/api/refresh", handlers.RefreshHandler)
	apiMux.Handle("/ws", auth.AuthMiddleware(handlers.WSHandler(gameHub)))
	r := rate.Limit(config.Get().Server.LIMITER_R)
	b := config.Get().Server.LIMITER_B
	limiter := middleware.NewLimiter(r, b)
	mainMux.Handle("/api/", limiter.Limit(apiMux))
	mainMux.Handle("/ws", limiter.Limit(apiMux))
	var MyServer Server = NewHTTPServer(addr, mainMux)
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
