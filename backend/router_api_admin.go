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
	mux.Handle("DELETE /api/admin/shipping-methods/{shippingMethodId}", cfg.checkAuth(cfg.checkAdmin(http.HandlerFunc(cfg.handleApiAdminDeleteShippingMethod))))
	mux.Handle("PATCH /api/admin/shipping-methods/{shippingMethodId}/status", cfg.checkAuth(cfg.checkAdmin(http.HandlerFunc(cfg.handleApiAdminToggleShippingMethodStatus))))
	mux.Handle("GET /api/admin/users", cfg.checkAuth(cfg.checkAdmin(http.HandlerFunc(cfg.handleApiAdminGetUsers))))
	mux.Handle("GET /api/admin/users/{userId}", cfg.checkAuth(cfg.checkAdmin(http.HandlerFunc(cfg.handleApiAdminGetUserDetails))))
	mux.Handle("PUT /api/admin/users/{userId}", cfg.checkAuth(cfg.checkAdmin(http.HandlerFunc(cfg.handleApiAdminUpdateUserDetails))))
	mux.Handle("PATCH /api/admin/users/{userId}/password", cfg.checkAuth(cfg.checkAdmin(http.HandlerFunc(cfg.handleApiAdminUpdateUserPassword))))
	mux.Handle("DELETE /api/admin/users/{userId}", cfg.checkAuth(cfg.checkAdmin(http.HandlerFunc(cfg.handleApiAdminDeleteUser))))
	mux.Handle("POST /api/admin/users/{userId}/disable", cfg.checkAuth(cfg.checkAdmin(http.HandlerFunc(cfg.handleApiAdminDisableUser))))
	mux.Handle("POST /api/admin/users/{userId}/enable", cfg.checkAuth(cfg.checkAdmin(http.HandlerFunc(cfg.handleApiAdminEnableUser))))
	mux.Handle("GET /api/admin/countries", cfg.checkAuth(cfg.checkAdmin(http.HandlerFunc(cfg.handleApiAdminGetCountries))))
	mux.Handle("POST /api/admin/countries", cfg.checkAuth(cfg.checkAdmin(http.HandlerFunc(cfg.handleApiAdminCreateCountry))))
	mux.Handle("GET /api/admin/countries/{countryId}", cfg.checkAuth(cfg.checkAdmin(http.HandlerFunc(cfg.handleApiAdminGetCountry))))
	mux.Handle("PUT /api/admin/countries/{countryId}", cfg.checkAuth(cfg.checkAdmin(http.HandlerFunc(cfg.handleApiAdminUpdateCountry))))
	mux.Handle("PATCH /api/admin/countries/{countryId}/status", cfg.checkAuth(cfg.checkAdmin(http.HandlerFunc(cfg.handleApiAdminToggleCountryStatus))))
	mux.Handle("DELETE /api/admin/countries/{countryId}", cfg.checkAuth(cfg.checkAdmin(http.HandlerFunc(cfg.handleApiAdminDeleteCountry))))
	mux.Handle("GET /api/admin/orders", cfg.checkAuth(cfg.checkAdmin(http.HandlerFunc(cfg.handleApiAdminListOrders))))
	mux.Handle("GET /api/admin/orders/{orderId}", cfg.checkAuth(cfg.checkAdmin(http.HandlerFunc(cfg.handleApiAdminGetSingleOrder))))
	log.Printf("Shop API routes registered")
}
