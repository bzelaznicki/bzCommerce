package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) registerApiAdminRoutes(mux *http.ServeMux) {
	log.Printf("Registering Shop API routes...")
	mux.Handle("GET /api/admin/products", cfg.checkAuth(cfg.checkAdmin(http.HandlerFunc(cfg.handleApiAdminGetProducts))))
	log.Printf("Shop API routes registered")
}
