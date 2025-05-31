package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/bzelaznicki/bzCommerce/internal/auth"
	"github.com/bzelaznicki/bzCommerce/internal/database"
)

func (cfg *apiConfig) handleApiLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	if params.Email == "" || params.Password == "" {
		respondWithError(w, http.StatusUnauthorized, "Email and password cannot be empty")
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	if err = auth.CheckPassword(user.PasswordHash, params.Password); err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	signedToken, err := auth.GenerateJWT(user.ID.String(), user.Email, cfg.jwtSecret, user.IsAdmin, CookieExpirationTime)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	refreshToken, err := auth.MakeRefreshToken()

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		log.Printf("error generating refresh token: %v", err)
		return
	}

	addedToken, err := cfg.db.AddRefreshToken(r.Context(), database.AddRefreshTokenParams{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().UTC().Add(refreshTokenExpiration),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		log.Printf("error adding refresh token: %v", err)
		return
	}

	type response struct {
		User         User   `json:"user"`
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	userResponse := User{
		ID:        user.ID,
		FullName:  user.FullName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		IsAdmin:   user.IsAdmin,
	}

	resp := response{
		User:         userResponse,
		Token:        signedToken,
		RefreshToken: addedToken.Token,
	}
	respondWithJSON(w, http.StatusOK, resp)
}
