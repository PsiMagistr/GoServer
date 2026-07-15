package database

import (
	"log"

	"GoServer/internal/models"

	"github.com/jmoiron/sqlx"
)

var AllSpells = make(map[int]models.Spell)

func InitSpells(db *sqlx.DB) error {
	var spells []models.Spell
	query := `SELECT * FROM spells`
	err := db.Select(&spells, query)
	if err != nil {
		return err
	}
	for _, s := range spells {
		AllSpells[s.ID] = s
	}
	log.Printf("Загружено %d заклинаний в справочник.", len(AllSpells))
	return nil
}
