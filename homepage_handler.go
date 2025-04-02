package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) handleHomePage(w http.ResponseWriter, r *http.Request) {
	dbProducts, err := cfg.db.ListProducts(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	data := struct {
		StoreName string
		Products  []Product
	}{
		StoreName: cfg.storeName,
		Products:  products,
	}

	tmpl := parsePageTemplate("templates/pages/home.html")

	err = tmpl.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		log.Printf("template error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

}
