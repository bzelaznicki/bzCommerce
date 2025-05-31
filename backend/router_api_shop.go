package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) registerApiShopRoutes(mux *http.ServeMux) {
	log.Printf("Registering Shop API routes...")
	mux.Handle("GET /api/products/{slug}", http.HandlerFunc(cfg.handleApiGetSingleProduct))
	mux.Handle("GET /api/products", http.HandlerFunc(cfg.handleApiGetProducts))
	mux.Handle("GET /api/categories/{slug}/products", http.HandlerFunc(cfg.handleApiGetCategoryProducts))
	log.Printf("Shop API routes registered")
}
