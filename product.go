package main

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID
	Name        string
	Slug        string
	ImagePath   string
	Description string
	CategoryID  uuid.UUID
	Variants    []Variant
}

type Variant struct {
	ID            uuid.UUID
	Name          string
	ProductID     uuid.UUID
	Price         float64
	StockQuantity int32
	ImageUrl      string
	VariantName   string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
