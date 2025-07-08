package main

import (
	"net/http"

	"github.com/bzelaznicki/bzCommerce/internal/database"
)

func (cfg *apiConfig) handleApiGetCountries(w http.ResponseWriter, r *http.Request) {
	countries, err := cfg.db.GetCountries(r.Context())

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting countries")
		return
	}

	if countries == nil {
		countries = []database.Country{}
	}

	respondWithJSON(w, http.StatusOK, countries)
}
