package database

import (
	"fmt"
	"log"

	"GoServer/internal/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

func InitDB() error {
	var err error
	params := config.Get()
	game := params.GAME
	// game config.Config = params.GAME
	dbHost := params.DB.HOST
	dbPort := params.DB.PORT
	dbUser := params.DB.USER
	dbName := params.DB.NAME
	dbPassword := params.DB.PASSWORD
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4",
		dbUser, dbPassword, dbHost, dbPort, dbName)
	// dsn := "root:@tcp(127.0.0.1:3306)/anhat_db?parseTime=true&charset=utf8mb4"
	DB, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		return err
	}
	if err = DB.Ping(); err != nil {
		return err
	}
	createUserTableQuery := `
		CREATE TABLE IF NOT EXISTS users (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			username VARCHAR(255) UNIQUE NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP			
		)ENGINE=InnoDB;`
	// Таблица рефрештокенов.

	refreshTokensTableOuery := `
	CREATE TABLE IF NOT EXISTS refresh_tokens (
        id BIGINT PRIMARY KEY AUTO_INCREMENT,
        user_id BIGINT NOT NULL,
        token VARCHAR(255) NOT NULL, -- Заменили TEXT на VARCHAR(255)
        expires_at DATETIME NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        UNIQUE (token), -- Явно указываем уникальность для VARCHAR
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB;`

	createCharacterTableQuery := fmt.Sprintf(
		`CREATE TABLE IF NOT EXISTS characters (
		id BIGINT PRIMARY KEY AUTO_INCREMENT,
		user_id BIGINT UNIQUE NOT NULL,
		name VARCHAR(255) UNIQUE NOT NULL,
		gender VARCHAR(10) NOT NULL,
		avatar_id VARCHAR(50) NOT NULL,
		state int DEFAULT 0,
		level INT DEFAULT 1,
		gold INT DEFAULT %v,
		free_points INT DEFAULT %v,
		strength INT DEFAULT %v,
		agility INT DEFAULT %v,
		intuition INT DEFAULT %v,
		vitality INT DEFAULT %v,
		wisdom INT DEFAULT %v,		
		charm INT DEFAULT %v,
		hp DOUBLE DEFAULT %v,
		max_hp DOUBLE DEFAULT %v,
		mana DOUBLE DEFAULT %v,	
		max_mana DOUBLE DEFAULT %v,
		exp DOUBLE DEFAULT 30.0,
		max_exp DOUBLE DEFAULT %v,
    	next_level_exp DOUBLE DEFAULT %v,              
        world_id VARCHAR(100) DEFAULT '%v',
    	location_id VARCHAR(100) DEFAULT '%v',
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE 
		)ENGINE=InnoDB;`,
		game.Gold,
		game.FREEPOINTS,
		game.Streпth,
		game.Agility,
		game.Intuition,
		game.Vitality,
		game.Wisdom,
		game.Charm,
		game.MAXHP,
		game.MAXHP,
		game.MAXMANA,
		game.MAXMANA,
		game.MAXEXP,
		game.NEXTEXP,
		game.WORLDID,
		game.LOCATIONID,
	)

	createSpellsTableQuery := `
	CREATE TABLE IF NOT EXISTS spells(
	id INT AUTO_INCREMENT PRIMARY KEY,
	name VARCHAR(100) NOT NULL,
	mana_cost DOUBLE DEFAULT 5.0,
	min_level INT DEFAULT 1,
	element ENUM('fire', 'water', 'air', 'earth') NOT NULL,
	type ENUM('attack', 'shield') NOT NULL,
	power DOUBLE DEFAULT 1.0,         -- Величина воздействия (коэффициент)
    description TEXT
	) ENGINE=InnoDB`

	createCharacterSpellsTableQuery := `
	CREATE TABLE IF NOT EXISTS character_spells (
    character_id BIGINT NOT NULL,
    spell_id INT NOT NULL,
    PRIMARY KEY (character_id, spell_id),
    FOREIGN KEY (character_id) REFERENCES characters(id) ON DELETE CASCADE,
    FOREIGN KEY (spell_id) REFERENCES spells(id) ON DELETE CASCADE
	) ENGINE=InnoDB;`

	_, err = DB.Exec(createUserTableQuery)
	if err != nil {
		log.Println("Ошибка создания таблицы users", err)
		return err // Если ошибка при создании таблицы, выходим сразу
	}
	_, err = DB.Exec(refreshTokensTableOuery)
	if err != nil {
		log.Println("Ошибка создания таблицы refresh-tokens", err)
		return err
	}
	_, err = DB.Exec(createCharacterTableQuery)
	if err != nil {
		log.Println("Ошибка создания таблицы characters", err)
		return err
	}
	_, err = DB.Exec(createSpellsTableQuery)
	if err != nil {
		log.Println("Ошибка создания таблицы spells", err)
		return err
	}
	_, err = DB.Exec(createCharacterSpellsTableQuery)
	if err != nil {
		log.Println("Ошибка создания таблицы character_spells", err)
		return err
	}
	log.Println("База данных успешно инициализирована")
	return nil
}
