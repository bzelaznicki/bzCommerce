package main

/*
func (cfg *apiConfig) handleApiAdminListOrders(w http.ResponseWriter, r *http.Request) {
	page, limit := getPaginationParams(r)
	offset := (page - 1) * limit

	q := r.URL.Query()

	var statusParam *string
	if status := q.Get("status"); status != "" {
		statusParam = &status
	}

	var paymentStatusParam *string
	if paymentStatus := q.Get("payment_status"); paymentStatus != "" {
		paymentStatusParam = &paymentStatus
	}

	var searchParam *string
	if search := q.Get("search"); search != "" {
		searchParam = &search
	}

	var dateFromParam *time.Time
	if dateFromStr := q.Get("date_from"); dateFromStr != "" {
		if t, err := time.Parse("2006-01-02", dateFromStr); err == nil {
			dateFromParam = &t
		}
	}

	var dateToParam *time.Time
	if dateToStr := q.Get("date_to"); dateToStr != "" {
		if t, err := time.Parse("2006-01-02", dateToStr); err == nil {
			dateToParam = &t
		}
	}

	ctx := r.Context()

	count, err := cfg.db.CountOrders(ctx, database.CountOrdersParams{
		Search:        searchParam,
		Status:        statusParam,
		PaymentStatus: paymentStatusParam,
		DateFrom:      dateFromParam,
		DateTo:        dateToParam,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to count orders")
		return
	}

	orders, err := cfg.db.ListOrders(ctx,
		database.ListOrdersParams{
			Limit:         limit,
			Offset:        offset,
			Search:        searchParam,
			Status:        statusParam,
			PaymentStatus: paymentStatusParam,
			DateFrom:      dateFromParam,
			DateTo:        dateToParam,
		},
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to list orders")
		return
	}

	// Return paginated response
	response := NewPaginatedResponse(orders, page, limit, count)
	respondWithJSON(w, http.StatusOK, response)
}
*/
