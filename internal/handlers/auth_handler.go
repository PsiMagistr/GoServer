package handlers

import (
	"GoServer/internal/auth"
	"GoServer/internal/database"
	"encoding/json"
	"fmt"
	"net/http"

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

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
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

//LoginHandler

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Ошибка в формате данных.", http.StatusBadRequest)
		return
	}
	fmt.Printf("Пришли данные логина: %+v\n", req)
	user, err := database.GetUserByEmail(req.Email)
	if err != nil {
		http.Error(w, "Неверный логин или пароль.", http.StatusBadRequest)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		http.Error(w, "Неверный логин или пароль.", http.StatusBadRequest)
		return
	}
	token, err := auth.GenerateToken(user.ID, user.Username)
	if err != nil {
		fmt.Println("JWT Error:", err)
		http.Error(w, "Ошибка создания сессии.", http.StatusConflict)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"token":    token,
	}
	_ = json.NewEncoder(w).Encode(response)
}

func MeHandler(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	w.Header().Set("Content-Type", "application/json")
	_, _ = fmt.Fprintf(w, `{"id": %d, "username": "%s", "status": "online"}`, claims.UserID, claims.Username)
}
