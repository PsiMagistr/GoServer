package models

type Spell struct {
	ID          int     `db:"id" json:"id"`
	CharID      int     `db:"character_id" json:"character_id"`
	SpellID     int     `db:"spell_id" json:"spell_id"`
	Name        string  `db:"name" json:"name"`
	ManaCost    float64 `db:"mana_cost" json:"mana_cost"`
	MinLevel    int     `db:"min_level" json:"min_level"`
	Element     string  `db:"element" json:"element"`
	Type        string  `db:"type" json:"type"`
	Power       float64 `db:"power" json:"power"`
	Description string  `db:"description" json:"description"`
}
