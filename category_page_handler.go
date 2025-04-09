package main

import (
	"database/sql"
	"log"
	"net/http"
)

type CategoryPageData struct {
	CategoryName string
	Products     []Product
}

func (cfg *apiConfig) handleCategoryPage(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	if slug == "" {
		http.NotFound(w, r)
		return
	}

	nsSlug := sql.NullString{
		String: slug,
		Valid:  true,
	}
	category, err := cfg.db.ListProductsByCategory(r.Context(), nsSlug)

	if err != nil {
		cfg.RenderError(w, r, http.StatusInternalServerError, "Internal Server Error")
		log.Printf("error loading category page: %v", err)
		return
	}
	products := make([]Product, 0, len(category))
	for _, dbProduct := range category {
		product := Product{
			ID:          dbProduct.ID,
			Name:        dbProduct.Name,
			Slug:        dbProduct.Slug,
			ImagePath:   dbProduct.ImageUrl.String,
			CategoryID:  dbProduct.CategoryID,
			Description: dbProduct.Description.String,
		}
		products = append(products, product)
	}
	categoryData := CategoryPageData{
		CategoryName: slug,
		Products:     products,
	}
	cfg.Render(w, r, "templates/pages/category.html", categoryData)

}
