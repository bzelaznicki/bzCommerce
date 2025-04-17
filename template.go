package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/bzelaznicki/bzCommerce/internal/database"
	"github.com/google/uuid"
)

type BasePageData struct {
	StoreName  string
	IsLoggedIn bool
	IsAdmin    bool
	Categories []CategoryGroup
	Data       any
}

func (cfg *apiConfig) NewPageData(ctx context.Context, data any) (BasePageData, error) {
	categories, err := cfg.db.GetCategories(ctx)
	if err != nil {
		return BasePageData{}, fmt.Errorf("error getting base page data: %v", err)
	}

	var isLoggedIn bool
	var isAdmin bool

	userIDStr, ok := ctx.Value(userIDContextKey).(string)
	if ok {
		isLoggedIn = true

		userID, err := uuid.Parse(userIDStr)
		if err == nil {
			user, err := cfg.db.GetUserById(ctx, userID)
			if err == nil {
				isAdmin = user.IsAdmin
			}
		}
	}

	cats := groupCategories(categories)

	return BasePageData{
		StoreName:  cfg.storeName,
		Categories: cats,
		IsLoggedIn: isLoggedIn,
		IsAdmin:    isAdmin,
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

type CategoryGroup struct {
	Parent   database.Category
	Children []database.Category
}

func groupCategories(cats []database.Category) []CategoryGroup {
	parentMap := make(map[uuid.UUID][]database.Category)
	var topLevel []database.Category

	for _, cat := range cats {
		if cat.ParentID.Valid {
			parentMap[cat.ParentID.UUID] = append(parentMap[cat.ParentID.UUID], cat)
		} else {
			topLevel = append(topLevel, cat)
		}
	}

	var groups []CategoryGroup
	for _, parent := range topLevel {
		children := parentMap[parent.ID]
		groups = append(groups, CategoryGroup{Parent: parent, Children: children})
	}

	return groups
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

func (cfg *apiConfig) RenderError(w http.ResponseWriter, r *http.Request, code int, message string) {
	w.WriteHeader(code)

	data := struct {
		Code    int
		Message string
	}{
		Code:    code,
		Message: message,
	}

	cfg.Render(w, r, "templates/pages/error.html", data)
}
