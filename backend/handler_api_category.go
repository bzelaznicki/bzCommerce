package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type CategoryResponse struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Slug        string     `json:"slug"`
	Description string     `json:"description"`
	ParentID    *uuid.UUID `json:"parent_id,omitzero"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (cfg *apiConfig) handleApiGetCategories(w http.ResponseWriter, r *http.Request) {
	dbCategories, err := cfg.db.GetCategories(r.Context())

	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, "Error getting categories")
		return
	}

	categories := make([]CategoryResponse, 0, len(dbCategories))
	for _, dbCategory := range dbCategories {
		var parentID *uuid.UUID
		if dbCategory.ParentID.Valid {
			parentID = &dbCategory.ParentID.UUID
		}
		categories = append(categories, CategoryResponse{
			ID:          dbCategory.ID,
			Name:        dbCategory.Name,
			Slug:        dbCategory.Slug,
			Description: dbCategory.Description.String,
			ParentID:    parentID,
			CreatedAt:   dbCategory.CreatedAt,
			UpdatedAt:   dbCategory.UpdatedAt,
		})
	}
	respondWithJSON(w, http.StatusOK, categories)
}

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
	dbChildren, err := cfg.db.GetChildCategories(r.Context(), uuid.NullUUID{
		UUID:  cat.ID,
		Valid: true,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to load subcategories")
		log.Printf("error loading subcategories for category %s: %v", slug, err)
		return
	}

	children := make([]Category, 0, len(dbChildren))
	for _, dbChild := range dbChildren {
		children = append(children, Category{
			ID:          dbChild.ID,
			Name:        dbChild.Name,
			Slug:        dbChild.Slug,
			Description: dbChild.Description.String,
			ParentID:    dbChild.ParentID.UUID,
			CreatedAt:   dbChild.CreatedAt,
			UpdatedAt:   dbChild.UpdatedAt,
		})
	}

	categoryData := CategoryPageData{
		CategoryName: cat.Name,
		Products:     products,
		Children:     children,
		Breadcrumbs:  breadcrumbs,
	}

	respondWithJSON(w, http.StatusOK, categoryData)
}
