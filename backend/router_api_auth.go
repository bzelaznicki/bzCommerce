package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) registerApiAuthRoutes(mux *http.ServeMux) {
	log.Printf("Registering Auth API routes...")
	
	// Apply rate limiting to authentication endpoints
	authLimiter := cfg.withRateLimit(cfg.authRateLimiter)
	
	mux.Handle("POST /api/login", authLimiter(http.HandlerFunc(cfg.handleApiLogin)))
	mux.Handle("POST /api/refresh", authLimiter(http.HandlerFunc(cfg.handleApiRefreshToken)))
	mux.Handle("GET /api/account", cfg.checkAuth(http.HandlerFunc(cfg.handleApiGetAccount)))
	mux.Handle("POST /api/logout", http.HandlerFunc(cfg.handleApiLogout))
	mux.Handle("POST /api/users", authLimiter(http.HandlerFunc(cfg.handlerApiRegister)))
	log.Printf("Auth API routes registered")
}
