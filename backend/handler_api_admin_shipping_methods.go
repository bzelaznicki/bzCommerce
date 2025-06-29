package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bzelaznicki/bzCommerce/internal/database"
)

type ShippingMethodRequest struct {
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Price         float64 `json:"price"`
	EstimatedDays string  `json:"estimated_days"`
	SortOrder     int32   `json:"sort_order"`
	IsActive      bool    `json:"is_active"`
}

func (cfg *apiConfig) handleApiAdminGetShippingMethods(w http.ResponseWriter, r *http.Request) {
	shippingOptions, err := cfg.db.GetShippingOptions(r.Context())

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get shipping methods")
		return
	}

	respondWithJSON(w, http.StatusOK, shippingOptions)
}

func (cfg *apiConfig) handleApiAdminCreateShippingMethod(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	params := ShippingMethodRequest{}

	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	if params.Name == "" || params.Price < 0 || params.EstimatedDays == "" {
		respondWithError(w, http.StatusBadRequest, "Invalid parameters")
		return
	}

	shippingMethod, err := cfg.db.CreateShippingOption(r.Context(), database.CreateShippingOptionParams{
		Name:          params.Name,
		Description:   sql.NullString{String: params.Description, Valid: params.Description != ""},
		Price:         fmt.Sprintf("%2f", params.Price), //TODO [BZC-364]: refactor once the HTML frontend is removed to use float64 instead
		EstimatedDays: params.EstimatedDays,
		SortOrder:     params.SortOrder,
		IsActive:      params.IsActive,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to add shipping method")
		return
	}

	respondWithJSON(w, http.StatusOK, shippingMethod)

}
