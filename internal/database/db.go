package database

import (
	"log"

	_ "github.com/glebarez/go-sqlite"
	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

func InitDB() error {
	var err error
	DB, err = sqlx.Connect("sqlite", "./game.db?_pragma=foreign_keys(1)")
	if err != nil {
		return err
	}
	if err = DB.Ping(); err != nil {
		return err
	}
	createUserTableQuery := `
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username VARCHAR(255) UNIQUE NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
			
		)`
	// Таблица рефрештокенов.

	refreshTokensTableOuery := `
	CREATE TABLE IF NOT EXISTS refresh_tokens (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,
        token TEXT NOT NULL UNIQUE,
        expires_at DATETIME NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);`
	createCharacterTableQuery := `
	CREATE TABLE IF NOT EXISTS caracters (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER UNIQUE NOT NULL,
		name VARCHAR (255) NOT NULL,
		gender VARCHAR (10) NOT NULL,
		avatar_id VARCHAR(50) NOT NULL,
		level INTEGER DEFAULT 1,
		max_hp INTEGER DEFAULT 150,
		hp INTEGER DEFAULT 150,
		max_mana INTEGER DEFAULT 200,
        mana INTEGER DEFAULT 200,      
        location_id VARCHAR(100) DEFAULT 'city_room',
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE 
	);`
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
