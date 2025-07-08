package main

import (
	"log"
	"net/http"

	"github.com/bzelaznicki/bzCommerce/internal/database"
	"github.com/google/uuid"
)

type CheckoutPageData struct {
	Items           []database.GetCartDetailsWithSnapshotPriceRow
	Total           float64
	ShippingOptions []database.ShippingOption
	PaymentOptions  []database.PaymentOption
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

	paymentOptions, err := cfg.db.GetPaymentOptions(r.Context())

	if err != nil {
		cfg.RenderError(w, r, http.StatusBadRequest, "Could not load payment options")
		log.Printf("Error loading payment options: %v", err)
	}

	data := CheckoutPageData{
		Items:           items,
		Total:           total,
		ShippingOptions: shippingOptions,
		PaymentOptions:  paymentOptions,
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
	shippingPhone := r.FormValue("shipping_phone")
	shippingOptionID := r.FormValue("shipping_method_id")
	paymentOptionID := r.FormValue("payment_method_id")

	sameAsShipping := r.PostFormValue("same_as_shipping") != ""

	billingName := r.PostFormValue("billing_name")
	billingAddress := r.PostFormValue("billing_address")
	billingCity := r.PostFormValue("billing_city")
	billingPostalCode := r.PostFormValue("billing_postal_code")

	if sameAsShipping {
		billingName = r.PostFormValue("shipping_name")
		billingAddress = r.PostFormValue("shipping_address")
		billingCity = r.PostFormValue("shipping_city")
		billingPostalCode = r.PostFormValue("shipping_postal_code")

	}

	cart, err := cfg.db.GetCartById(r.Context(), cartId)
	if err != nil {
		cfg.RenderError(w, r, http.StatusInternalServerError, "Could not load cart")
		log.Printf("Error loading cart: %v", err)
		return
	}

	items, err := cfg.db.GetCartDetailsWithSnapshotPrice(r.Context(), cartId)
	if err != nil {
		cfg.RenderError(w, r, http.StatusInternalServerError, "Could not load cart items")
		log.Printf("Error loading cart items: %v", err)
		return
	}

	var subtotal float64
	for _, item := range items {
		subtotal += float64(item.Quantity) * item.PricePerItem
	}

	var shippingPrice float64
	if shippingOptionID != "" {
		shippingUUID, err := uuid.Parse(shippingOptionID)
		if err != nil {
			cfg.RenderError(w, r, http.StatusBadRequest, "Invalid shipping method selected")
			return
		}

		option, err := cfg.db.SelectShippingOptionById(r.Context(), shippingUUID)
		if err != nil {
			cfg.RenderError(w, r, http.StatusInternalServerError, "Could not find shipping option")
			log.Printf("Error finding shipping option: %v", err)
			return
		}

		shippingPrice = option.Price
	}

	totalPrice := subtotal + shippingPrice

	_, err = cfg.db.CreateOrder(r.Context(), database.CreateOrderParams{
		UserID:             cart.UserID,
		CustomerEmail:      customerEmail,
		TotalPrice:         totalPrice,
		ShippingPrice:      shippingPrice,
		ShippingName:       shippingName,
		ShippingAddress:    shippingAddress,
		ShippingCity:       shippingCity,
		ShippingPostalCode: shippingPostalCode,
		ShippingCountryID:  uuid.UUID{},
		ShippingPhone:      shippingPhone,
		ShippingOptionID:   uuid.MustParse(shippingOptionID),
		PaymentOptionID:    uuid.MustParse(paymentOptionID),
		BillingName:        billingName,
		BillingAddress:     billingAddress,
		BillingCity:        billingCity,
		BillingPostalCode:  billingPostalCode,
		BillingCountryID:   uuid.UUID{},
	})
	if err != nil {
		cfg.RenderError(w, r, http.StatusInternalServerError, "Could not create order")
		log.Printf("Error creating order: %v", err)
		return
	}
	http.Redirect(w, r, "/checkout/payment", http.StatusSeeOther)
}
