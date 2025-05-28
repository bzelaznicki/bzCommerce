package main

import (
	"log"
	"net/http"
	"time"

	"github.com/bzelaznicki/bzCommerce/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handleLoginGet(w http.ResponseWriter, r *http.Request) {
	cfg.Render(w, r, "templates/pages/login.html", nil)
}
func (cfg *apiConfig) handleLoginPost(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		cfg.RenderError(w, r, http.StatusBadRequest, "Bad Request")
		log.Printf("bad request: %v", err)
		return
	}
	email := r.FormValue("email")
	password := r.FormValue("password")
	if email == "" || password == "" {
		cfg.RenderError(w, r, http.StatusBadRequest, "Email and password are required")
		log.Println("email or password is empty")
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), email)
	if err != nil {
		cfg.RenderError(w, r, http.StatusUnauthorized, "Invalid email or password")
		log.Printf("error getting user %s: %v", email, err)
		return
	}
	if err := auth.CheckPassword(user.PasswordHash, password); err != nil {
		cfg.RenderError(w, r, http.StatusUnauthorized, "Invalid email or password")
		log.Printf("error authenticating user %s: %v", email, err)
		return
	}

	signedToken, err := auth.GenerateJWT(user.ID.String(), user.Email, cfg.jwtSecret, user.IsAdmin, CookieExpirationTime)
	if err != nil {
		cfg.RenderError(w, r, http.StatusInternalServerError, "Internal Server Error")
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
	cartID, ok := getCartIDFromCookie(r, cfg.cartCookieKey)
	if ok && cartID != uuid.Nil {
		err := cfg.mergeAnonymousCart(r.Context(), cartID, user.ID)
		if err != nil {
			log.Printf("error merging carts: %v", err)
		}
	}

	userCartID, err := cfg.getOrCreateCartIDForUser(r.Context(), user.ID)
	if err == nil {
		setCartIDCookie(w, userCartID, cfg.cartCookieKey)
	}

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

	clearCartIDCookie(w)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
