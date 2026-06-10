package models

type PlayerStatus int

const (
	StatusFree = iota
	StatusMoving
	StatusBattle
	StatusDead
)
