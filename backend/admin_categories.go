package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/bzelaznicki/bzCommerce/internal/database"
	"github.com/google/uuid"
)

type AdminCategoriesListPageData struct {
	Categories  []AdminCategoryRow
	Breadcrumbs []Breadcrumb
}

type AdminCategoryFormPageData struct {
	IsEdit          bool
	FormAction      string
	Category        AdminCategoryForm
	CategoryOptions []CategoryOption
	Breadcrumbs     []Breadcrumb
}

type AdminCategoryRow struct {
	ID          uuid.UUID
	Name        string
	Slug        string
	Description string
	ParentName  sql.NullString
	ParentID    uuid.NullUUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type AdminCategoryForm struct {
	ID          uuid.UUID
	Name        string
	Slug        string
	Description string
	ParentID    uuid.NullUUID
}

func (cfg *apiConfig) handleAdminCategoryList(w http.ResponseWriter, r *http.Request) {
	rows, err := cfg.db.ListCategoriesWithParent(r.Context())

	if err != nil {
		cfg.RenderError(w, r, http.StatusInternalServerError, "Could not load categories")
		log.Printf("could not load categories: %v", err)
		return
	}

	categories := make([]AdminCategoryRow, 0, len(rows))
	for _, row := range rows {
		categories = append(categories, AdminCategoryRow{
			ID:          row.ID,
			Name:        row.Name,
			Slug:        row.Slug,
			Description: row.Description.String,
			ParentName:  row.ParentName,
			ParentID:    row.ParentID,
			CreatedAt:   row.CreatedAt,
			UpdatedAt:   row.UpdatedAt,
		})
	}

	data := AdminCategoriesListPageData{
		Categories: categories,
		Breadcrumbs: NewBreadcrumbTrail(
			Breadcrumb{Label: "Categories"},
		),
	}

	cfg.Render(w, r, "templates/pages/admin_categories.html", data)

}

func (cfg *apiConfig) handleAdminCategoryNewForm(w http.ResponseWriter, r *http.Request) {
	categories, err := cfg.db.GetCategories(r.Context())
	if err != nil {
		cfg.RenderError(w, r, http.StatusInternalServerError, "Failed to load categories")
		log.Printf("failed to load categories: %v", err)
		return
	}

	data := AdminCategoryFormPageData{
		IsEdit:          false,
		FormAction:      "/admin/categories",
		Category:        AdminCategoryForm{},
		CategoryOptions: flattenCategoryOptions(groupCategories(categories)),
		Breadcrumbs: NewBreadcrumbTrail(
			Breadcrumb{Label: "Categories", URL: "/admin/categories"},
			Breadcrumb{Label: "New"},
		),
	}

	cfg.Render(w, r, "templates/pages/admin_category_form.html", data)

}

func (cfg *apiConfig) handleAdminCategoryCreate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		cfg.RenderError(w, r, http.StatusBadRequest, "Invalid form")
		log.Printf("invalid form: %v", err)
		return
	}

	parentID := r.FormValue("parent_id")

	var parsedParent uuid.NullUUID
	if parentID != "" {
		if id, err := uuid.Parse(parentID); err == nil {
			parsedParent = uuid.NullUUID{UUID: id, Valid: true}
		}
	}

	_, err := cfg.db.CreateCategory(r.Context(), database.CreateCategoryParams{
		Name:        r.FormValue("name"),
		Slug:        r.FormValue("slug"),
		Description: sql.NullString{String: r.FormValue("description"), Valid: r.FormValue("description") != ""},
		ParentID:    parsedParent,
	})
	if err != nil {
		cfg.RenderError(w, r, http.StatusInternalServerError, "Failed to create category")
		log.Printf("failed to create category: %v", err)
		return
	}

	http.Redirect(w, r, "/admin/categories", http.StatusSeeOther)
}
func (cfg *apiConfig) handleAdminCategoryEditForm(w http.ResponseWriter, r *http.Request) {

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		cfg.RenderError(w, r, http.StatusNotFound, "Category not found")
		log.Printf("failed to parse category UUID: %v", err)
		return
	}

	cat, err := cfg.db.GetCategoryById(r.Context(), id)

	if err != nil {
		cfg.RenderError(w, r, http.StatusNotFound, fmt.Sprintf("Category %s not found", id))
		log.Printf("failed to find category %s: %v", id, err)
		return
	}

	categories, err := cfg.db.GetCategories(r.Context())
	if err != nil {
		cfg.RenderError(w, r, http.StatusInternalServerError, "Failed to load categories")
		log.Printf("failed to load categories: %v", err)
		return
	}

	data := AdminCategoryFormPageData{
		IsEdit:     true,
		FormAction: fmt.Sprintf("/admin/categories/%s", id),
		Category: AdminCategoryForm{
			ID:          cat.ID,
			Name:        cat.Name,
			Slug:        cat.Slug,
			Description: cat.Description.String,
			ParentID:    cat.ParentID,
		},
		CategoryOptions: flattenCategoryOptions(groupCategories(categories)),
		Breadcrumbs: NewBreadcrumbTrail(
			Breadcrumb{Label: "Categories", URL: "/admin/categories"},
			Breadcrumb{Label: "Edit"},
		),
	}

	cfg.Render(w, r, "templates/pages/admin_category_form.html", data)

}

