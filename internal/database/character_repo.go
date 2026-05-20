package database

import (
	"database/sql"
	"errors"

	"GoServer/internal/models"
)

func CreateCharacter(char *models.Character) error {
	query := `INSERT INTO characters (user_id, name, gender, avatar_id)
	VALUES (?, ?, ?, ?)`
	_, err := DB.Exec(query, char.UserID, char.Name, char.Gender, char.AvatarID)
	return err
}

func GetCharacterByUserID(userID int64) (*models.Character, error) {
	char := &models.Character{}
	query := "SELECT * FROM characters WHERE user_id = ?"
	err := DB.Get(char, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Персонажа нет, это не ошибка БД
		}
		return nil, err
	}
	return char, nil
}

func GetPlayersInLocation(locationID string) ([]models.Character, error) {
	var players []models.Character
	query := `SELECT id, name, level, avatar_id, gender FROM characters WHERE location_id = ?`
	err := DB.Select(&players, query, locationID)
	return players, err
}

func UpdateCharacterLocation(charID int64, locationID string) error {
	query := `UPDATE characters SET location_id = ? WHERE id = ?`
	_, err := DB.Exec(query, locationID, charID)
	return err
}
