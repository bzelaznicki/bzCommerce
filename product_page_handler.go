package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) handleProductPage(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if slug == "" {
		http.NotFound(w, r)
		return
	}

	product, err := cfg.db.GetProductBySlug(r.Context(), slug)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	productData := Product{
		ID:          product.ID,
		Name:        product.Name,
		Slug:        product.Slug,
		ImagePath:   product.ImageUrl.String,
		CategoryID:  product.CategoryID,
		Description: product.Description.String,
	}

	data := struct {
		StoreName string
		Product   Product
	}{
		StoreName: cfg.storeName,
		Product:   productData,
	}

	tmpl := parsePageTemplate("templates/pages/product.html")

	err = tmpl.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		log.Printf("template error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

}
