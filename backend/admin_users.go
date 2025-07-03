package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bzelaznicki/bzCommerce/internal/auth"
	"github.com/bzelaznicki/bzCommerce/internal/database"
	"github.com/google/uuid"
)

type AdminUserRow struct {
	ID        uuid.UUID `json:"id"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	IsAdmin   bool      `json:"is_admin"`
}

type AdminUsersListPageData struct {
	Users       []AdminUserRow
	Breadcrumbs []Breadcrumb
}

type AdminUserFormPageData struct {
	User        User
	Breadcrumbs []Breadcrumb
}

func (cfg *apiConfig) handleAdminUsersList(w http.ResponseWriter, r *http.Request) {
	rows, err := cfg.db.ListUsers(r.Context())

	if err != nil {
		cfg.RenderError(w, r, http.StatusInternalServerError, "Could not load users")
		log.Printf("could not load users: %v", err)
		return
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

	data := AdminUsersListPageData{
		Users: users,
		Breadcrumbs: NewBreadcrumbTrail(
			Breadcrumb{Label: "Users"},
		),
	}

	cfg.Render(w, r, "templates/pages/admin_users.html", data)

}

func (cfg *apiConfig) handleAdminUserEditForm(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		cfg.RenderError(w, r, http.StatusNotFound, fmt.Sprintf("Unable to find user: %v", err))
		log.Printf("Error parsing UUID: %v", err)
		return
	}

	u, err := cfg.db.GetUserById(r.Context(), id)
	if err != nil {
		cfg.RenderError(w, r, http.StatusNotFound, "User not found")
		log.Printf("user not found (id: %s): %v", id, err)
		return
	}

	user := User{
		ID:        u.ID,
		FullName:  u.FullName,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		IsAdmin:   u.IsAdmin,
	}

	data := AdminUserFormPageData{
		User: user,
		Breadcrumbs: NewBreadcrumbTrail(
			Breadcrumb{Label: "Users", URL: "/admin/users"},
			Breadcrumb{Label: "Edit"},
		),
	}

	cfg.Render(w, r, "templates/pages/admin_user_edit.html", data)

}

func (cfg *apiConfig) handleAdminUserUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		cfg.RenderError(w, r, http.StatusNotFound, fmt.Sprintf("Unable to find user: %v", err))
		log.Printf("Error parsing UUID: %v", err)
		return
	}

	if err := r.ParseForm(); err != nil {
		cfg.RenderError(w, r, http.StatusBadRequest, "Invalid form")
		log.Printf("error parsing form: %v", err)
		return
	}

	fullName := r.FormValue("full_name")
	email := r.FormValue("email")
	isAdmin := r.FormValue("is_admin") == "true"
	password := r.FormValue("password")
	confirm := r.FormValue("confirm")

	_, err = cfg.db.UpdateUserById(r.Context(), database.UpdateUserByIdParams{
		ID:       id,
		FullName: fullName,
		Email:    email,
		IsAdmin:  isAdmin,
	})

	if err != nil {
		cfg.RenderError(w, r, http.StatusInternalServerError, "Error updating user")
		log.Printf("error updating user id %s: %v", id, err)
		return
	}

	if password != "" {
		if password != confirm {
			cfg.RenderError(w, r, http.StatusBadRequest, "Passwords don't match")

			return
		}
		hashedPassword, err := auth.HashPassword(password)
		if err != nil {
			cfg.RenderError(w, r, http.StatusInternalServerError, "Internal Server Error")
			log.Printf("error hashing password: %v", err)
			return
		}

		_, err = cfg.db.UpdateUserPassword(r.Context(), database.UpdateUserPasswordParams{
			ID:       id,
			Password: hashedPassword,
		})

		if err != nil {
			cfg.RenderError(w, r, http.StatusInternalServerError, "Failed to update password")
			log.Printf("error updating password: %v", err)
		}
	}

	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}

func (cfg *apiConfig) handleAdminUserDelete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		cfg.RenderError(w, r, http.StatusNotFound, fmt.Sprintf("Unable to find user: %v", err))
		log.Printf("Error parsing UUID: %v", err)
		return
	}

	err = cfg.db.DeleteUserById(r.Context(), id)

	if err != nil {
		cfg.RenderError(w, r, http.StatusInternalServerError, "Failed to delete user")
		log.Printf("failed to delete user (id %s): %v", id, err)
		return
	}

	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}
