package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/bzelaznicki/bzCommerce/internal/database"
	"github.com/google/uuid"
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
		IsActive:      params.IsActive,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to add shipping method")
		return
	}

	respondWithJSON(w, http.StatusOK, shippingMethod)

}

func (cfg *apiConfig) handleApiAdminGetShippingMethod(w http.ResponseWriter, r *http.Request) {
	shippingMethodId, err := uuid.Parse(r.PathValue("shippingMethodId"))

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid shipping method ID")
		return
	}

	shippingMethod, err := cfg.db.SelectShippingOptionById(r.Context(), shippingMethodId)

	if err != nil {
		respondWithError(w, http.StatusNotFound, "Shipping method not found")
		return
	}

	respondWithJSON(w, http.StatusOK, shippingMethod)
}

func (cfg *apiConfig) handleApiAdminUpdateShippingMethod(w http.ResponseWriter, r *http.Request) {
	shippingMethodId, err := uuid.Parse(r.PathValue("shippingMethodId"))

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid shipping method ID")
		return
	}

	decoder := json.NewDecoder(r.Body)

	params := ShippingMethodRequest{}

	err = decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	if params.Name == "" || params.Price < 0 || params.EstimatedDays == "" {
		respondWithError(w, http.StatusBadRequest, "Invalid parameters")
		return
	}

	shippingMethod, err := cfg.db.UpdateShippingOption(r.Context(), database.UpdateShippingOptionParams{
		ID:            shippingMethodId,
		Name:          params.Name,
		Description:   sql.NullString{String: params.Description, Valid: params.Description != ""},
		Price:         fmt.Sprintf("%2f", params.Price), //TODO [BZC-364]: refactor once the HTML frontend is removed to use float64 instead
		EstimatedDays: params.EstimatedDays,
		SortOrder:     params.SortOrder,
		IsActive:      params.IsActive,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update shipping method")
		return
	}

	respondWithJSON(w, http.StatusOK, shippingMethod)
}

func (cfg *apiConfig) handleApiAdminToggleShippingMethodStatus(w http.ResponseWriter, r *http.Request) {
	shippingMethodId, err := uuid.Parse(r.PathValue("shippingMethodId"))

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid shipping method ID")
		return
	}
	decoder := json.NewDecoder(r.Body)
	params := struct {
		Status bool `json:"status"`
	}{}

	err = decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	row, err := cfg.db.ToggleShippingOptionStatus(r.Context(), database.ToggleShippingOptionStatusParams{
		IsActive: params.Status,
		ID:       shippingMethodId,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Shipping method not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to update status")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"id":      row.ID,
		"status":  row.IsActive,
	})

}
