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

func MakeRefreshToken() (string, error) {

	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)

	if err != nil {
		return "", fmt.Errorf("error generating random key: %s", err)
	}
	randomString := hex.EncodeToString(bytes)

	return randomString, nil
}
