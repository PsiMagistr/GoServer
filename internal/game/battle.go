package game

import (
	"sync"
	"time"
)

type Battle struct {
	ID              int64
	Attacker        *Client
	Defender        *Client
	AttackerHP      float64
	DefenderHP      float64
	AttackerMana    float64
	DefenderMana    float64
	AttackerMaxHP   float64
	DefenderMaxHP   float64
	AttackerMaxMana float64
	DefenderMaxMana float64
	AttackerTurn    []string
	DefenderTurn    []string
	Round           int
	Log             []string  // Лог боя для этого матча
	ExpiresAt       time.Time // Таймер на совершение хода
	mu              sync.RWMutex
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
