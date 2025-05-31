package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/bzelaznicki/bzCommerce/internal/database"
	"github.com/google/uuid"
)

type CategoryPageData struct {
	CategoryName string                            `json:"category_name"`
	Products     []Product                         `json:"products"`
	Children     []Category                        `json:"children"`
	Breadcrumbs  []database.GetCategoryPathByIDRow `json:"breadcrumbs"`
}

func (cfg *apiConfig) handleCategoryPage(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if slug == "" {
		cfg.RenderError(w, r, http.StatusNotFound, "Category not found")
		return
	}

	cat, err := cfg.db.GetCategoryBySlug(r.Context(), slug)
	if err != nil {
		cfg.RenderError(w, r, http.StatusNotFound, "Category not found")
		log.Printf("error finding category with slug '%s': %v", slug, err)
		return
	}

	breadcrumbs, err := cfg.db.GetCategoryPathByID(r.Context(), uuid.NullUUID{UUID: cat.ID, Valid: true})
	if err != nil {
		cfg.RenderError(w, r, http.StatusInternalServerError, "Internal Server Error")
		log.Printf("failed to get breadcrumbs: %v", err)
	}

	for i, j := 0, len(breadcrumbs)-1; i < j; i, j = i+1, j-1 {
		breadcrumbs[i], breadcrumbs[j] = breadcrumbs[j], breadcrumbs[i]
	}

	nsSlug := sql.NullString{String: slug, Valid: true}

	productRows, err := cfg.db.ListProductsByCategoryRecursive(r.Context(), nsSlug)
	if err != nil {
		cfg.RenderError(w, r, http.StatusInternalServerError, "Internal Server Error")
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
		cfg.RenderError(w, r, http.StatusInternalServerError, "Failed to load subcategories")
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

	cfg.Render(w, r, "templates/pages/category.html", categoryData)
}
