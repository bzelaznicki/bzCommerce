package main

import (
	"net/http"
)

func (cfg *apiConfig) handleHomePage(w http.ResponseWriter, r *http.Request) {
	dbProducts, err := cfg.db.ListProducts(r.Context())
	if err != nil {
		cfg.RenderError(w, r, http.StatusInternalServerError, "Internal Server Error")

		return
	}

	products := make([]Product, 0, len(dbProducts))

	for _, dbProduct := range dbProducts {
		product := Product{
			ID:          dbProduct.ID,
			Name:        dbProduct.Name,
			Slug:        dbProduct.Slug,
			ImagePath:   dbProduct.ImageUrl.String,
			CategoryID:  dbProduct.CategoryID,
			Description: dbProduct.Description.String,
		}
		products = append(products, product)
	}
	cfg.Render(w, r, "templates/pages/home.html", products)
}
