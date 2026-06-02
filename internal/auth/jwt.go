package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"GoServer/internal/config"

	"GoServer/internal/database"

	"github.com/golang-jwt/jwt/v5"
)

// var jwtSecret = []byte(j)
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
	jwtSecret := []byte(config.Get().JWT.SECRET)
	accessTokenDuration, err := time.ParseDuration(config.Get().JWT.ACCESSTIME)
	if err != nil {
		accessTokenDuration = 1 * time.Minute
		fmt.Println("В конфиге неверно установленно время аксес-токена")
	}
	now := time.Now()
	expirationTime := now.Add(accessTokenDuration)
	claims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(now), // Время создания
			NotBefore: jwt.NewNumericDate(now), // Нельзя использовать до этого момента
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
	refreshTokenTime, err := time.ParseDuration(config.Get().JWT.REFRESHTIME)
	if err != nil {
		refreshTokenTime = 10 * time.Minute
		fmt.Println("В конфиге неверно установленно время рефреш-токена")
	}
	err = database.SaveRefreshToken(UserID, refreshToken, refreshTokenTime)
	if err != nil {
		return nil, fmt.Errorf("ошибка сохранения refresh токена: %w", err)
	}
	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
