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

func UpdateCharacterWorld(charID int64, worldID string, locationID string) error {
	query := `UPDATE characters SET world_id = ?, location_id = ? WHERE id = ?`
	_, err := DB.Exec(query, worldID, locationID, charID)
	return err
}

func UpdateCharacterHpMana(charID int64, hp float64, mana float64) error {
	query := "UPDATE characters SET hp = ?, mana = ? WHERE id = ?"
	_, err := DB.Exec(query, hp, mana, charID)
	return err
}

func UpdateCharacterStats(char *models.Character) error {
	query := "UPDATE characters SET strength = ?, agility = ?, intuition = ?, wisdom = ?, charm = ?, vitality = ?, free_points = ?, max_mana = ?, max_hp = ? WHERE id = ?"
	_, err := DB.Exec(query, char.Strength, char.Agility, char.Intuition, char.Wisdom, char.Charm, char.Vitality, char.FreePoints, char.MaxMana, char.MaxHP, char.ID)
	return err
}

func UpdateCharacter(c *models.Character) error {
	query := `
	UPDATE characters SET
	name = ?,
	level = ?, gold = ?, free_points = ?,
	strength = ?, agility = ?, intuition = ?, vitality = ?, charm = ?, wisdom = ?,
	hp = ?, max_hp = ?, mana = ?, max_mana = ?,
	exp = ?, max_exp = ?, next_level_exp = ?,
	world_id = ?, location_id = ?, state = ?
	WHERE id = ?`
	_, err := DB.Exec(query,
		c.Name,
		c.Level, c.Gold, c.FreePoints,
		c.Strength, c.Agility, c.Intuition, c.Vitality, c.Charm, c.Wisdom,
		c.HP, c.MaxHP, c.Mana, c.MaxMana,
		c.Exp, c.MaxExp, c.NextLevelExp,
		c.WorldID, c.LocationID, c.State,
		c.ID,
	)
	return err
}

func GetCharacterSpells(charID int64) ([]models.Spell, error) {
	var spells []models.Spell
	query := `
	SELECT * FROM spells
	JOIN character_spells ON spells.id = character_spells.spell_id
	WHERE character_spells.character_id = ?`
	err := DB.Select(&spells, query, charID)
	return spells, err
}
