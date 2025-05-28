package main

import (
	"net/http"

	"log"

	"github.com/bzelaznicki/bzCommerce/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handleApiGetProducts(w http.ResponseWriter, r *http.Request) {
	dbProducts, err := cfg.db.ListProducts(r.Context())

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		log.Printf("Error rendering products: %v", err)
		return
	}

	products := make([]Product, 0, len(dbProducts))

	for _, dbProduct := range dbProducts {
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

	respondWithJSON(w, http.StatusOK, products)
}

func (cfg *apiConfig) handleApiGetSingleProduct(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	product, err := cfg.db.GetProductBySlug(r.Context(), slug)

	if err != nil {
		respondWithError(w, http.StatusNotFound, "Page not found")
		return
	}

	breadcrumbs, err := cfg.db.GetCategoryPathByID(r.Context(), uuid.NullUUID{UUID: product.CategoryID, Valid: true})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get category path")
		return
	}

	for i, j := 0, len(breadcrumbs)-1; i < j; i, j = i+1, j-1 {
		breadcrumbs[i], breadcrumbs[j] = breadcrumbs[j], breadcrumbs[i]
	}

	dbVariants, err := cfg.db.GetProductVariantsByProductId(r.Context(), product.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get product variants")
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

	productData := Product{
		ID:          product.ID,
		Name:        product.Name,
		Slug:        product.Slug,
		ImagePath:   product.ImageUrl.String,
		CategoryID:  product.CategoryID,
		Description: product.Description.String,
		Variants:    variants,
	}

	type response struct {
		Product     Product                           `json:"product"`
		Breadcrumbs []database.GetCategoryPathByIDRow `json:"breadcrumbs"`
	}

	resp := response{
		Product:     productData,
		Breadcrumbs: breadcrumbs,
	}

	respondWithJSON(w, http.StatusOK, resp)
}
