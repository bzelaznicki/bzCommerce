package main

import (
	"fmt"
	"net/mail"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func validateUserDetails(email, password string) error {
	if len(password) == 0 {
		return fmt.Errorf("password missing")
	}

	if !validateEmail(email) {
		return fmt.Errorf("invalid email address: %s", email)

	}
	return nil
}

func validateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil

}
