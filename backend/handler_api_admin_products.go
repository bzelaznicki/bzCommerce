package main

import "net/http"

func (cfg *apiConfig) handleApiAdminGetProducts(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, struct {
		Works bool `json:"works"`
	}{
		Works: true,
	})
}
