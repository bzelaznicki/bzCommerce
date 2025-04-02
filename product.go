package main

import "github.com/google/uuid"

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
	ID   uuid.UUID
	Name string
}
