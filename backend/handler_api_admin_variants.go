package main

import (
	"net/http"

	"github.com/bzelaznicki/bzCommerce/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handleApiAdminGetVariants(w http.ResponseWriter, r *http.Request) {
	productId, err := uuid.Parse(r.PathValue("id"))

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	product, err := cfg.db.GetProductById(r.Context(), productId)

	if err != nil {
		respondWithError(w, http.StatusNotFound, "Product not foundf")
		return
	}

	variants, err := cfg.db.GetVariantsByProductID(r.Context(), product.ID)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get variants")
		return
	}

	if variants == nil {
		variants = []database.ProductVariant{}
	}

	type response struct {
		ProductId   uuid.UUID                 `json:"product_id"`
		ProductName string                    `json:"product_name"`
		Variants    []database.ProductVariant `json:"product_variants"`
	}

	resp := response{
		ProductId:   product.ID,
		ProductName: product.Name,
		Variants:    variants,
	}

	respondWithJSON(w, http.StatusOK, resp)
}
