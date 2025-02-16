package main

import (
	"encoding/json"
	"net/http"

	"github.com/bzelaznicki/bzCommerce/internal/database"
	"github.com/google/uuid"
)

func (cfg apiConfig) handleUserCreate(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		FullName string `json:"full_name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding JSON data", err)
		return
	}

	err = validateUserDetails(params.Email, params.Password)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid credentials, check username and password", err)
		return
	}

	createUser := database.CreateUserParams{
		ID:           uuid.New(),
		Email:        params.Email,
		FullName:     params.FullName,
		PasswordHash: "UPDATEME",
	}
	dbUser, err := cfg.db.CreateUser(r.Context(), createUser)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create user", err)
		return
	}

	user := User{
		ID:        dbUser.ID,
		FullName:  dbUser.FullName,
		Email:     dbUser.Email,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
	}
	respondWithJSON(w, http.StatusCreated, user)
}
