package main

import (
	"net/http"

	"github.com/bzelaznicki/bzCommerce/internal/database"
)

func (cfg *apiConfig) handleApiGetPaymentMethods(w http.ResponseWriter, r *http.Request) {
	paymentMethods, err := cfg.db.GetActivePaymentOptions(r.Context())

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get payment methods")
		return
	}
	if paymentMethods == nil {
		paymentMethods = []database.PaymentOption{}
	}

	respondWithJSON(w, http.StatusOK, paymentMethods)
}
