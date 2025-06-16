package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("Failed to set Goose dialect: %v", err)
	}

	migrationsDir := "backend/sql/schema"

	if err := goose.Up(db, migrationsDir); err != nil {
		log.Fatalf("Goose migration failed: %v", err)
	}

	log.Println("Database migrated successfully.")
}
