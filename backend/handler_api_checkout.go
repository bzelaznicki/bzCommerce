package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/bzelaznicki/bzCommerce/internal/database"
	"github.com/google/uuid"
)

type CheckoutParams struct {
	CustomerEmail      string    `json:"customer_email"`
	ShippingName       string    `json:"shipping_name"`
	ShippingAddress    string    `json:"shipping_address"`
	ShippingCity       string    `json:"shipping_city"`
	ShippingPostalCode string    `json:"shipping_postal_code"`
	ShippingCountryID  uuid.UUID `json:"shipping_country_id"`
	ShippingPhone      string    `json:"shipping_phone"`
	BillingName        string    `json:"billing_name"`
	BillingAddress     string    `json:"billing_address"`
	BillingCity        string    `json:"billing_city"`
	BillingPostalCode  string    `json:"billing_postal_code"`
	BillingCountryID   uuid.UUID `json:"billing_country_id"`
	ShippingMethodID   uuid.UUID `json:"shipping_method_id"`
	PaymentMethodID    uuid.UUID `json:"payment_method_id"`
}

type OrderResponse struct {
	ID                 uuid.UUID                `json:"id"`
	UserID             *uuid.UUID               `json:"user_id"`
	Status             string                   `json:"status"`
	TotalPrice         float64                  `json:"total_price"`
	CreatedAt          *time.Time               `json:"created_at"`
	UpdatedAt          *time.Time               `json:"updated_at"`
	CustomerEmail      string                   `json:"customer_email"`
	ShippingName       string                   `json:"shipping_name"`
	ShippingAddress    string                   `json:"shipping_address"`
	ShippingCity       string                   `json:"shipping_city"`
	ShippingPostalCode string                   `json:"shipping_postal_code"`
	ShippingPhone      string                   `json:"shipping_phone"`
	BillingName        string                   `json:"billing_name"`
	BillingAddress     string                   `json:"billing_address"`
	BillingCity        string                   `json:"billing_city"`
	BillingPostalCode  string                   `json:"billing_postal_code"`
	ShippingMethodID   uuid.UUID                `json:"shipping_method_id"`
	ShippingPrice      float64                  `json:"shipping_price"`
	PaymentMethodID    uuid.UUID                `json:"payment_method_id"`
	ShippingCountryID  uuid.UUID                `json:"shipping_country_id"`
	BillingCountryID   uuid.UUID                `json:"billing_country_id"`
	CartItems          []database.OrdersVariant `json:"cart_items"`
}

