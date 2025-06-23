package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) registerApiAdminRoutes(mux *http.ServeMux) {
	log.Printf("Registering Shop API routes...")
	mux.Handle("GET /api/admin/products", cfg.checkAuth(cfg.checkAdmin(http.HandlerFunc(cfg.handleApiAdminGetProducts))))
	mux.Handle("POST /api/admin/products", cfg.checkAuth(cfg.checkAdmin(http.HandlerFunc(cfg.handleApiCreateProduct))))
	mux.Handle("POST /api/admin/cloudinary", cfg.checkAuth(cfg.checkAdmin(http.HandlerFunc(cfg.handleCloudinarySignUpload))))
	mux.Handle("PUT /api/admin/products/{id}", cfg.checkAuth(cfg.checkAdmin(http.HandlerFunc(cfg.handleApiUpdateProduct))))
	mux.Handle("GET /api/admin/products/{id}", cfg.checkAuth(cfg.checkAdmin(http.HandlerFunc(cfg.handleApiAdminGetProductDetails))))
	mux.Handle("DELETE /api/admin/products/{id}", cfg.checkAuth(cfg.checkAdmin(http.HandlerFunc(cfg.handleApiAdminDeleteProduct))))
	mux.Handle("GET /api/admin/products/{id}/variants", cfg.checkAuth(cfg.checkAdmin(http.HandlerFunc(cfg.handleApiAdminGetVariants))))
	log.Printf("Shop API routes registered")
}
