package main

import (
	"net/http"

	"github.com/bzelaznicki/bzCommerce/internal/database"
)

func (cfg *apiConfig) handleApiAdminGetProducts(w http.ResponseWriter, r *http.Request) {
	page, limit := getPaginationParams(r)
	offset := (page - 1) * limit

	params := database.ListProductsWithCategoryPaginatedParams{
		Limit:  limit,
		Offset: offset,
	}

	rows, err := cfg.db.ListProductsWithCategoryPaginated(r.Context(), params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to load products")
		return
	}

	totalCount, err := cfg.db.CountProducts(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error counting products")
		return
	}

	products := make([]AdminProductRow, 0, len(rows))
	for _, row := range rows {
		products = append(products, AdminProductRow{
			ID:           row.ID,
			Name:         row.Name,
			Slug:         row.Slug,
			Description:  row.Description.String,
			ImagePath:    row.ImageUrl.String,
			CategoryName: row.CategoryName,
			CategorySlug: row.CategorySlug,
		})
	}

	resp := NewPaginatedResponse(products, page, limit, totalCount)
	respondWithJSON(w, http.StatusOK, resp)
}
