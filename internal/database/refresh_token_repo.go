package database

import (
	"fmt"
	"time"
)

func SaveRefreshToken(userID int, token string, duration time.Duration) error {
	expiresAt := time.Now().Add(duration)
	query := `INSERT INTO refresh_tokens (userID, token, expires_at) VALUES (?, ?, ?)`
	_, err := DB.Exec(query, userID, token, expiresAt)
	return err
}

func GetUserByRefreshToken(token string) (int, error) {
	type result struct {
		UserID    int       `db:"user_id"`
		ExpiresAt time.Time `db:"expires_at"`
	}
	res := result{}
	query := `SELECT user_id, expires_at FROM refresh_tokens WHERE token = ? LIMIT 1`
	err := DB.Get(&res, query, token)
	if err != nil {
		return 0, err
	}
	if time.Now().After(res.ExpiresAt) {
		//_ = DeleteRefreshToken(token)
		return 0, fmt.Errorf("token expired")
	}
	return res.UserID, nil
}

func DeleteRefreshToken(token string) error {
	query := `DELETE FROM refresh_tokens WHERE token = ?`
	_, err := DB.Exec(query, token)
	return err
}
