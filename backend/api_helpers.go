package main

import (
	"net/http"
	"strconv"
)

type PaginatedResponse[T any] struct {
	Data       []T   `json:"data"`
	Page       int64 `json:"page"`
	Limit      int64 `json:"limit"`
	TotalCount int64 `json:"total_count"`
	TotalPages int64 `json:"total_pages"`
}

func NewPaginatedResponse[T any](items []T, page, limit, totalCount int64) PaginatedResponse[T] {
	totalPages := (totalCount + limit - 1) / limit
	return PaginatedResponse[T]{
		Data:       items,
		Page:       page,
		Limit:      limit,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}
}

func getPaginationParams(r *http.Request) (int64, int64) {
	query := r.URL.Query()

	page := 1
	limit := 20

	if p := query.Get("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if l := query.Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	return int64(page), int64(limit)
}
