package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) registerApiShopRoutes(mux *http.ServeMux) {
	log.Printf("Registering Shop API routes...")
	mux.Handle("GET /api/products/{slug}", withCORS(http.HandlerFunc(cfg.handleApiGetSingleProduct)))
	mux.Handle("GET /api/products", withCORS(http.HandlerFunc(cfg.handleApiGetProducts)))
	mux.Handle("GET /api/categories", withCORS(http.HandlerFunc(cfg.handleApiGetCategories)))
	mux.Handle("GET /api/categories/{slug}/products", withCORS(http.HandlerFunc(cfg.handleApiGetCategoryProducts)))
	mux.Handle("POST /api/login", withCORS(http.HandlerFunc(cfg.handleApiLogin)))
	log.Printf("Shop API routes registered")
}
