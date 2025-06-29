package main

import "net/http"

func (cfg *apiConfig) handleApiAdminGetShippingMethods(w http.ResponseWriter, r *http.Request) {
	shippingOptions, err := cfg.db.GetShippingOptions(r.Context())

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get shipping methods")
		return
	}

	respondWithJSON(w, http.StatusOK, shippingOptions)
}
