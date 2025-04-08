package main

import (
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handleAdminUsersList(w http.ResponseWriter, r *http.Request) {
	rows, err := cfg.db.ListUsers(r.Context())

	if err != nil {
		http.Error(w, "Could not load users", http.StatusInternalServerError)
		log.Printf("could not load users: %v", err)
		return
	}

	type AdminUserRow struct {
		ID        uuid.UUID
		FullName  string
		Email     string
		CreatedAt string
		UpdatedAt string
		IsAdmin   bool
	}

	users := make([]AdminUserRow, 0, len(rows))

	for _, row := range rows {
		users = append(users, AdminUserRow{
			ID:        row.ID,
			Email:     row.Email,
			FullName:  row.FullName,
			CreatedAt: row.CreatedAt.Format(TimeFormat),
			UpdatedAt: row.UpdatedAt.Format(TimeFormat),
			IsAdmin:   row.IsAdmin,
		})
	}

	cfg.Render(w, r, "templates/pages/admin_users.html", users)
}
