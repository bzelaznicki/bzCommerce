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
	mux.Handle("POST /cart/add", cfg.maybeWithAuth(http.HandlerFunc(cfg.handleAddToCart)))
	mux.Handle("GET /cart", cfg.maybeWithAuth(http.HandlerFunc(cfg.handleViewCart)))
	mux.Handle("GET /cart/remove/{id}", cfg.maybeWithAuth(http.HandlerFunc(cfg.handleDeleteCartItem)))
	mux.Handle("POST /cart/update", cfg.maybeWithAuth(http.HandlerFunc(cfg.handleCartUpdate)))
	mux.Handle("GET /checkout", cfg.maybeWithAuth(http.HandlerFunc(cfg.handleViewCheckout)))
	mux.Handle("POST /checkout", cfg.maybeWithAuth(http.HandlerFunc(cfg.handleCheckout)))
	log.Printf("Shop routes registered")

}
