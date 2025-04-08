package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/bzelaznicki/bzCommerce/internal/database"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	db        *database.Queries
	jwtSecret string
	platform  string
	templates *template.Template
	storeName string
}

func main() {
	logger()
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

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("PORT cannot be empty")
	}
	storeName := os.Getenv("STORE_NAME")
	if storeName == "" {
		storeName = "bzCommerce" // fallback default
	}
	templates := template.Must(template.ParseFiles(
		"templates/base.html",
	))

	cfg := apiConfig{
		db:        dbQueries,
		jwtSecret: jwtSecret,
		platform:  platform,
		templates: templates,
		storeName: storeName,
	}

	mux := http.NewServeMux()

	srv := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	cfg.registerRoutes(mux)

	fmt.Printf("serving on port %s\n", port)

	log.Fatal(srv.ListenAndServe())
}
