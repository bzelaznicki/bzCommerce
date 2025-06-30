package main

import (
	"fmt"
	"log"
	"net/http"
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

	query := fmt.Sprintf(`
	SELECT
	id, full_name, email, created_at, updated_at, is_admin
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
			&u.ID, &u.FullName, &u.Email, &u.CreatedAt, &u.UpdatedAt, &u.IsAdmin,
		); err != nil {
			respondWithError(w, http.StatusInternalServerError, "Scan failed")
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
