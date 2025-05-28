package main

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID
	FullName  string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
	IsAdmin   bool
}
