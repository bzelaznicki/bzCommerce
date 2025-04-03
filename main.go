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
	db           *database.Queries
	jwtSecret    string
	filepathRoot string
	platform     string
	templates    *template.Template
	storeName    string
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
	storeName := os.Getenv("STORE_NAME")
	if storeName == "" {
		storeName = "bzCommerce" // fallback default
	}
	templates := template.Must(template.ParseFiles(
		"templates/base.html",
		"templates/pages/home.html",
		"templates/pages/product.html",
		"templates/pages/category.html",
	))

	for _, tmpl := range templates.Templates() {
		log.Println("Loaded template:", tmpl.Name())
	}
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
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	mux.Handle("GET /", cfg.maybeWithAuth(http.HandlerFunc(cfg.handleHomePage)))
	mux.Handle("GET /product/{slug}", cfg.maybeWithAuth(http.HandlerFunc(cfg.handleProductPage)))
	mux.Handle("GET /category/{slug}", cfg.maybeWithAuth(http.HandlerFunc(cfg.handleCategoryPage)))

	mux.HandleFunc("POST /api/users", cfg.handleUserCreate)
	mux.Handle("GET /login", cfg.redirectIfAuthenticated(http.HandlerFunc(cfg.handleLoginGet)))
	mux.Handle("GET /register", cfg.redirectIfAuthenticated(http.HandlerFunc(cfg.handleRegisterGet)))
	mux.Handle("POST /login", cfg.redirectIfAuthenticated(http.HandlerFunc(cfg.handleLoginPost)))
	mux.Handle("POST /register", cfg.redirectIfAuthenticated(http.HandlerFunc(cfg.handleRegisterPost)))

	mux.Handle("GET /account", cfg.withAuth(http.HandlerFunc(cfg.handleAccountPage)))

	mux.HandleFunc("GET /logout", cfg.handleLogout)

	fmt.Printf("serving files from %s on port %s\n", filepathRoot, port)

	log.Fatal(srv.ListenAndServe())
}
