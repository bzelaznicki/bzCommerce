package main

import (
	"encoding/json"
	"net/http"

	"github.com/bzelaznicki/bzCommerce/internal/auth"
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
		User:  userResponse,
		Token: signedToken,
	}
	respondWithJSON(w, http.StatusOK, resp)
}
