package main

import (
	"log"
	"os"

	"github.com/bzelaznicki/bzCommerce/internal/database"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	db database.Client
}

func main() {
	godotenv.Load(".env")

	pathToDB := os.Getenv("DB_PATH")
	if pathToDB == "" {
		log.Fatal("DB_URL must be set")
	}

	db, err := database.NewClient(pathToDB)

	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET cannot be empty")
	}
}
