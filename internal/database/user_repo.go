package database

import (
	"GoServer/internal/models"

	"golang.org/x/crypto/bcrypt"
)

func CreateUser(username, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	query := `INSERT INTO users (username, email, password) VALUES (?, ?, ?)`
	_, err = DB.Exec(query, username, email, string(hashedPassword))
	return err
}

func GetUserByEmail(email string) (*models.User, error) {
	query := `SELECT id, username, email, password FROM users WHERE email = ?`
	u := models.User{}
	err := DB.Get(&u, query, email)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
