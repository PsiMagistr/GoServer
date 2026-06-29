package config

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type Server struct {
	IP        string  `json:"ip"`
	HOST      string  `json:"host"`
	LIMITER_R float64 `json:"limiter_r"`
	LIMITER_B int     `json:"limiter_b"`
}

type DB struct {
	PORT     string `json:"port"`
	HOST     string `json:"host"`
	NAME     string `json:"name"`
	USER     string `json:"user"`
	PASSWORD string `json:"password"`
}

type JWT struct {
	SECRET      string `json:"jwt_secret"`
	ACCESSTIME  string `json:"access_token_time"`
	REFRESHTIME string `json:"refresh_token_time"`
}

type GAME struct {
	MINCHARLEN int     `json:"min_char_len"`
	MAXCHARLEN int     `json:"max_char_len"`
	MAXHP      float64 `json:"max_hp"`
	MAXMANA    float64 `json:"max_mana"`
	MAXEXP     float64 `json:"max_exp"`
	NEXTEXP    float64 `json:"next_exp"`
	FREEPOINTS int     `json:"free_points"`
	Gold       int     `json:"gold"`
	Streпth    int     `json:"strength"`
	Agility    int     `json:"agility"`
	Vitality   int     `json:"Vitality"`
	Intuition  int     `json:"intuition"`
	Wisdom     int     `json:"wisdom"`
	Charm      int     `json:"charm"`
	WORLDID    string  `json:"world_id"`
	LOCATIONID string  `json:"location_id"`
	ROUNDTIME  int     `json:"round_time"`
}

type Config struct {
	Server Server `json:"server"`
	DB     DB     `json:"db"`
	JWT    JWT    `json:"jwt"`
	GAME   GAME   `json:"game"`
}

var (
	once     sync.Once
	instance *Config
)

func LoadFile(name string) error {
	var err error
	once.Do(func() {
		var file *os.File
		file, err = os.Open(name)
		if err != nil {
			fmt.Println("Ошибка открытия конфига", err)
			return
		}
		defer file.Close()
		var cfg Config
		decoder := json.NewDecoder(file)
		err = decoder.Decode(&cfg)
		if err != nil {
			return
		}
		instance = &cfg
	})
	return err
}

func Get() *Config {
	if instance == nil {
		panic("Конфигурация не инициализирована. Вызовите config.LoadFIle")
	}
	return instance
}