func (cfg *apiConfig) handleAdminCategoryUpdate(w http.ResponseWriter, r *http.Request) {

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		cfg.RenderError(w, r, http.StatusNotFound, "Category not found")
		log.Printf("invalid product id: %v", err)
		return
	}

	if err := r.ParseForm(); err != nil {
		cfg.RenderError(w, r, http.StatusBadRequest, "Invalid form")
		log.Printf("error parsing form: %v", err)
		return
	}

	parentID := r.FormValue("parent_id")
	var parsedParent uuid.NullUUID
	if parentID != "" {
		if id, err := uuid.Parse(parentID); err == nil {
			parsedParent = uuid.NullUUID{UUID: id, Valid: true}
		}
	}

	err = cfg.db.UpdateCategoryById(r.Context(), database.UpdateCategoryByIdParams{
		ID:          id,
		Name:        r.FormValue("name"),
		Slug:        r.FormValue("slug"),
		Description: sql.NullString{String: r.FormValue("description"), Valid: r.FormValue("description") != ""},
		ParentID:    parsedParent,
	})

	if err != nil {
		cfg.RenderError(w, r, http.StatusInternalServerError, fmt.Sprintf("Failed to update category %s", id))
		log.Printf("failed to update category: %v", err)
		return
	}

	http.Redirect(w, r, "/admin/categories", http.StatusSeeOther)

}
func (cfg *apiConfig) handleAdminCategoryDelete(w http.ResponseWriter, r *http.Request) {

	id, err := uuid.Parse(r.PathValue("id"))

	if err != nil {
		http.NotFound(w, r)
		log.Printf("invalid category id: %v", err)
		return
	}

	_, err = cfg.db.GetCategoryById(r.Context(), id)

	if err != nil {
		cfg.RenderError(w, r, http.StatusNotFound, fmt.Sprintf("Category %s not found", id))
		log.Printf("category not found (id %s): %v", id, err)
		return
	}

	err = cfg.db.DeleteCategoryById(r.Context(), id)

	if err != nil {
		cfg.RenderError(w, r, http.StatusInternalServerError, fmt.Sprintf("Failed to delete category %s", id))
		log.Printf("failed to delete category (id %s): %v", id, err)
		return
	}

	http.Redirect(w, r, "/admin/products", http.StatusSeeOther)
}

type CategoryOption struct {
	ID    uuid.UUID
	Name  string
	Depth int
}

func flattenCategoryOptions(groups []CategoryGroup) []CategoryOption {
	var options []CategoryOption

	walk := func(cat database.Category, depth int) {
		prefix := strings.Repeat("--", depth)
		options = append(options, CategoryOption{
			ID:    cat.ID,
			Name:  prefix + " " + cat.Name,
			Depth: depth,
		})
	}

	for _, group := range groups {
		walk(group.Parent, 0)
		for _, child := range group.Children {
			walk(child, 1)
		}
	}

	return options
}
