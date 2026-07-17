package game

import (
	"sync"
	"time"

	"GoServer/internal/models"
)

type Battle struct {
	ID           int64
	AttackerData models.Character
	DefenderData models.Character
	AttackerTurn []int
	DefenderTurn []int
	Round        int
	Finished     bool
	Log          []string  // Лог боя для этого матча
	ExpiresAt    time.Time // Таймер на совершение хода
	mu           sync.RWMutex
}

type BattleFighterDTO struct {
	ID       int64   `json:"id"`
	Name     string  `json:"name"`
	Level    int     `json:"level"`
	HP       float64 `json:"hp"`
	MaxHP    float64 `json:"max_hp"`
	Mana     float64 `json:"mana"`
	MaxMana  float64 `json:"max_mana"`
	AvatarID string  `json:"avatar_id"`
	Gender   string  `json:"gender"`
}

type BattleSnapshot struct {
	BattleID int64            `json:"battle_id"`
	Round    int              `json:"round"`
	TimeLeft int              `json:"time_left"`
	You      BattleFighterDTO `json:"you"`
	Opponent BattleFighterDTO `json:"opponent"`
}

type BattleTurnRequest struct {
	BattleID int64 `json:"battle_id"`
	Round    int   `json:"round"`
	Spells   []int `json:"spells"`
}
