package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/bzelaznicki/bzCommerce/internal/auth"
	"github.com/bzelaznicki/bzCommerce/internal/database"
)

func (cfg *apiConfig) handlerApiRegister(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		FullName string `json:"full_name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding parameters")
		return
	}

	// Sanitize inputs
	params.FullName = sanitizeName(params.FullName)
	params.Email = strings.ToLower(strings.TrimSpace(params.Email))

	if params.FullName == "" {
		respondWithError(w, http.StatusBadRequest, "Name cannot be empty")
		return
	}

	if params.Email == "" || !isValidEmail(params.Email) {
		respondWithError(w, http.StatusBadRequest, "Invalid email address")
		return
	}

	if len(params.Password) < MinPasswordLength {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Password must be at least %d characters long", MinPasswordLength))
		return
	}

	if !validatePasswordStrength(params.Password) {
		respondWithError(w, http.StatusBadRequest, "Password must contain at least 3 of the following: uppercase letter, lowercase letter, number, special character")
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	user := database.CreateUserParams{
		FullName:     params.FullName,
		Email:        params.Email,
		PasswordHash: hashedPassword,
	}

	dbUser, err := cfg.db.CreateUser(r.Context(), user)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "User already exists or cannot be created")
		return
	}

	userResponse := User{
		ID:        dbUser.ID,
		Email:     dbUser.Email,
		FullName:  dbUser.FullName,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		IsAdmin:   dbUser.IsAdmin,
	}

	respondWithJSON(w, http.StatusCreated, userResponse)
}
