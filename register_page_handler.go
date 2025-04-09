package main

import (
	"log"
	"net/http"

	"github.com/bzelaznicki/bzCommerce/internal/auth"
	"github.com/bzelaznicki/bzCommerce/internal/database"
)

func (cfg *apiConfig) handleRegisterGet(w http.ResponseWriter, r *http.Request) {
	cfg.Render(w, r, "templates/pages/register.html", nil)
}

func (cfg *apiConfig) handleRegisterPost(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		cfg.RenderError(w, r, http.StatusBadRequest, "Bad Request")
		return
	}

	fullName := r.FormValue("full_name")
	email := r.FormValue("email")
	password := r.FormValue("password")
	confirm := r.FormValue("confirm")
	if password != confirm {
		cfg.RenderError(w, r, http.StatusBadRequest, "Passwords do not match")
		return
	}
	if fullName == "" || email == "" || password == "" {
		cfg.RenderError(w, r, http.StatusBadRequest, "All fields are required")
		return
	}

	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		cfg.RenderError(w, r, http.StatusInternalServerError, "Internal Server Error")
		log.Printf("error during registration: %v", err)
		return
	}

	user := database.CreateUserParams{
		FullName:     fullName,
		Email:        email,
		PasswordHash: hashedPassword,
	}

	_, err = cfg.db.CreateUser(r.Context(), user)
	if err != nil {
		cfg.RenderError(w, r, http.StatusBadRequest, "User already exists or cannot be created")
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
