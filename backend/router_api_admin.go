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
	mux.Handle("POST /api/admin/products/{productId}/variants", cfg.checkAuth(cfg.checkAdmin(http.HandlerFunc(cfg.handleApiAdminCreateVariant))))
	mux.Handle("GET /api/admin/products/{productId}/variants/{variantId}", cfg.checkAuth(cfg.checkAdmin(http.HandlerFunc(cfg.handleApiAdminGetVariant))))
	mux.Handle("PUT /api/admin/products/{productId}/variants/{variantId}", cfg.checkAuth(cfg.checkAdmin(http.HandlerFunc(cfg.handleApiAdminUpdateVariant))))
	mux.Handle("DELETE /api/admin/products/{productId}/variants/{variantId}", cfg.checkAuth(cfg.checkAdmin(http.HandlerFunc(cfg.handleApiAdminDeleteVariant))))
	mux.Handle("GET /api/admin/categories", cfg.checkAuth(cfg.checkAdmin(http.HandlerFunc(cfg.handleApiAdminGetCategories))))
	mux.Handle("GET /api/admin/categories/{categoryId}", cfg.checkAuth(cfg.checkAdmin(http.HandlerFunc(cfg.handleApiAdminGetCategory))))
	mux.Handle("PUT /api/admin/categories/{categoryId}", cfg.checkAuth(cfg.checkAdmin(http.HandlerFunc(cfg.handleApiAdminUpdateCategory))))
	mux.Handle("DELETE /api/admin/categories/{categoryId}", cfg.checkAuth(cfg.checkAdmin(http.HandlerFunc(cfg.handleApiAdminDeleteCategory))))
	mux.Handle("POST /api/admin/categories", cfg.checkAuth(cfg.checkAdmin(http.HandlerFunc(cfg.handleApiAdminCreateCategory))))

	mux.Handle("GET /api/admin/shipping-methods", cfg.checkAuth(cfg.checkAdmin(http.HandlerFunc(cfg.handleApiAdminGetShippingMethods))))
	mux.Handle("POST /api/admin/shipping-methods", cfg.checkAuth(cfg.checkAdmin(http.HandlerFunc(cfg.handleApiAdminCreateShippingMethod))))
	mux.Handle("GET /api/admin/shipping-methods/{shippingMethodId}", cfg.checkAuth(cfg.checkAdmin(http.HandlerFunc(cfg.handleApiAdminGetShippingMethod))))
	mux.Handle("PUT /api/admin/shipping-methods/{shippingMethodId}", cfg.checkAuth(cfg.checkAdmin(http.HandlerFunc(cfg.handleApiAdminUpdateShippingMethod))))
	log.Printf("Shop API routes registered")
}
