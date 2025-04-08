package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) registerShopRoutes(mux *http.ServeMux) {
	log.Printf("Registering Shop routes...")
	mux.Handle("GET /", cfg.maybeWithAuth(http.HandlerFunc(cfg.handleHomePage)))
	mux.Handle("GET /product/{slug}", cfg.maybeWithAuth(http.HandlerFunc(cfg.handleProductPage)))
	mux.Handle("GET /category/{slug}", cfg.maybeWithAuth(http.HandlerFunc(cfg.handleCategoryPage)))
	log.Printf("Shop routes registered")

}
