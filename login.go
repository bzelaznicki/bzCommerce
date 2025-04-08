package main

import (
	"log"
	"net/http"
	"time"

	"github.com/bzelaznicki/bzCommerce/internal/auth"
)

func (cfg *apiConfig) handleLoginGet(w http.ResponseWriter, r *http.Request) {
	cfg.Render(w, r, "templates/pages/login.html", nil)
}
func (cfg *apiConfig) handleLoginPost(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	email := r.FormValue("email")
	password := r.FormValue("password")
	if email == "" || password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), email)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		log.Printf("error getting user %s: %v", email, err)
		return
	}
	if err := auth.CheckPassword(user.PasswordHash, password); err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		log.Printf("error authenticating user %s: %v", email, err)
		return
	}

	signedToken, err := auth.GenerateJWT(user.ID.String(), user.Email, cfg.jwtSecret, user.IsAdmin, CookieExpirationTime)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("error signing token for %s: %v", email, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    signedToken,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
func (cfg *apiConfig) handleLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false,
	})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
