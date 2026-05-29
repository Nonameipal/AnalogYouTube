package pkg

import (
	"fmt"
	"time"

	"github.com/Nonameipal/AnalogYouTube/internal/configs"
	"github.com/dgrijalva/jwt-go"
)

type CustomClaims struct {
	jwt.StandardClaims
	UserID    int    `json:"user_id"`
	Role      string `json:"role"`
	IsRefresh bool   `json:"is_refresh"`
}

func GenerateToken(userID int, ttl int, role string, isRefresh bool) (string, error) {
	claims := CustomClaims{
		StandardClaims: jwt.StandardClaims{},
		UserID:         userID,
		IsRefresh:      isRefresh,
		Role:           role,
	}

	if isRefresh {
		claims.StandardClaims.ExpiresAt = time.Now().Add(time.Duration(ttl) * 24 * time.Hour).Unix()
	} else {
		claims.StandardClaims.ExpiresAt = time.Now().Add(time.Duration(ttl) * time.Minute).Unix()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(configs.AppSettings.AuthParams.JwtSecret))
}

func ParseToken(tokenString string) (int, bool, string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(configs.AppSettings.AuthParams.JwtSecret), nil
	})
	if err != nil {
		return 0, false, "", err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return 0, false, "", fmt.Errorf("invalid token")
	}

	return claims.UserID, claims.IsRefresh, claims.Role, nil
}
