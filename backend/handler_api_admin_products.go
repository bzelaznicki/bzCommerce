package main

import (
	"database/sql"
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handleApiAdminGetProducts(w http.ResponseWriter, r *http.Request) {
	page, limit := getPaginationParams(r)
	offset := (page - 1) * limit

	q := r.URL.Query()
	search := q.Get("search")
	sortBy := q.Get("sort_by")
	sortOrder := q.Get("sort_order")

	sortColumn := "p.created_at"
	switch sortBy {
	case "name":
		sortColumn = "p.name"
	case "category_name":
		sortColumn = "c.name"
	case "updated_at":
		sortColumn = "p.updated_at"
	}

	direction := "ASC"
	if sortOrder == "desc" {
		direction = "DESC"
	}
	// #nosec G201 -- sortColumn and direction are strictly whitelisted to prevent SQL injection
	query := fmt.Sprintf(`
		SELECT
		p.id, p.name, p.slug, p.description, p.image_url AS image_path,
		c.name AS category_name, c.slug AS category_slug,
		p.created_at, p.updated_at
		FROM products p
		JOIN categories c ON p.category_id = c.id
		WHERE p.name ILIKE $1
		ORDER BY %s %s
		LIMIT $2 OFFSET $3`, sortColumn, direction)

	rows, err := cfg.sqlDB.QueryContext(r.Context(), query, "%"+search+"%", limit, offset)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Query failed")
		return
	}
	defer rows.Close()

	var products []AdminProductRow
	for rows.Next() {
		var p AdminProductRow
		var desc, image sql.NullString
		if err := rows.Scan(
			&p.ID, &p.Name, &p.Slug, &desc, &image,
			&p.CategoryName, &p.CategorySlug,
			&p.CreatedAt, &p.UpdatedAt,
		); err != nil {

			respondWithError(w, http.StatusInternalServerError, "Scan failed")
			return
		}
		p.Description = desc.String
		p.ImagePath = image.String
		products = append(products, p)
	}

	count, err := cfg.db.CountFilteredProducts(r.Context(), "%"+search+"%")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Count failed")
		return
	}

	resp := NewPaginatedResponse(products, page, limit, count)
	respondWithJSON(w, http.StatusOK, resp)
}
