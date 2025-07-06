package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/bzelaznicki/bzCommerce/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func (cfg *apiConfig) handleApiAdminGetCountries(w http.ResponseWriter, r *http.Request) {
	countries, err := cfg.db.GetCountries(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get countries")
		return
	}

	if countries == nil {
		countries = []database.Country{}
	}

	respondWithJSON(w, http.StatusOK, countries)
}

func (cfg *apiConfig) handleApiAdminCreateCountry(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	params := database.CreateCountryParams{}

	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if params.IsoCode == "" || params.Name == "" {
		respondWithError(w, http.StatusBadRequest, "ISO code and name cannot be empty")
		return
	}
	params.IsoCode = strings.ToUpper(params.IsoCode)
	createdCountry, err := cfg.db.CreateCountry(r.Context(), params)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			if pqErr.Constraint == "countries_iso_code_key" {
				respondWithError(w, http.StatusConflict, "Country with the same ISO code already exists")
			} else {
				respondWithError(w, http.StatusConflict, "Country already exists (duplicate field)")
			}
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to create country")
		return
	}

	respondWithJSON(w, http.StatusOK, createdCountry)
}

func (cfg *apiConfig) handleApiAdminGetCountry(w http.ResponseWriter, r *http.Request) {
	countryId, err := uuid.Parse(r.PathValue("countryId"))

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid country ID")
		return
	}

	country, err := cfg.db.GetCountryById(r.Context(), countryId)

	if err != nil {
		respondWithError(w, http.StatusNotFound, "Country not found")
		return
	}

	respondWithJSON(w, http.StatusOK, country)
}

func (cfg *apiConfig) handleApiAdminUpdateCountry(w http.ResponseWriter, r *http.Request) {
	countryId, err := uuid.Parse(r.PathValue("countryId"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid country ID")
		return
	}

	decoder := json.NewDecoder(r.Body)

	params := database.UpdateCountryByIdParams{}

	err = decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}
	if params.IsoCode == "" || params.Name == "" {
		respondWithError(w, http.StatusBadRequest, "ISO code and name cannot be empty")
		return
	}
	params.ID = countryId
	params.IsoCode = strings.ToUpper(params.IsoCode)

	country, err := cfg.db.UpdateCountryById(r.Context(), params)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			if pqErr.Constraint == "countries_iso_code_key" {
				respondWithError(w, http.StatusConflict, "Country with the same ISO code already exists")
			} else {
				respondWithError(w, http.StatusConflict, "Country already exists (duplicate field)")
			}
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to update country")
		return
	}

	respondWithJSON(w, http.StatusOK, country)
}

func (cfg *apiConfig) handleApiAdminToggleCountryStatus(w http.ResponseWriter, r *http.Request) {
	countryId, err := uuid.Parse(r.PathValue("countryId"))

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid country ID")
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

	row, err := cfg.db.ToggleCountryStatus(r.Context(), database.ToggleCountryStatusParams{
		ID:       countryId,
		IsActive: params.Status,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Country not found")
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

func (cfg *apiConfig) handleApiAdminDeleteCountry(w http.ResponseWriter, r *http.Request) {
	countryId, err := uuid.Parse(r.PathValue("countryId"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid country ID")
		return
	}

	rows, err := cfg.db.DeleteCountryById(r.Context(), countryId)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23503" {
				respondWithError(w, http.StatusConflict, "Cannot delete country - in use")
				return
			}
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to delete country")
		return
	}

	if rows == 0 {
		respondWithError(w, http.StatusNotFound, "Country not found")
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
