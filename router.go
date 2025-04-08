package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) registerRoutes(mux *http.ServeMux) {
	cfg.registerStaticRoutes(mux)
	cfg.registerShopRoutes(mux)
	cfg.registerAdminRoutes(mux)
	cfg.registerAuthRoutes(mux)
	log.Printf("All routes registered")
}
