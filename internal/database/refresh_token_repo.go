package database

import (
	"GoServer/internal/models"
	"fmt"
	"time"
)

func SaveRefreshToken(userID int, token string, duration time.Duration) error {
	expiresAt := time.Now().Add(duration)
	query := `INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES (?, ?, ?)`
	_, err := DB.Exec(query, userID, token, expiresAt)
	return err
}

func GetUserByRefreshToken(token string) (*models.User, error) {
	type result struct {
		models.User
		ExpiresAt time.Time `db:"expires_at"`
	}
	res := &result{}
	query := `
	SELECT u.id, u.username, u.email, rt.expires_at
    FROM users AS u
	JOIN refresh_tokens AS rt
	ON u.id = rt.user_id
	WHERE rt.token = ?
	LIMIT 1`
	err := DB.Get(res, query, token)
	if err != nil {
		return nil, err
	}
	if time.Now().After(res.ExpiresAt) {
		return nil, fmt.Errorf("refresh token expired")
	}
	return &res.User, nil
}

func DeleteRefreshToken(token string) error {
	query := `DELETE FROM refresh_tokens WHERE token = ?`
	_, err := DB.Exec(query, token)
	return err
}
