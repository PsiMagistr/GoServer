package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"GoServer/internal/auth"
	"GoServer/internal/database"
	"GoServer/internal/models"
)

type CreateCharacterRequest struct {
	Name   string `json:"name"`
	Gender string `json:"gender"`
	Avatar string `json:"avatar_id"`
}

func CheckCharacterHandler(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	char, err := database.GetCharacterByUserID(claims.UserID)
	if err != nil {
		http.Error(w, "Внутренняя ошибка сервера.", http.StatusInternalServerError)
		return
	}
	if char == nil {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"id":      claims.UserID,
			"message": "Персонаж не создан",
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(char)
}

func CreateCharacterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "CreateCharacter Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.Body == nil {
		http.Error(w, "CreateCharacter Body is empty", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	//////
	var req CreateCharacterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Ошибка в формате данных.", http.StatusBadRequest)
		return
	}
	///////////////////
	_, _ = fmt.Printf("Пришли данные создания персонажа. %+v\n", req)
	// Забираем контекст.
	claims := r.Context().Value(auth.UserContextKey).(*auth.Claims)
	character := &models.Character{
		UserID:   claims.UserID,
		Name:     req.Name,
		Gender:   req.Gender,
		AvatarID: req.Avatar,
	}
	err = database.CreateCharacter(character)
	if err != nil {
		http.Error(w, "Имя персонажа занято.", http.StatusConflict)
		return
	}
	_, _ = fmt.Printf("Персонаж создан. %+v\n", character)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "Персонаж создан"})
}
