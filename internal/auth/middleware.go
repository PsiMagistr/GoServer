package auth

import (
	"context"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserContextKey contextKey = "user"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Requested-With") == "" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		cookie, err := r.Cookie("access_token")
		if err != nil {
			http.Error(w, "Требуется авторизация (кука не найдена)", http.StatusUnauthorized)
			return
		}
		claims := Claims{}
		tokenString := cookie.Value
		token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Неверный или просроченный токен", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), UserContextKey, &claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
