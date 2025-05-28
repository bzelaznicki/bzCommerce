package main

import (
	"log"
	"net/http"

	"github.com/bzelaznicki/bzCommerce/internal/database"
	"github.com/google/uuid"
)

type ProductPageData struct {
	Product     Product
	Breadcrumbs []database.GetCategoryPathByIDRow
}

func (cfg *apiConfig) handleProductPage(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if slug == "" {
		cfg.RenderError(w, r, http.StatusNotFound, "Page not found")
		return
	}

	product, err := cfg.db.GetProductBySlug(r.Context(), slug)
	if err != nil {
		cfg.RenderError(w, r, http.StatusNotFound, "Page not found")
		return
	}

	breadcrumbs, err := cfg.db.GetCategoryPathByID(r.Context(), uuid.NullUUID{UUID: product.CategoryID, Valid: true})

	if err != nil {
		cfg.RenderError(w, r, http.StatusInternalServerError, "Internal Server Error")
		log.Printf("failed to get breadcrumbs: %v", err)
	}

	for i, j := 0, len(breadcrumbs)-1; i < j; i, j = i+1, j-1 {
		breadcrumbs[i], breadcrumbs[j] = breadcrumbs[j], breadcrumbs[i]
	}

	dbVariants, err := cfg.db.GetProductVariantsByProductId(r.Context(), product.ID)
	if err != nil {
		log.Printf("error fetching product variants: %v", err)
		cfg.RenderError(w, r, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	variants := make([]Variant, 0, len(dbVariants))
	for _, dbVariant := range dbVariants {
		variant := Variant{
			ID:            dbVariant.ID,
			Name:          dbVariant.VariantName.String,
			Price:         dbVariant.Price,
			StockQuantity: dbVariant.StockQuantity,
			ImageUrl:      dbVariant.ImageUrl.String,
			VariantName:   dbVariant.VariantName.String,
			CreatedAt:     dbVariant.CreatedAt,
			UpdatedAt:     dbVariant.UpdatedAt,
		}
		variants = append(variants, variant)
	}

	if len(variants) == 0 {
		log.Printf("no variants found for product ID %d", product.ID)
	}
	productData := Product{
		ID:          product.ID,
		Name:        product.Name,
		Slug:        product.Slug,
		ImagePath:   product.ImageUrl.String,
		CategoryID:  product.CategoryID,
		Description: product.Description.String,
		Variants:    variants,
	}

	cfg.Render(w, r, "templates/pages/product.html", ProductPageData{
		Product:     productData,
		Breadcrumbs: breadcrumbs,
	})

}
