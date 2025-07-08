package main

import (
	"net/http"

	"github.com/bzelaznicki/bzCommerce/internal/database"
)

func (cfg *apiConfig) handleApiGetShippingMethods(w http.ResponseWriter, r *http.Request) {
	shippingMethods, err := cfg.db.GetActiveShippingOptions(r.Context())

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get shipping methods")
		return
	}

	if shippingMethods == nil {
		shippingMethods = []database.ShippingOption{}
	}

	respondWithJSON(w, http.StatusOK, shippingMethods)
}
