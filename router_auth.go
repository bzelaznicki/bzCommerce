package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) registerAuthRoutes(mux *http.ServeMux) {
	log.Printf("Registering Auth routes...")
	mux.Handle("GET /login", cfg.redirectIfAuthenticated(http.HandlerFunc(cfg.handleLoginGet)))
	mux.Handle("GET /register", cfg.redirectIfAuthenticated(http.HandlerFunc(cfg.handleRegisterGet)))
	mux.Handle("POST /login", cfg.redirectIfAuthenticated(http.HandlerFunc(cfg.handleLoginPost)))
	mux.Handle("POST /register", cfg.redirectIfAuthenticated(http.HandlerFunc(cfg.handleRegisterPost)))
	mux.Handle("GET /account", cfg.withAuth(http.HandlerFunc(cfg.handleAccountPage)))
	mux.Handle("GET /account/addresses", cfg.withAuth(http.HandlerFunc(cfg.handleManageAddressesList)))
	mux.HandleFunc("GET /logout", cfg.handleLogout)
	log.Printf("Auth routes registered")
}
