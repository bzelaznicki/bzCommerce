package main

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handleApiGetAccount(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(contextKeyUserID).(string)
	userUUID, err := uuid.Parse(userID)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to decode userId")
		return
	}
	dbUser, err := cfg.db.GetUserById(r.Context(), userUUID)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	type response struct {
		UserId    uuid.UUID `json:"user_id"`
		Email     string    `json:"email"`
		FullName  string    `json:"full_name"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		IsAdmin   bool      `json:"is_admin"`
	}
	resp := response{
		UserId:    dbUser.ID,
		Email:     dbUser.Email,
		FullName:  dbUser.FullName,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		IsAdmin:   dbUser.IsAdmin,
	}

	respondWithJSON(w, http.StatusOK, resp)
}
