package database

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"GoServer/internal/config"
	"GoServer/internal/models"

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

func runInitialSeeds() {
	data, err := os.ReadFile("data/spells.json")
	if err != nil {
		log.Printf("⚠️ Файл заклинаний не найден: %v", err)
		return
	}
	var spells []models.Spell
	err = json.Unmarshal(data, &spells)
	if err != nil {
		log.Printf("❌ Ошибка парсинга JSON: %v", err)
		return
	}
	if len(spells) == 0 {
		return
	}
	query := `INSERT INTO spells (id, name, mana_cost, min_level, element, type, power, description)
	VALUES (:id, :name, :mana_cost, :min_level, :element, :type, :power, :description)
	ON DUPLICATE KEY UPDATE 
			name=VALUES(name), 
			mana_cost=VALUES(mana_cost), 
			power=VALUES(power),
			description=VALUES(description)`
	_, err = DB.NamedExec(query, spells)
	if err != nil {
		log.Printf("❌ Ошибка sqlx bulk insert: %v", err)
		return
	}
	log.Printf("✅ Справочник магии синхронизирован через sqlx (%d записей)", len(spells))
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
	runInitialSeeds()
	log.Println("База данных успешно инициализирована")
	return nil
}
