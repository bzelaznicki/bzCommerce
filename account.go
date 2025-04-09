package main

import (
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handleAccountPage(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDContextKey).(string)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	user, err := cfg.db.GetUserById(r.Context(), uuid.MustParse(userID))
	if err != nil {
		cfg.RenderError(w, r, http.StatusNotFound, "User not found")
		log.Printf("user not found: %v", err)
		return
	}

	cfg.Render(w, r, "templates/pages/account.html", user)
}
