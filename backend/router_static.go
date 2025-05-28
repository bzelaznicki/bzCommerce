package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) registerStaticRoutes(mux *http.ServeMux) {
	log.Printf("Registering static routes...")
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	log.Printf("Static routes registered")
}
