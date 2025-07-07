package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bzelaznicki/bzCommerce/internal/auth"
	"github.com/bzelaznicki/bzCommerce/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func (cfg *apiConfig) handleApiAdminGetUsers(w http.ResponseWriter, r *http.Request) {
	page, limit := getPaginationParams(r)

	offset := (page - 1) * limit

	q := r.URL.Query()
	search := q.Get("search")
	sortBy := q.Get("sort_by")
	sortOrder := q.Get("sort_order")

	sortColumn := "created_at"

	switch sortBy {
	case "name":
		sortColumn = "full_name"
	case "email":
		sortColumn = "email"
	case "updated_at":
		sortColumn = "updated_at"
	}

	direction := "ASC"

	if sortOrder == "desc" {
		direction = "DESC"
	}
	// #nosec G201 -- sortColumn and direction are strictly whitelisted to prevent SQL injection
	query := fmt.Sprintf(`
	SELECT
	id, full_name, email, created_at, updated_at, is_admin, is_active, disabled_at
	FROM users 
	WHERE email ILIKE $1 OR full_name ILIKE $1
	ORDER BY %s %s
	LIMIT $2 OFFSET $3
	`, sortColumn, direction)
	rows, err := cfg.sqlDB.QueryContext(r.Context(), query, "%"+search+"%", limit, offset)
	if err != nil {
		log.Print(err)
		respondWithError(w, http.StatusInternalServerError, "Query failed")
		return
	}

	defer rows.Close()

	users := []AdminUserRow{}

	for rows.Next() {
		u := AdminUserRow{}
		if err := rows.Scan(
			&u.ID, &u.FullName, &u.Email, &u.CreatedAt, &u.UpdatedAt, &u.IsAdmin, &u.IsActive, &u.DisabledAt,
		); err != nil {
			respondWithError(w, http.StatusInternalServerError, "Scan failed")
			log.Printf("Error: %v", err)
			return
		}

		users = append(users, u)
	}

	count, err := cfg.db.CountFilteredUsers(r.Context(), "%"+search+"%")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Count failed")
		return
	}

	if users == nil {
		users = []AdminUserRow{}
	}

	resp := NewPaginatedResponse(users, page, limit, count)

	respondWithJSON(w, http.StatusOK, resp)
}

func (cfg *apiConfig) handleApiAdminGetUserDetails(w http.ResponseWriter, r *http.Request) {
	userId, err := uuid.Parse(r.PathValue("userId"))

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := cfg.db.GetUserAccountById(r.Context(), userId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "User not found")
		return
	}
	var disabledAt *string
	if user.DisabledAt.Valid {
		s := user.DisabledAt.Time.Format(time.RFC3339)
		disabledAt = &s
	}

	resp := AdminUserRow{
		ID:         user.ID,
		FullName:   user.FullName,
		Email:      user.Email,
		CreatedAt:  user.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  user.UpdatedAt.Format(time.RFC3339),
		IsAdmin:    user.IsAdmin,
		IsActive:   user.IsActive,
		DisabledAt: disabledAt,
	}
	respondWithJSON(w, http.StatusOK, resp)
}

func (cfg *apiConfig) handleApiAdminUpdateUserDetails(w http.ResponseWriter, r *http.Request) {
	userId, err := uuid.Parse(r.PathValue("userId"))

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := struct {
		FullName string `json:"full_name"`
		Email    string `json:"email"`
		IsAdmin  bool   `json:"is_admin"`
	}{}

	err = decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if params.Email == "" || params.FullName == "" {
		respondWithError(w, http.StatusBadRequest, "Name or email cannot be empty")
		return
	}

	updatedUser, err := cfg.db.UpdateUserById(r.Context(), database.UpdateUserByIdParams{
		ID:       userId,
		FullName: params.FullName,
		Email:    params.Email,
		IsAdmin:  params.IsAdmin,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "User not found")
			return
		}
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			if pqErr.Constraint == "users_email_key" {
				respondWithError(w, http.StatusConflict, "User with this email already exists")
			} else {
				respondWithError(w, http.StatusConflict, "User already exists (duplicate field)")
			}
			return
		}

	}

	respondWithJSON(w, http.StatusOK, updatedUser)
}

func (cfg *apiConfig) handleApiAdminUpdateUserPassword(w http.ResponseWriter, r *http.Request) {
	userId, err := uuid.Parse(r.PathValue("userId"))

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	params := struct {
		NewPassword string `json:"new_password"`
	}{}

	decoder := json.NewDecoder(r.Body)

	err = decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if len(params.NewPassword) < MinPasswordLength {
		respondWithError(w, http.StatusBadRequest, "Password too short")
		return
	}

	hashedPassword, err := auth.HashPassword(params.NewPassword)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error saving password")
		return
	}

	rows, err := cfg.db.UpdateUserPassword(r.Context(), database.UpdateUserPasswordParams{
		ID:       userId,
		Password: hashedPassword,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update password")
		return
	}

	if rows == 0 {
		respondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}

func (cfg *apiConfig) handleApiAdminDeleteUser(w http.ResponseWriter, r *http.Request) {
	userId, err := uuid.Parse(r.PathValue("userId"))

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	rows, err := cfg.db.DeleteUserById(r.Context(), userId)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete user")
		return
	}

	if rows == 0 {
		respondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}

func (cfg *apiConfig) handleApiAdminDisableUser(w http.ResponseWriter, r *http.Request) {
	userId, err := uuid.Parse(r.PathValue("userId"))

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := cfg.db.DisableUser(r.Context(), userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "User not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to disable user")
		return
	}

	var disabledAt *string
	if user.DisabledAt.Valid {
		s := user.DisabledAt.Time.Format(time.RFC3339)
		disabledAt = &s
	}

	resp := AdminUserRow{
		ID:         user.ID,
		FullName:   user.FullName,
		Email:      user.Email,
		CreatedAt:  user.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  user.UpdatedAt.Format(time.RFC3339),
		IsAdmin:    user.IsAdmin,
		IsActive:   user.IsActive,
		DisabledAt: disabledAt,
	}
	respondWithJSON(w, http.StatusOK, resp)
}

func (cfg *apiConfig) handleApiAdminEnableUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := cfg.db.EnableUser(r.Context(), userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "User not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to enable user")
		return
	}

	resp := AdminUserRow{
		ID:        user.ID,
		FullName:  user.FullName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
		IsAdmin:   user.IsAdmin,
		IsActive:  user.IsActive,
	}
	respondWithJSON(w, http.StatusOK, resp)
}
