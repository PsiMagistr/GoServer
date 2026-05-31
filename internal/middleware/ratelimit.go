package middleware

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type IPLimiter struct {
	visitors map[string]*visitor
	mu       sync.Mutex
	r        rate.Limit
	b        int
}

func NewLimiter(r rate.Limit, b int) *IPLimiter {
	i := &IPLimiter{
		visitors: make(map[string]*visitor),
		r:        r,
		b:        b,
	}
	go i.cleanupVisitors()
	return i
}

func (i *IPLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()
	_, exists := i.visitors[ip]
	if !exists {
		limiter := rate.NewLimiter(i.r, i.b)
		i.visitors[ip] = &visitor{
			limiter:  limiter,
			lastSeen: time.Now(),
		}
	}
	i.visitors[ip].lastSeen = time.Now()
	return i.visitors[ip].limiter
}

func (i *IPLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, "Не удалось определить адрес отправителя", http.StatusInternalServerError)
			return
		}
		limiter := i.GetLimiter(ip)
		if !limiter.Allow() {
			fmt.Println("Попытка ddos атаки")
			http.Error(w, "Слишком много запросов. Подождите немного.", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (i *IPLimiter) cleanupVisitors() {
	for {
		time.Sleep(time.Hour)
		i.mu.Lock()
		fmt.Println("Запуск плановой очистки Rate Limiter...")
		for ip, v := range i.visitors {
			if time.Since(v.lastSeen) > time.Hour {
				delete(i.visitors, ip)
			}
		}
		fmt.Printf("Очистка завершена. Активных IP в памяти: %d", len(i.visitors))
		i.mu.Unlock()
	}
}
