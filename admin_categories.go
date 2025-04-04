package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bzelaznicki/bzCommerce/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handleAdminCategoryList(w http.ResponseWriter, r *http.Request) {
	rows, err := cfg.db.ListCategoriesWithParent(r.Context())

	if err != nil {
		http.Error(w, "Could not load categories", http.StatusInternalServerError)
		log.Printf("could not load categories: %v", err)
		return
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

	cfg.Render(w, r, "templates/pages/admin_categories.html", categories)
}

func (cfg *apiConfig) handleAdminCategoryNewForm(w http.ResponseWriter, r *http.Request) {
	categories, err := cfg.db.GetCategories(r.Context())
	if err != nil {
		http.Error(w, "Failed to load categories", http.StatusInternalServerError)
		log.Printf("failed to load categories: %v", err)
		return
	}

	data := struct {
		IsEdit     bool
		FormAction string
		Category   struct {
			ID          uuid.UUID
			Name        string
			Slug        string
			Description string
			ParentID    uuid.NullUUID
		}
		Categories []database.Category
	}{
		IsEdit:     false,
		FormAction: "/admin/categories",
		Category: struct {
			ID          uuid.UUID
			Name        string
			Slug        string
			Description string
			ParentID    uuid.NullUUID
		}{},
		Categories: categories,
	}

	cfg.Render(w, r, "templates/pages/admin_category_form.html", data)
}

func (cfg *apiConfig) handleAdminCategoryCreate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
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
		http.Error(w, "Failed to create category", http.StatusInternalServerError)
		log.Printf("failed to create category: %v", err)
		return
	}

	http.Redirect(w, r, "/admin/categories", http.StatusSeeOther)
}
func (cfg *apiConfig) handleAdminCategoryEditForm(w http.ResponseWriter, r *http.Request) {

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.NotFound(w, r)
		log.Printf("failed to parse UUID: %v", err)
		return
	}

	cat, err := cfg.db.GetCategoryById(r.Context(), id)

	if err != nil {
		http.Error(w, "Category not found", http.StatusNotFound)
		log.Printf("failed to find category: %v", err)
		return
	}

	categories, err := cfg.db.GetCategories(r.Context())
	if err != nil {
		http.Error(w, "Failed to load categories", http.StatusInternalServerError)
		log.Printf("failed to load categories: %v", err)
		return
	}

	data := struct {
		IsEdit     bool
		FormAction string
		Category   struct {
			ID          uuid.UUID
			Name        string
			Slug        string
			Description string
			ParentID    uuid.NullUUID
		}
		Categories []database.Category
	}{
		IsEdit:     true,
		FormAction: fmt.Sprintf("/admin/categories/%s", id),
		Category: struct {
			ID          uuid.UUID
			Name        string
			Slug        string
			Description string
			ParentID    uuid.NullUUID
		}{
			ID:          cat.ID,
			Name:        cat.Name,
			Slug:        cat.Slug,
			Description: cat.Description.String,
			ParentID:    cat.ParentID,
		},
		Categories: categories,
	}

	cfg.Render(w, r, "templates/pages/admin_category_form.html", data)
}

func (cfg *apiConfig) handleAdminCategoryUpdate(w http.ResponseWriter, r *http.Request) {

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.NotFound(w, r)
		log.Printf("invalid product id: %v", err)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
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
		http.Error(w, "failed to update category", http.StatusInternalServerError)
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
		http.Error(w, "Category not found", http.StatusNotFound)
		log.Printf("category not found (id %s): %v", id, err)
		return
	}

	err = cfg.db.DeleteCategoryById(r.Context(), id)

	if err != nil {
		http.Error(w, "Failed to delete category", http.StatusInternalServerError)
		log.Printf("failed to delete category (id %s): %v", id, err)
		return
	}

	http.Redirect(w, r, "/admin/products", http.StatusSeeOther)
}
