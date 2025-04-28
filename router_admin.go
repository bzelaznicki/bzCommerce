package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) registerAdminRoutes(mux *http.ServeMux) {
	log.Printf("Registering Admin routes...")
	//ADMIN
	mux.Handle("GET /admin", cfg.requireAdmin(http.HandlerFunc(cfg.handleAdminDashboard)))
	// View all products
	mux.Handle("GET /admin/products", cfg.requireAdmin(http.HandlerFunc(cfg.handleAdminProductList)))

	// New product form
	mux.Handle("GET /admin/products/new", cfg.requireAdmin(http.HandlerFunc(cfg.handleAdminProductNewForm)))
	mux.Handle("POST /admin/products", cfg.requireAdmin(http.HandlerFunc(cfg.handleAdminProductCreate)))

	// Edit product
	mux.Handle("GET /admin/products/{id}/edit", cfg.requireAdmin(http.HandlerFunc(cfg.handleAdminProductEditForm)))
	mux.Handle("POST /admin/products/{id}", cfg.requireAdmin(http.HandlerFunc(cfg.handleAdminProductUpdate)))

	// Delete product
	mux.Handle("POST /admin/products/{id}/delete", cfg.requireAdmin(http.HandlerFunc(cfg.handleAdminProductDelete)))

	//List product variants
	mux.Handle("GET /admin/products/{id}/variants", cfg.requireAdmin(http.HandlerFunc(cfg.handleAdminVariantList)))

	//New product form
	mux.Handle("GET /admin/products/{id}/variants/new", cfg.requireAdmin(http.HandlerFunc(cfg.handleAdminVariantNewForm)))
	mux.Handle("POST /admin/products/{id}/variants", cfg.requireAdmin(http.HandlerFunc(cfg.handleAdminVariantCreate)))

	//Edit Variants
	mux.Handle("GET /admin/variants/{id}/edit", cfg.requireAdmin(http.HandlerFunc(cfg.handleAdminVariantEditForm)))
	mux.Handle("POST /admin/variants/{id}", cfg.requireAdmin(http.HandlerFunc(cfg.handleAdminVariantUpdate)))

	//Delete variant
	mux.Handle("POST /admin/variants/{id}/delete", cfg.requireAdmin(http.HandlerFunc(cfg.handleAdminVariantDelete)))

	// View all categories
	mux.Handle("GET /admin/categories", cfg.requireAdmin(http.HandlerFunc(cfg.handleAdminCategoryList)))

	// New category form + create
	mux.Handle("GET /admin/categories/new", cfg.requireAdmin(http.HandlerFunc(cfg.handleAdminCategoryNewForm)))

	mux.Handle("POST /admin/categories", cfg.requireAdmin(http.HandlerFunc(cfg.handleAdminCategoryCreate)))

	// Edit category form + update
	mux.Handle("GET /admin/categories/{id}/edit", cfg.requireAdmin(http.HandlerFunc(cfg.handleAdminCategoryEditForm)))
	mux.Handle("POST /admin/categories/{id}", cfg.requireAdmin(http.HandlerFunc(cfg.handleAdminCategoryUpdate)))

	// Delete category
	mux.Handle("POST /admin/categories/{id}/delete", cfg.requireAdmin(http.HandlerFunc(cfg.handleAdminCategoryDelete)))

	// View users
	mux.Handle("GET /admin/users", cfg.requireAdmin(http.HandlerFunc(cfg.handleAdminUsersList)))
	// Edit users page
	mux.Handle("GET /admin/users/{id}/edit", cfg.requireAdmin(http.HandlerFunc(cfg.handleAdminUserEditForm)))
	mux.Handle("POST /admin/users/{id}", cfg.requireAdmin(http.HandlerFunc(cfg.handleAdminUserUpdate)))

	// Delete user
	mux.Handle("POST /admin/users/{id}/delete", cfg.requireAdmin(http.HandlerFunc(cfg.handleAdminUserDelete)))

	// View shipping options
	mux.Handle("GET /admin/shipping", cfg.requireAdmin(http.HandlerFunc(cfg.handleAdminShippingOptionsList)))

	// New shipping option form
	mux.Handle("GET /admin/shipping/new", cfg.requireAdmin(http.HandlerFunc(cfg.handleAdminShippingOptionsNew)))

	// Create shipping option
	mux.Handle("POST /admin/shipping/new", cfg.requireAdmin(http.HandlerFunc(cfg.handleAdminShippingOptionCreate)))

	// Edit shipping option form
	mux.Handle("GET /admin/shipping/{id}/edit", cfg.requireAdmin(http.HandlerFunc(cfg.handleAdminShippingOptionsEdit)))
	// Update shipping option
	mux.Handle("POST /admin/shipping/{id}/edit", cfg.requireAdmin(http.HandlerFunc(cfg.handleAdminShippingOptionUpdate)))
	// Delete shipping option
	mux.Handle("POST /admin/shipping/{id}/delete", cfg.requireAdmin(http.HandlerFunc(cfg.handleAdminShippingOptionDelete)))

	// Payment options list
	mux.Handle("GET /admin/payment", cfg.requireAdmin(http.HandlerFunc(cfg.handleAdminPaymentOptionsList)))

	// New payment option form
	mux.Handle("GET /admin/payment/new", cfg.requireAdmin(http.HandlerFunc(cfg.handleAdminPaymentOptionsNew)))

	// Create payment option
	mux.Handle("POST /admin/payment/new", cfg.requireAdmin(http.HandlerFunc(cfg.handleAdminPaymentOptionsCreate)))

	// Edit payment option form
	mux.Handle("GET /admin/payment/{id}/edit", cfg.requireAdmin(http.HandlerFunc(cfg.handleAdminPaymentOptionEdit)))

	// Update payment option
	mux.Handle("POST /admin/payment/{id}/edit", cfg.requireAdmin(http.HandlerFunc(cfg.handleAdminPaymentOptionUpdate)))

	// Delete payment option
	mux.Handle("POST /admin/payment/{id}/delete", cfg.requireAdmin(http.HandlerFunc(cfg.handleAdminPaymentOptionDelete)))

	log.Printf("Admin routes registered")
}
