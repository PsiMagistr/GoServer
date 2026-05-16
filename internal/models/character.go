package models

type Character struct {
	ID         int64  `db:"id" json:"id"`
	UserID     int64  `db:"user_id" json:"user_id"`
	Name       string `db:"name" json:"name"`
	Gender     string `db:"gender" json:"gender"`
	AvatarID   string `db:"avatar_id" json:"avatar_id"`
	Level      int    `db:"level" json:"level"`
	HP         int    `db:"hp" json:"hp"`
	MaxHP      int    `db:"max_hp" json:"max_hp"`
	Mana       int    `db:"mana" json:"mana"`
	MaxMana    int    `db:"max_mana" json:"max_mana"`
	WorldID    string `db:"world_id" json:"world_id"`
	LocationID string `db:"location_id" json:"location_id"`
}
