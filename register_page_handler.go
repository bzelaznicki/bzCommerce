package main

import (
	"net/http"

	"github.com/bzelaznicki/bzCommerce/internal/auth"
	"github.com/bzelaznicki/bzCommerce/internal/database"
)

func (cfg *apiConfig) handleRegisterGet(w http.ResponseWriter, r *http.Request) {
	cfg.Render(w, r, "templates/pages/register.html", nil)
}

func (cfg *apiConfig) handleRegisterPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	fullName := r.FormValue("full_name")
	email := r.FormValue("email")
	password := r.FormValue("password")
	confirm := r.FormValue("confirm")
	if password != confirm {
		http.Error(w, "Passwords do not match", http.StatusBadRequest)
		return
	}
	if fullName == "" || email == "" || password == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	user := database.CreateUserParams{
		FullName:     fullName,
		Email:        email,
		PasswordHash: hashedPassword,
	}

	_, err = cfg.db.CreateUser(r.Context(), user)
	if err != nil {
		http.Error(w, "User already exists or cannot be created", http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