func (cfg *apiConfig) handleApiCheckout(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)

	decoder.DisallowUnknownFields()

	params := CheckoutParams{}

	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}
	cartId, err := cfg.getOrCreateCartID(w, r)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to load cart")
		return
	}
	email := params.CustomerEmail
	userId := getUserIDFromContext(r.Context())

	if userId != uuid.Nil {
		user, err := cfg.db.GetUserById(r.Context(), userId)

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Unable to load user")
			return
		}

		email = user.Email
	} else {

		if !isValidEmail(email) {
			respondWithError(w, http.StatusBadRequest, "Invalid email")
			return
		}
	}

	if params.ShippingName == "" || params.ShippingAddress == "" ||
		params.ShippingCity == "" || params.ShippingPostalCode == "" ||
		params.ShippingCountryID == uuid.Nil || params.ShippingPhone == "" ||
		params.ShippingMethodID == uuid.Nil || params.PaymentMethodID == uuid.Nil ||
		params.BillingName == "" || params.BillingAddress == "" ||
		params.BillingCity == "" || params.BillingPostalCode == "" ||
		params.BillingCountryID == uuid.Nil {
		respondWithError(w, http.StatusBadRequest, "Missing required fields")
		return
	}

	shippingCountry, err := cfg.db.GetCountryById(r.Context(), params.ShippingCountryID)
	if err != nil || !shippingCountry.IsActive {
		respondWithError(w, http.StatusBadRequest, "Invalid shipping country")
		return
	}
	billingCountry, err := cfg.db.GetCountryById(r.Context(), params.BillingCountryID)
	if err != nil || !billingCountry.IsActive {
		respondWithError(w, http.StatusBadRequest, "Invalid billing country")
		return
	}

	items, err := cfg.db.GetCartDetailsWithSnapshotPrice(r.Context(), cartId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to load cart details")
		return
	}
	if len(items) == 0 {
		respondWithError(w, http.StatusBadRequest, "Cart is empty")
		return
	}

	subtotal := 0.0

	for _, item := range items {
		subtotal += float64(item.Quantity) * item.PricePerItem
	}

	shippingMethod, err := cfg.db.SelectShippingOptionById(r.Context(), params.ShippingMethodID)

	if err != nil || !shippingMethod.IsActive {
		respondWithError(w, http.StatusNotFound, "Shipping method not found")
		return
	}

	paymentMethod, err := cfg.db.GetPaymentOptionById(r.Context(), params.PaymentMethodID)
	if err != nil || !paymentMethod.IsActive {
		respondWithError(w, http.StatusNotFound, "Payment method not found")
		return
	}

	shippingPrice := shippingMethod.Price

	totalPrice := subtotal + shippingPrice
	tx, err := cfg.sqlDB.BeginTx(r.Context(), nil)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create order")
		return
	}

	defer func() {
		if tx != nil {
			_ = tx.Rollback()
		}
	}()

	qtx := cfg.db.WithTx(tx)

	for _, item := range items {
		_, err := qtx.DecreaseVariantStock(r.Context(), database.DecreaseVariantStockParams{
			Quantity:  item.Quantity,
			VariantID: item.ProductVariantID,
		})
		if err != nil {
			if err == sql.ErrNoRows {
				respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Insufficient stock for SKU %s", item.Sku))
				return
			}
			respondWithError(w, http.StatusInternalServerError, "Failed to reserve stock")
			return
		}
	}

	order, err := qtx.CreateOrder(r.Context(), database.CreateOrderParams{
		UserID:             uuid.NullUUID{UUID: userId, Valid: userId != uuid.Nil},
		TotalPrice:         totalPrice,
		CustomerEmail:      email,
		ShippingName:       params.ShippingName,
		ShippingAddress:    params.ShippingAddress,
		ShippingCity:       params.ShippingCity,
		ShippingPostalCode: params.ShippingPostalCode,
		ShippingCountryID:  params.ShippingCountryID,
		ShippingPhone:      params.ShippingPhone,
		BillingName:        params.BillingName,
		BillingAddress:     params.BillingAddress,
		BillingCity:        params.BillingCity,
		BillingPostalCode:  params.BillingPostalCode,
		BillingCountryID:   params.BillingCountryID,
		ShippingOptionID:   params.ShippingMethodID,
		PaymentOptionID:    params.PaymentMethodID,
		ShippingPrice:      shippingPrice,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Cannot create order")
		return
	}

	cartItems, err := qtx.CopyCartDataIntoOrder(r.Context(), database.CopyCartDataIntoOrderParams{OrderID: order.ID, CartID: cartId})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Cannot copy cart items into order")
		return
	}

	if _, err := qtx.UpdateCartStatus(r.Context(), database.UpdateCartStatusParams{
		Status: "completed",
		CartID: cartId,
	}); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update cart")
		return
	}

	if err := tx.Commit(); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Cannot finalize order")
		return
	}

	tx = nil

	var userIDPtr *uuid.UUID
	if order.UserID.Valid {
		userIDPtr = &order.UserID.UUID
	} else {
		userIDPtr = nil
	}

	var createdAtPtr, updatedAtPtr *time.Time
	if order.CreatedAt.Valid {
		createdAtPtr = &order.CreatedAt.Time
	} else {
		createdAtPtr = nil
	}
	if order.UpdatedAt.Valid {
		updatedAtPtr = &order.UpdatedAt.Time
	} else {
		updatedAtPtr = nil
	}

	resp := OrderResponse{
		ID:                 order.ID,
		UserID:             userIDPtr,
		Status:             string(order.Status),
		TotalPrice:         order.TotalPrice,
		CreatedAt:          createdAtPtr,
		UpdatedAt:          updatedAtPtr,
		CustomerEmail:      order.CustomerEmail,
		ShippingName:       order.ShippingName,
		ShippingAddress:    order.ShippingAddress,
		ShippingCity:       order.ShippingCity,
		ShippingPostalCode: order.ShippingPostalCode,
		ShippingPhone:      order.ShippingPhone,
		BillingName:        order.BillingName,
		BillingAddress:     order.BillingAddress,
		BillingCity:        order.BillingCity,
		BillingPostalCode:  order.BillingPostalCode,
		ShippingMethodID:   order.ShippingOptionID,
		ShippingPrice:      order.ShippingPrice,
		PaymentMethodID:    order.PaymentOptionID,
		ShippingCountryID:  order.ShippingCountryID,
		BillingCountryID:   order.BillingCountryID,
		CartItems:          cartItems,
	}

	respondWithJSON(w, http.StatusOK, resp)
}
