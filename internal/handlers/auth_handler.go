package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"GoServer/internal/auth"
	"GoServer/internal/database"

	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Username        string `json:"username"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func GetCookieParams(name string, value string, maxAge int) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   maxAge,
	}
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.Body == nil {
		http.Error(w, "Body is empty", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()
	var req RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	fmt.Printf("Пришли данные регистрации: %+v\n", req)
	if err != nil {
		http.Error(w, "Ошибка в формате данных.", http.StatusBadRequest)
		return
	}

	if req.ConfirmPassword != req.Password {
		http.Error(w, "Пароли не совпадают!", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Email == "" || req.Password == "" {
		http.Error(w, "Все поля должны быть заполнены", http.StatusBadRequest)
		return
	}
	err = database.CreateUser(req.Username, req.Email, req.Password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Пользователь с таким именем или email уже существует", http.StatusConflict)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "Регистрация успешна!"})
}

// LoginHandler

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.Body == nil {
		http.Error(w, "Body is empty", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()
	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Ошибка в формате данных.", http.StatusBadRequest)
		return
	}
	fmt.Printf("Пришли данные логина: %+v\n", req)
	user, err := database.GetUserByEmail(req.Email)
	if err != nil {
		http.Error(w, "Неверный логин или пароль.", http.StatusUnauthorized)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		http.Error(w, "Неверный логин или пароль.", http.StatusBadRequest)
		return
	}
	tokens, err := auth.GetTokenPair(user.ID, user.Username)
	if err != nil {
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}
	accessCookie := GetCookieParams("access_token", tokens.AccessToken, 365*24*3600)
	refreshCookie := GetCookieParams("refresh_token", tokens.RefreshToken, 365*24*3600)

	http.SetCookie(w, accessCookie)
	http.SetCookie(w, refreshCookie)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"username": user.Username,
		"id":       user.ID,
	})
}

func MeHandler(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	w.Header().Set("Content-Type", "application/json")
	_, _ = fmt.Fprintf(w, `{"id": %d, "username": "%s", "status": "online"}`, claims.UserID, claims.Username)
}

// Рефрещш-токен, проверка.
func RefreshHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "Рефреш-токен не найден", http.StatusUnauthorized)
		return
	}
	user, err := database.GetUserByRefreshToken(cookie.Value)
	if err != nil {
		http.Error(w, "Сессия истекла, войдите заново", http.StatusUnauthorized)
		return
	}
	_ = database.DeleteRefreshToken(cookie.Value)
	tokens, err := auth.GetTokenPair(user.ID, user.Username)
	if err != nil {
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}
	accessCookie := GetCookieParams("access_token", tokens.AccessToken, 365*24*3600)
	refreshCookie := GetCookieParams("refresh_token", tokens.RefreshToken, 365*24*3600)
	http.SetCookie(w, accessCookie)
	http.SetCookie(w, refreshCookie)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "refreshed"})
}
