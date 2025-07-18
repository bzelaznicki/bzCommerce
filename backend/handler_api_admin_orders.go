package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/bzelaznicki/bzCommerce/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handleApiAdminListOrders(w http.ResponseWriter, r *http.Request) {
	page, limit := getPaginationParams(r)
	offset := (page - 1) * limit

	q := r.URL.Query()

	var statusParam string
	if status := q.Get("status"); status != "" {
		statusParam = status
	}

	var paymentStatusParam string
	if paymentStatus := q.Get("payment_status"); paymentStatus != "" {
		paymentStatusParam = paymentStatus
	}

	var searchParam string
	if search := q.Get("search"); search != "" {
		searchParam = search
	}

	var dateFromParam time.Time
	if dateFromStr := q.Get("date_from"); dateFromStr != "" {
		if t, err := time.Parse("2006-01-02", dateFromStr); err == nil {
			dateFromParam = t
		}
	}

	var dateToParam time.Time
	if dateToStr := q.Get("date_to"); dateToStr != "" {
		if t, err := time.Parse("2006-01-02", dateToStr); err == nil {
			dateToParam = t
		}
	}

	ctx := r.Context()

	count, err := cfg.db.CountOrders(ctx, database.CountOrdersParams{
		Search:        sql.NullString{String: searchParam, Valid: searchParam != ""},
		Status:        database.NullOrderStatus{OrderStatus: database.OrderStatus(statusParam), Valid: statusParam != ""},
		PaymentStatus: database.NullPaymentStatus{PaymentStatus: database.PaymentStatus(paymentStatusParam), Valid: paymentStatusParam != ""},
		DateFrom:      sql.NullTime{Time: dateFromParam, Valid: !dateFromParam.IsZero()},
		DateTo:        sql.NullTime{Time: dateToParam, Valid: !dateToParam.IsZero()},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to count orders")
		log.Printf("Count orders error: %v", err)
		return
	}

	orders, err := cfg.db.ListOrders(ctx,
		database.ListOrdersParams{
			Limit:         limit,
			Offset:        offset,
			Search:        sql.NullString{String: searchParam, Valid: searchParam != ""},
			Status:        database.NullOrderStatus{OrderStatus: database.OrderStatus(statusParam), Valid: statusParam != ""},
			PaymentStatus: database.NullPaymentStatus{PaymentStatus: database.PaymentStatus(paymentStatusParam), Valid: paymentStatusParam != ""},
			DateFrom:      sql.NullTime{Time: dateFromParam, Valid: !dateFromParam.IsZero()},
			DateTo:        sql.NullTime{Time: dateToParam, Valid: !dateToParam.IsZero()},
		},
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to list orders")
		log.Printf("List orders error: %v", err)
		return
	}

	if orders == nil {
		orders = []database.ListOrdersRow{}
	}

	response := NewPaginatedResponse(orders, page, limit, count)
	respondWithJSON(w, http.StatusOK, response)
}

func (cfg *apiConfig) handleApiAdminGetSingleOrder(w http.ResponseWriter, r *http.Request) {
	orderId, err := uuid.Parse(r.PathValue("orderId"))

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid order ID")
		return
	}

	order, err := cfg.db.GetOrderWithUserById(r.Context(), orderId)

	if err != nil {
		respondWithError(w, http.StatusNotFound, "Order not found")
		return
	}

	orderItems, err := cfg.db.GetOrderItemsByOrderIdWithVariants(r.Context(), orderId)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to get order items")
		return
	}

	if orderItems == nil {
		orderItems = []database.GetOrderItemsByOrderIdWithVariantsRow{}
	}

	resp := struct {
		OrderID            uuid.UUID                                        `json:"order_id"`
		UserID             uuid.NullUUID                                    `json:"user_id"`
		Status             string                                           `json:"status"`
		PaymentStatus      string                                           `json:"payment_status"`
		TotalPrice         float64                                          `json:"total_price"`
		CreatedAt          time.Time                                        `json:"created_at"`
		UpdatedAt          time.Time                                        `json:"updated_at"`
		CustomerEmail      string                                           `json:"customer_email"`
		ShippingName       string                                           `json:"shipping_name"`
		ShippingAddress    string                                           `json:"shipping_address"`
		ShippingCity       string                                           `json:"shipping_city"`
		ShippingPostalCode string                                           `json:"shipping_postal_code"`
		ShippingPhone      string                                           `json:"shipping_phone"`
		BillingName        string                                           `json:"billing_name"`
		BillingAddress     string                                           `json:"billing_address"`
		BillingCity        string                                           `json:"billing_city"`
		BillingPostalCode  string                                           `json:"billing_postal_code"`
		ShippingOptionID   uuid.UUID                                        `json:"shipping_option_id"`
		ShippingPrice      float64                                          `json:"shipping_price"`
		PaymentOptionID    uuid.UUID                                        `json:"payment_option_id"`
		ShippingCountryID  uuid.UUID                                        `json:"shipping_country_id"`
		BillingCountryID   uuid.UUID                                        `json:"billing_country_id"`
		ShippingMethodName sql.NullString                                   `json:"shipping_method_name"`
		PaymentMethodName  sql.NullString                                   `json:"payment_method_name"`
		UserEmail          sql.NullString                                   `json:"user_email"`
		UserCreatedAt      sql.NullTime                                     `json:"user_created_at"`
		OrderItems         []database.GetOrderItemsByOrderIdWithVariantsRow `json:"order_items"`
	}{
		OrderID:            order.ID,
		UserID:             order.UserID,
		Status:             string(order.Status),
		PaymentStatus:      string(order.PaymentStatus),
		TotalPrice:         order.TotalPrice,
		CreatedAt:          order.CreatedAt,
		UpdatedAt:          order.UpdatedAt,
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
		ShippingOptionID:   order.ShippingOptionID,
		ShippingPrice:      order.ShippingPrice,
		PaymentOptionID:    order.PaymentOptionID,
		ShippingCountryID:  order.ShippingCountryID,
		BillingCountryID:   order.BillingCountryID,
		ShippingMethodName: order.ShippingMethodName,
		PaymentMethodName:  order.PaymentMethodName,
		UserEmail:          order.UserEmail,
		UserCreatedAt:      order.UserCreatedAt,
		OrderItems:         orderItems,
	}
	respondWithJSON(w, http.StatusOK, resp)
}
