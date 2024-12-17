package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}


func GenerateJWT(userID int64, username string) (string, error) {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "default_secret_key_for_development" 
	}
	claims := TokenClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "mobile-banking-backend",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("gagal membuat token: %v", err)
	}

	return signedToken, nil
}

func VerifyJWT(tokenString string) (*TokenClaims, error) {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "default_secret_key_for_development" 
	}

	token, err := jwt.ParseWithClaims(
		tokenString, 
		&TokenClaims{}, 
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("metode penandatanganan token tidak valid")
			}
			return []byte(secretKey), nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("token tidak valid: %v", err)
	}

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("token tidak dapat diverifikasi")
}

func RefreshJWT(tokenString string) (string, error) {
	claims, err := VerifyJWT(tokenString)
	if err != nil {
		return "", fmt.Errorf("gagal me-refresh token: %v", err)
	}

	newToken, err := GenerateJWT(claims.UserID, claims.Username)
	if err != nil {
		return "", fmt.Errorf("gagal membuat token baru: %v", err)
	}

	return newToken, nil
}