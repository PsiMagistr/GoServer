package database

import (
	"fmt"

	"GoServer/internal/config"
)

// GetSchema возвращает слайс всех запросов для создания таблиц в правильном порядке
func GetSchema() []string {
	game := config.Get().GAME
	return []string{
		// 1. Users
		`CREATE TABLE IF NOT EXISTS users (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			username VARCHAR(255) UNIQUE NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP			
		) ENGINE=InnoDB;`,

		// 2. Tokens
		`CREATE TABLE IF NOT EXISTS refresh_tokens (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			user_id BIGINT NOT NULL,
			token VARCHAR(255) NOT NULL UNIQUE,
			expires_at DATETIME NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		) ENGINE=InnoDB;`,

		// 3. Characters (с подстановкой данных из конфига)
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS characters (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			user_id BIGINT UNIQUE NOT NULL,
			name VARCHAR(255) UNIQUE NOT NULL,
			gender VARCHAR(10) NOT NULL,
			avatar_id VARCHAR(50) NOT NULL,
			state INT DEFAULT 0,
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
			exp DOUBLE DEFAULT 0.0,
			max_exp DOUBLE DEFAULT %v,
			next_level_exp DOUBLE DEFAULT %v,              
			world_id VARCHAR(100) DEFAULT '%s',
			location_id VARCHAR(100) DEFAULT '%s',
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE 
		) ENGINE=InnoDB;`,
			game.Gold, game.FREEPOINTS, game.Streпth, game.Agility,
			game.Intuition, game.Vitality, game.Wisdom, game.Charm,
			game.MAXHP, game.MAXHP, game.MAXMANA, game.MAXMANA,
			game.MAXEXP, game.NEXTEXP, game.WORLDID, game.LOCATIONID),

		// 4. Spells
		`CREATE TABLE IF NOT EXISTS spells (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			mana_cost DOUBLE DEFAULT 5.0,
			min_level INT DEFAULT 1,
			element ENUM('fire', 'water', 'air', 'earth') NOT NULL,
			type ENUM('attack', 'shield') NOT NULL,
			power DOUBLE DEFAULT 1.0,
			description TEXT			
		) ENGINE=InnoDB;`,

		// 5. Character Spells
		`CREATE TABLE IF NOT EXISTS character_spells (
			character_id BIGINT NOT NULL,
			spell_id INT NOT NULL,
			PRIMARY KEY (character_id, spell_id),
			FOREIGN KEY (character_id) REFERENCES characters(id) ON DELETE CASCADE,
			FOREIGN KEY (spell_id) REFERENCES spells(id) ON DELETE CASCADE
		) ENGINE=InnoDB;`,
	}
}
