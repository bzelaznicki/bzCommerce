package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/bzelaznicki/bzCommerce/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type VariantRequest struct {
	Sku           string  `json:"sku"`
	Price         float64 `json:"price"`
	StockQuantity int32   `json:"stock_quantity"`
	ImageUrl      string  `json:"image_url"`
	Name          string  `json:"name"`
}

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

func (cfg *apiConfig) handleApiAdminCreateVariant(w http.ResponseWriter, r *http.Request) {
	productId, err := uuid.Parse(r.PathValue("productId"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}
	product, err := cfg.db.GetProductById(r.Context(), productId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Product not found")
		return
	}

	decoder := json.NewDecoder(r.Body)

	params := VariantRequest{}

	err = decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if params.Name == "" || params.Sku == "" || params.Price <= 0 || params.StockQuantity < 0 {
		respondWithError(w, http.StatusBadRequest, "Invalid variant data")
		return
	}

	variant := database.CreateVariantParams{
		ProductID:     product.ID,
		Sku:           params.Sku,
		Price:         params.Price,
		StockQuantity: params.StockQuantity,
		ImageUrl:      sql.NullString{String: params.ImageUrl, Valid: params.ImageUrl != ""},
		VariantName:   sql.NullString{String: params.Name, Valid: params.Name != ""},
	}

	addedVariant, err := cfg.db.CreateVariant(r.Context(), variant)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			if pqErr.Constraint == "products_slug_key" {
				respondWithError(w, http.StatusConflict, "Variant with the same SKU already exists")
			} else {
				respondWithError(w, http.StatusConflict, "Variant already exists (duplicate field)")
			}
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to create variant")
		return
	}

	respondWithJSON(w, http.StatusOK, addedVariant)
}

func (cfg *apiConfig) handlerApiAdminGetVariant(w http.ResponseWriter, r *http.Request) {
	variantId, err := uuid.Parse(r.PathValue("variantId"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid variant ID")
		return
	}

	variant, err := cfg.db.GetVariantByID(r.Context(), variantId)

	if err != nil {
		respondWithError(w, http.StatusNotFound, "Variant not found")
		return
	}

	respondWithJSON(w, http.StatusOK, variant)

}

func (cfg *apiConfig) handlerApiAdminUpdateVariant(w http.ResponseWriter, r *http.Request) {
	variantId, err := uuid.Parse(r.PathValue("variantId"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid variant ID")
		return
	}

	decoder := json.NewDecoder(r.Body)

	params := VariantRequest{}

	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	variant := database.UpdateVariantParams{
		Sku:           params.Sku,
		Price:         params.Price,
		StockQuantity: params.StockQuantity,
		VariantName:   sql.NullString{String: params.Name, Valid: params.Name != ""},
		ImageUrl:      sql.NullString{String: params.ImageUrl, Valid: params.ImageUrl != ""},
		ID:            variantId,
	}

	updatedVariant, err := cfg.db.UpdateVariant(r.Context(), variant)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Variant not found")
			return
		}

		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			switch pqErr.Constraint {
			case "products_slug_key":
				respondWithError(w, http.StatusConflict, "Variant with the same SKU already exists")
			case "products_sku_key":
				respondWithError(w, http.StatusConflict, "Variant with the same SKU already exists")
			default:
				respondWithError(w, http.StatusConflict, "Variant already exists (duplicate field)")
			}
			return
		}

		respondWithError(w, http.StatusInternalServerError, "Failed to update variant")
		return
	}

	respondWithJSON(w, http.StatusOK, updatedVariant)
}
