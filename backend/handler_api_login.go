package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/bzelaznicki/bzCommerce/internal/auth"
	"github.com/bzelaznicki/bzCommerce/internal/database"
	"github.com/google/uuid"
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
		ExpiresAt: time.Now().Add(refreshTokenExpiration),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		log.Printf("error adding refresh token: %v", err)
		return
	}

	type response struct {
		User  User   `json:"user"`
		Token string `json:"token"`
	}

	userResponse := User{
		ID:        user.ID,
		FullName:  user.FullName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		IsAdmin:   user.IsAdmin,
	}

	secureCookie := cfg.platform != "dev" // non-secure only on "dev" environment, set in .env
	expiresAt := time.Now().Add(refreshTokenExpiration)
	log.Println("Setting refresh token expiration to:", expiresAt)

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    addedToken.Token,
		Path:     "/",
		HttpOnly: true,
		Secure:   secureCookie,
		SameSite: http.SameSiteStrictMode,
		Expires:  expiresAt,
	})

	resp := response{
		User:  userResponse,
		Token: signedToken,
	}

	cartID, ok := getCartIDFromCookie(r, cfg.cartCookieKey)

	if ok && cartID != uuid.Nil {
		err := cfg.mergeAnonymousCart(r.Context(), cartID, user.ID)
		if err != nil {
			log.Printf("error merging carts: %v", err)
		}
	}

	userCartID, err := cfg.getOrCreateCartIDForUser(r.Context(), user.ID)

	if err == nil {
		cfg.setCartIDCookie(w, userCartID, cfg.cartCookieKey)
	}
	respondWithJSON(w, http.StatusOK, resp)
}

func (cfg *apiConfig) handleApiRefreshToken(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")

	if err != nil || cookie.Value == "" {
		respondWithError(w, http.StatusUnauthorized, "No refresh token provided")
		return
	}

	refreshToken := cookie.Value

	dbToken, err := cfg.db.GetRefreshToken(r.Context(), refreshToken)

	if err != nil || dbToken.ExpiresAt.Before(time.Now().UTC()) {
		respondWithError(w, http.StatusUnauthorized, "Invalid or expired token")
		return
	}

	user, err := cfg.db.GetUserById(r.Context(), dbToken.UserID)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid user")
		return
	}

	signedToken, err := auth.GenerateJWT(user.ID.String(), user.Email, cfg.jwtSecret, user.IsAdmin, CookieExpirationTime)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error generating access token")
		return
	}

	type response struct {
		Token string `json:"token"`
	}

	respondWithJSON(w, http.StatusOK, response{Token: signedToken})
}

func (cfg *apiConfig) handleApiLogout(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("refresh_token")

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Refresh token missing")
		return
	}

	err = cfg.db.InvalidateRefreshToken(r.Context(), cookie.Value)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not revoke token")
		return
	}
	secureCookie := cfg.platform != "dev" // non-secure only on "dev" environment, set in .env

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		HttpOnly: secureCookie,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Unix(0, 0),
	})

	respondWithJSON(w, http.StatusOK, nil)
}
