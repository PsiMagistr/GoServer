package database

import (
	"fmt"
	"log"

	"GoServer/internal/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

func ensureDatabaseExists(user, password, host, port, dbName string) error {
	// Подключаемся к MySQL БЕЗ указания имени базы
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/?parseTime=true", user, password, host, port)
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return fmt.Errorf("ошибка подключения к серверу для создания БД: %w", err)
	}
	defer db.Close()
	// Создаем базу, если её нет
	query := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", dbName)
	_, err = db.Exec(query)
	return err
}

func seedTable(tableName string, query string) {
	var count int
	err := DB.Get(&count, fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName))
	if err != nil {
		log.Printf("Ошибка проверки таблицы %s: %v", tableName, err)
		return
	}
	if count == 0 {
		log.Printf("Таблица %s пуста. Загрузка начальных данных...", tableName)
		_, err := DB.Exec(query)
		if err != nil {
			log.Printf("Ошибка при сидинге таблицы %s: %v", tableName, err)
		}
	}
}

func createTables(queries []string) error {
	for _, q := range queries {
		_, err := DB.Exec(q)
		if err != nil {
			return fmt.Errorf("ошибка выполнения запроса: %v \nЗапрос: %s", err, q)
		}
	}
	return nil
}

func InitDB() error {
	var err error
	params := config.Get()
	dbCfg := params.DB
	err = ensureDatabaseExists(dbCfg.USER, dbCfg.PASSWORD, dbCfg.HOST, dbCfg.PORT, dbCfg.NAME)
	if err != nil {
		return err
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4",
		dbCfg.USER, dbCfg.PASSWORD, dbCfg.HOST, dbCfg.PORT, dbCfg.NAME)

	DB, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		return err
	}
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("база данных не отвечает на пинг: %w", err)
	}
	schema := GetSchema()
	err = createTables(schema)
	if err != nil {
		return err
	}
	log.Println("База данных успешно инициализирована")
	return nil
}
