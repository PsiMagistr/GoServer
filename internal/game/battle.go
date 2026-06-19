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
	mu              sync.Mutex
}
