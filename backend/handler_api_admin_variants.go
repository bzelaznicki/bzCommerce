package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/bzelaznicki/bzCommerce/internal/database"
	"github.com/google/uuid"
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
	decoder := json.NewDecoder(r.Body)

	params := VariantRequest{}

	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if params.Name == "" || params.Sku == "" || params.Price <= 0 || params.StockQuantity < 0 {
		respondWithError(w, http.StatusBadRequest, "Invalid variant data")
		return
	}

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
		respondWithError(w, http.StatusInternalServerError, "Failed to create variant")
		return
	}

	respondWithJSON(w, http.StatusOK, addedVariant)
}
