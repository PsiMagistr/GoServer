package database

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

func InitDB() error {
	var err error
	dsn := "root:@tcp(127.0.0.1:3306)/anhat_db?parseTime=true&charset=utf8mb4"
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
		max_hp INT DEFAULT 150,
		hp INT DEFAULT 150,
		max_mana INT DEFAULT 200,
        mana INT DEFAULT 200,      
        location_id VARCHAR(100) DEFAULT 'city_room',
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
