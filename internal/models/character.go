package models

type Character struct {
	ID           int64        `db:"id" json:"id"`
	UserID       int64        `db:"user_id" json:"user_id"`
	Name         string       `db:"name" json:"name"`
	Gender       string       `db:"gender" json:"gender"`
	AvatarID     string       `db:"avatar_id" json:"avatar_id"`
	State        PlayerStatus `db:"state" json:"state"` // Состояние
	Level        int          `db:"level" json:"level"`
	Gold         int          `db:"gold" json:"gold"`
	FreePoints   int          `db:"free_points" json:"free_points"`
	Strength     int          `db:"strength" json:"strength"`
	Agility      int          `db:"agility" json:"agility"`
	Intuition    int          `db:"intuition" json:"intuition"`
	Vitality     int          `db:"vitality" json:"vitality"`
	Wisdom       int          `db:"wisdom" json:"wisdom"`
	Charm        int          `db:"charm" json:"charm"`
	HP           float64      `db:"hp" json:"hp"`
	MaxHP        float64      `db:"max_hp" json:"max_hp"`
	Mana         float64      `db:"mana" json:"mana"`
	MaxMana      float64      `db:"max_mana" json:"max_mana"`
	Exp          float64      `db:"exp" json:"exp"`
	MaxExp       float64      `db:"max_exp" json:"max_exp"`
	NextLevelExp float64      `db:"next_level_exp" json:"next_level_exp"`
	WorldID      string       `db:"world_id" json:"world_id"`
	LocationID   string       `db:"location_id" json:"location_id"`
}
