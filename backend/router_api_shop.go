package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) registerApiShopRoutes(mux *http.ServeMux) {
	log.Printf("Registering Shop API routes...")
	mux.Handle("GET /api/products/{slug}", cfg.withCORS(http.HandlerFunc(cfg.handleApiGetSingleProduct)))
	mux.Handle("GET /api/products", cfg.withCORS(http.HandlerFunc(cfg.handleApiGetProducts)))
	mux.Handle("GET /api/categories", cfg.withCORS(http.HandlerFunc(cfg.handleApiGetCategories)))
	mux.Handle("GET /api/categories/{slug}/products", cfg.withCORS(http.HandlerFunc(cfg.handleApiGetCategoryProducts)))
	mux.Handle("POST /api/login", cfg.withCORS(http.HandlerFunc(cfg.handleApiLogin)))
	mux.Handle("POST /api/refresh", cfg.withCORS(http.HandlerFunc(cfg.handleApiRefreshToken)))
	mux.Handle("GET /api/account", cfg.checkAuth(http.HandlerFunc(cfg.handleApiGetAccount)))
	mux.Handle("POST /api/logout", cfg.withCORS(http.HandlerFunc(cfg.handleApiLogout)))
	log.Printf("Shop API routes registered")
}
