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
	dbHost := config.Get().DB.HOST
	dbPort := config.Get().DB.PORT
	dbUser := config.Get().DB.USER
	dbName := config.Get().DB.NAME
	dbPassword := config.Get().DB.PASSWORD
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

	createCharacterTableQuery := `
	CREATE TABLE IF NOT EXISTS characters (
		id BIGINT PRIMARY KEY AUTO_INCREMENT,
		user_id BIGINT UNIQUE NOT NULL,
		name VARCHAR(255) UNIQUE NOT NULL,
		gender VARCHAR(10) NOT NULL,
		avatar_id VARCHAR(50) NOT NULL,
		level INT DEFAULT 1,
		gold INT DEFAULT 100,
		free_points INT DEFAULT 30,
		strength INT DEFAULT 3,
		agility INT DEFAULT 3,
		intuition INT DEFAULT 3,
		vitality INT DEFAULT 3,
		wisdom INT DEFAULT 3,		
		charm INT DEFAULT 3,
		hp DOUBLE DEFAULT 150.0,
		max_hp DOUBLE DEFAULT 150.0,
		mana DOUBLE DEFAULT 200.0,	
		max_mana DOUBLE DEFAULT 200.0,
		exp DOUBLE DEFAULT 0.0,
    	next_level_exp DOUBLE DEFAULT 1000.0,              
        world_id VARCHAR(100) DEFAULT 'main_city',
    	location_id VARCHAR(100) DEFAULT 'start_glade',
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE 
	)ENGINE=InnoDB;`

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
	log.Println("База данных успешно инициализирована")
	return nil
}
