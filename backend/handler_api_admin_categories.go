package main

import "net/http"

func (cfg *apiConfig) handleApiAdminGetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := cfg.db.ListCategoriesWithParent(r.Context())

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to load categories")
		return
	}

	respondWithJSON(w, http.StatusOK, categories)
}
