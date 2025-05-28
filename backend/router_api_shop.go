package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) registerApiShopRoutes(mux *http.ServeMux) {
	log.Printf("Registering Shop API routes...")
	mux.Handle("GET /api/products/{slug}", http.HandlerFunc(cfg.handleApiGetSingleProduct))
	log.Printf("Shop API routes registered")
}
