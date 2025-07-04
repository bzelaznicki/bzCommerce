package main

import (
	"encoding/json"
	"net/http"

	"github.com/bzelaznicki/bzCommerce/internal/database"
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
