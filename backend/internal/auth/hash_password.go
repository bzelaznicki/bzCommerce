package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", fmt.Errorf("error generating password: %v", err)
	}
	hashedString := string(hashedPassword)

	return hashedString, nil

}

func CheckPassword(userPassword, hashedInput string) error {

	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(hashedInput))
	if err != nil {
		return fmt.Errorf("passwords don't match: %v", err)
	}
	return nil
}
