package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/bzelaznicki/bzCommerce/internal/database"
	"github.com/google/uuid"
)

type CheckoutPageData struct {
	Items           []database.GetCartDetailsWithSnapshotPriceRow
	Total           float64
	ShippingOptions []database.ShippingOption
}

func (cfg *apiConfig) handleViewCheckout(w http.ResponseWriter, r *http.Request) {
	cartID, err := cfg.getOrCreateCartID(w, r)

	if err != nil {
		cfg.Render(w, r, "templates/pages/cart.html", struct {
			Items []database.GetCartDetailsWithSnapshotPriceRow
			Total float64
		}{})
		return
	}

	items, err := cfg.db.GetCartDetailsWithSnapshotPrice(r.Context(), cartID)
	if err != nil {
		cfg.RenderError(w, r, http.StatusInternalServerError, "Could not load cart")
		log.Printf("Error loading cart: %v", err)
		return
	}

	var total float64
	for _, item := range items {
		total += float64(item.Quantity) * item.PricePerItem
	}

	shippingOptions, err := cfg.db.GetShippingOptions(r.Context())
	if err != nil {
		cfg.RenderError(w, r, http.StatusInternalServerError, "Could not load shipping options")
		log.Printf("Error loading shipping options: %v", err)
		return
	}

	data := CheckoutPageData{
		Items:           items,
		Total:           total,
		ShippingOptions: shippingOptions,
	}
	cfg.Render(w, r, "templates/pages/checkout.html", data)

}

func (cfg *apiConfig) handleCheckout(w http.ResponseWriter, r *http.Request) {
	cartId, err := cfg.getOrCreateCartID(w, r)
	if err != nil {
		cfg.RenderError(w, r, http.StatusInternalServerError, "Could not load cart")
		log.Printf("Error loading cart: %v", err)
		return
	}

	customerEmail := r.FormValue("customer_email")
	shippingName := r.FormValue("shipping_name")
	shippingAddress := r.FormValue("shipping_address")
	shippingCity := r.FormValue("shipping_city")
	shippingPostalCode := r.FormValue("shipping_postal_code")
	shippingCountry := r.FormValue("shipping_country")
	shippingPhone := r.FormValue("shipping_phone")
	shippingOptionID := r.FormValue("shipping_option_id")

	sameAsShipping := r.PostFormValue("same_as_shipping") != ""

	billingName := r.PostFormValue("billing_name")
	billingAddress := r.PostFormValue("billing_address")
	billingCity := r.PostFormValue("billing_city")
	billingPostalCode := r.PostFormValue("billing_postal_code")
	billingCountry := r.PostFormValue("billing_country")

	if sameAsShipping {
		billingName = r.PostFormValue("shipping_name")
		billingAddress = r.PostFormValue("shipping_address")
		billingCity = r.PostFormValue("shipping_city")
		billingPostalCode = r.PostFormValue("shipping_postal_code")
		billingCountry = r.PostFormValue("shipping_country")
	}

	cart, err := cfg.db.GetCartById(r.Context(), cartId)
	if err != nil {
		cfg.RenderError(w, r, http.StatusInternalServerError, "Could not load cart")
		log.Printf("Error loading cart: %v", err)
		return
	}

	_, err = cfg.db.CreateOrder(r.Context(), database.CreateOrderParams{
		UserID:        cart.UserID,
		CustomerEmail: customerEmail,
		ShippingName: sql.NullString{
			String: shippingName,
			Valid:  shippingName != "",
		},
		ShippingAddress: sql.NullString{
			String: shippingAddress,
			Valid:  shippingAddress != "",
		},
		ShippingCity: sql.NullString{
			String: shippingCity,
			Valid:  shippingCity != "",
		},
		ShippingPostalCode: sql.NullString{
			String: shippingPostalCode,
			Valid:  shippingPostalCode != "",
		},

		ShippingCountry: sql.NullString{
			String: shippingCountry,
			Valid:  shippingCountry != "",
		},
		ShippingPhone: shippingPhone,
		ShippingOptionID: uuid.NullUUID{
			UUID:  uuid.MustParse(shippingOptionID),
			Valid: shippingOptionID != "",
		},

		BillingName: sql.NullString{
			String: billingName,
			Valid:  billingName != "",
		},
		BillingAddress: sql.NullString{
			String: billingAddress,
			Valid:  billingAddress != "",
		},
		BillingCity: sql.NullString{
			String: billingCity,
			Valid:  billingCity != "",
		},
		BillingPostalCode: sql.NullString{
			String: billingPostalCode,
			Valid:  billingPostalCode != "",
		},
		BillingCountry: sql.NullString{
			String: billingCountry,
			Valid:  billingCountry != "",
		},
	},
	)
	if err != nil {
		cfg.RenderError(w, r, http.StatusInternalServerError, "Could not create order")
		log.Printf("Error creating order: %v", err)
		return
	}

}
