package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) registerApiAuthRoutes(mux *http.ServeMux) {
	log.Printf("Registering Auth API routes...")
	mux.Handle("POST /api/login", http.HandlerFunc(cfg.handleApiLogin))
	mux.Handle("POST /api/refresh", http.HandlerFunc(cfg.handleApiRefreshToken))
	mux.Handle("GET /api/account", cfg.checkAuth(http.HandlerFunc(cfg.handleApiGetAccount)))
	mux.Handle("POST /api/logout", http.HandlerFunc(cfg.handleApiLogout))
	mux.Handle("POST /api/users", http.HandlerFunc(cfg.handlerApiRegister))
	log.Printf("Auth API routes registered")
}
