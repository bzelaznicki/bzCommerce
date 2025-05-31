package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/google/uuid"
)

//TODO: refactor into separate calls for breadcrumbs and child categories

func (cfg *apiConfig) handleApiGetCategoryProducts(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	if slug == "" {
		respondWithError(w, http.StatusNotFound, "Category not found")
		return
	}

	cat, err := cfg.db.GetCategoryBySlug(r.Context(), slug)

	if err != nil {
		respondWithError(w, http.StatusNotFound, "Category not found")
		return
	}

	breadcrumbs, err := cfg.db.GetCategoryPathByID(r.Context(), uuid.NullUUID{UUID: cat.ID, Valid: true})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	for i, j := 0, len(breadcrumbs)-1; i < j; i, j = i+1, j-1 {
		breadcrumbs[i], breadcrumbs[j] = breadcrumbs[j], breadcrumbs[i]
	}

	nsSlug := sql.NullString{String: slug, Valid: true}

	productRows, err := cfg.db.ListProductsByCategoryRecursive(r.Context(), nsSlug)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
		log.Printf("error loading recursive products: %v", err)
		return
	}

	products := make([]Product, 0, len(productRows))
	for _, dbProduct := range productRows {
		products = append(products, Product{
			ID:          dbProduct.ID,
			Name:        dbProduct.Name,
			Slug:        dbProduct.Slug,
			ImagePath:   dbProduct.ImageUrl.String,
			CategoryID:  dbProduct.CategoryID,
			Description: dbProduct.Description.String,
		})
	}
	children, err := cfg.db.GetChildCategories(r.Context(), uuid.NullUUID{
		UUID:  cat.ID,
		Valid: true,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to load subcategories")
		log.Printf("error loading subcategories for category %s: %v", slug, err)
		return
	}

	categoryData := CategoryPageData{
		CategoryName: cat.Name,
		Products:     products,
		Children:     children,
		Breadcrumbs:  breadcrumbs,
	}

	respondWithJSON(w, http.StatusOK, categoryData)
}
