package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) handleProductPage(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if slug == "" {
		http.NotFound(w, r)
		return
	}

	product, err := cfg.db.GetProductBySlug(r.Context(), slug)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	dbVariants, err := cfg.db.GetProductVariantsByProductId(r.Context(), product.ID)
	if err != nil {
		log.Printf("error fetching product variants: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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

	cfg.Render(w, r, "templates/pages/product.html", productData)
}
