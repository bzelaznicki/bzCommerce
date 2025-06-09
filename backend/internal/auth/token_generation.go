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

func ValidateJWT(tokenString, tokenSecret string) (string, error) {
	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})

	if err != nil || !token.Valid {
		return "", fmt.Errorf("invalid token: %v", err)
	}
	exp, ok := claims["exp"].(float64)
	if !ok {
		return "", fmt.Errorf("expiration (exp) not found in token claims")
	}
	if int64(exp) < time.Now().Unix() {
		return "", fmt.Errorf("token has expired")
	}

	userId, ok := claims["user_id"].(string)
	if !ok {
		return "", fmt.Errorf("user_id not found in token claims")
	}

	return userId, nil
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
