package handlers

import (
	"encoding/json"
	"net/http"

	"GoServer/internal/auth"
	"GoServer/internal/database"
)

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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "Персонаж создан"})
}
