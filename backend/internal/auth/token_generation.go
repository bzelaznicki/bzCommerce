package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateJWT(userId, email, tokenSecret string, isAdmin bool, tokenExpiration time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  userId,
		"email":    email,
		"is_admin": isAdmin,
		"exp":      time.Now().Add(tokenExpiration).Unix(),
	})

	signedToken, err := token.SignedString([]byte(tokenSecret))

	if err != nil {
		return "", fmt.Errorf("error signing token: %v", err)
	}

	return signedToken, nil
}

func ValidateJWT(tokenString, tokenSecret string) (string, bool, error) {
	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})

	if err != nil || !token.Valid {
		return "", false, fmt.Errorf("invalid token: %v", err)
	}
	exp, ok := claims["exp"].(float64)
	if !ok {
		return "", false, fmt.Errorf("expiration (exp) not found in token claims")
	}
	if int64(exp) < time.Now().Unix() {
		return "", false, fmt.Errorf("token has expired")
	}

	userId, ok := claims["user_id"].(string)
	if !ok {
		return "", false, fmt.Errorf("user_id not found in token claims")
	}

	isAdmin, _ := claims["is_admin"].(bool)

	return userId, isAdmin, nil
}

func MakeRefreshToken() (string, error) {

	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)

	if err != nil {
		return "", fmt.Errorf("error generating random key: %s", err)
	}
	randomString := hex.EncodeToString(bytes)

	return randomString, nil
}
