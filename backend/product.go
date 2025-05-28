package main

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	ImagePath   string    `json:"imagePath"`
	Description string    `json:"description"`
	CategoryID  uuid.UUID `json:"categoryId"`
	Variants    []Variant `json:"variants"`
}

type Variant struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	ProductID     uuid.UUID `json:"productId"`
	Price         float64   `json:"price"`
	StockQuantity int32     `json:"stockQuantity"`
	ImageUrl      string    `json:"imageUrl"`
	VariantName   string    `json:"variantName"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
