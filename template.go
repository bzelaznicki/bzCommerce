package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/bzelaznicki/bzCommerce/internal/database"
)

type BasePageData struct {
	StoreName  string
	Categories []database.Category
	IsLoggedIn bool
	Data       any
}

func (cfg *apiConfig) NewPageData(ctx context.Context, data any) (BasePageData, error) {
	categories, err := cfg.db.GetCategories(ctx)
	if err != nil {
		return BasePageData{}, fmt.Errorf("error getting base page data: %v", err)
	}
	_, isLoggedIn := ctx.Value(userIDContextKey).(string)
	return BasePageData{
		StoreName:  cfg.storeName,
		Categories: categories,
		IsLoggedIn: isLoggedIn,
		Data:       data,
	}, nil
}

func parsePageTemplate(filename string) *template.Template {
	tmpl := template.Must(template.ParseFiles(
		"templates/base.html",
		filename,
	))
	return tmpl
}

func (cfg *apiConfig) Render(w http.ResponseWriter, r *http.Request, page string, data any) {
	pageData, err := cfg.NewPageData(r.Context(), data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	tmpl := parsePageTemplate(page)
	err = tmpl.ExecuteTemplate(w, "base.html", pageData)
	if err != nil {
		log.Printf("template error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
