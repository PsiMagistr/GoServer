package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"GoServer/internal/database"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("my_super_secret_rpg_key_12345")

type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

func GenerateToken(userID int64, username string) (string, error) {
	expirationTime := time.Now().Add(1 * time.Minute)
	claims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func GetTokenPair(UserID int64, username string) (*TokenPair, error) {
	accessToken, err := GenerateToken(UserID, username)
	if err != nil {
		return nil, fmt.Errorf("Ошибка генерации аксесс-токена: %w", err)
	}
	refreshToken, err := GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("ошибка генерации рефреш-токена: %w", err)
	}
	err = database.SaveRefreshToken(UserID, refreshToken, 10*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("ошибка сохранения refresh токена: %w", err)
	}
	return &TokenPair{
		accessToken,
		refreshToken,
	}, nil
}
