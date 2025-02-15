package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/bzelaznicki/bzCommerce/internal/database"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	db           *database.Queries
	jwtSecret    string
	filepathRoot string
	platform     string
}

func main() {
	godotenv.Load(".env")

	pathToDB := os.Getenv("DB_URL")
	if pathToDB == "" {
		log.Fatal("DB_URL must be set")
	}

	db, err := sql.Open("postgres", pathToDB)

	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	dbQueries := database.New(db)

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET cannot be empty")
	}

	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM cannot be empty")
	}

	filepathRoot := os.Getenv("FILEPATH_ROOT")

	if filepathRoot == "" {
		log.Fatal("FILEPATH_ROOT cannot be empty")
	}

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("PORT cannot be empty")
	}

	cfg := apiConfig{
		db:        dbQueries,
		jwtSecret: jwtSecret,
		platform:  platform,
	}

	mux := http.NewServeMux()
	srv := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	appHandler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))

	mux.Handle("/app/", appHandler)
	mux.HandleFunc("POST /api/users", cfg.handleUserCreate)

	fmt.Printf("serving files from %s on port %s\n", filepathRoot, port)

	log.Fatal(srv.ListenAndServe())
}
