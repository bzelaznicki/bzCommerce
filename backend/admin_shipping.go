package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/bzelaznicki/bzCommerce/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handleAdminShippingOptionsList(w http.ResponseWriter, r *http.Request) {
	rows, err := cfg.db.GetShippingOptions(r.Context())

	if err != nil {
		cfg.RenderError(w, r, http.StatusInternalServerError, "Could not load shipping options")
		log.Printf("could not load shipping options: %v", err)
		return
	}

	type AdminShippingListPageData struct {
		ShippingOptions []database.ShippingOption
		Breadcrumbs     []Breadcrumb
	}

	data := AdminShippingListPageData{
		ShippingOptions: rows,
		Breadcrumbs: NewBreadcrumbTrail(
			Breadcrumb{Label: "Shipping Options"},
		),
	}
	cfg.Render(w, r, "templates/pages/admin_shipping_list.html", data)
}

type AdminShippingFormPageData struct {
	ShippingOption database.ShippingOption
	IsEdit         bool
	Breadcrumbs    []Breadcrumb
}

func (cfg *apiConfig) handleAdminShippingOptionsNew(w http.ResponseWriter, r *http.Request) {
	data := AdminShippingFormPageData{
		ShippingOption: database.ShippingOption{},
		IsEdit:         false,
		Breadcrumbs: NewBreadcrumbTrail(
			Breadcrumb{Label: "Shipping Options", URL: "/admin/shipping"},
			Breadcrumb{Label: "New"},
		),
	}

	cfg.Render(w, r, "templates/pages/admin_shipping_form.html", data)
}

func (cfg *apiConfig) handleAdminShippingOptionCreate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		cfg.RenderError(w, r, http.StatusBadRequest, "Invalid form")
		log.Printf("error parsing form: %v", err)
		return
	}

	name := r.FormValue("name")
	description := r.FormValue("description")
	price := r.FormValue("price")
	estimatedDays := r.FormValue("estimated_days")
	sortOrder := r.FormValue("sort_order")
	isActive := r.FormValue("is_active") != ""

	sortOrderInt64, err := strconv.ParseInt(sortOrder, 10, 32)
	if err != nil {
		cfg.RenderError(w, r, http.StatusBadRequest, "Invalid sort order")
		log.Printf("invalid sort order: %v", err)
		return
	}

	sortOrderInt32 := int32(sortOrderInt64)
	_, err = cfg.db.CreateShippingOption(r.Context(), database.CreateShippingOptionParams{
		Name: name,
		Description: sql.NullString{
			String: description,
			Valid:  true,
		},
		Price:         price,
		EstimatedDays: estimatedDays,
		SortOrder:     sortOrderInt32,
		IsActive:      isActive,
	})

	if err != nil {
		cfg.RenderError(w, r, http.StatusInternalServerError, "Failed to create shipping option")
		log.Printf("failed to create shipping option: %v", err)
		return
	}
	http.Redirect(w, r, "/admin/shipping", http.StatusSeeOther)
}

func (cfg *apiConfig) handleAdminShippingOptionsEdit(w http.ResponseWriter, r *http.Request) {

	idStr := r.PathValue("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		cfg.RenderError(w, r, http.StatusBadRequest, "Invalid shipping option ID")
		log.Printf("invalid shipping option ID: %v", err)
		return
	}
	shippingOption, err := cfg.db.SelectShippingOptionById(r.Context(), id)
	if err != nil {
		cfg.RenderError(w, r, http.StatusInternalServerError, "Failed to load shipping option")
		log.Printf("failed to load shipping option: %v", err)
		return
	}
	data := AdminShippingFormPageData{
		ShippingOption: shippingOption,
		IsEdit:         true,
		Breadcrumbs: NewBreadcrumbTrail(
			Breadcrumb{Label: "Shipping Options", URL: "/admin/shipping"},
			Breadcrumb{Label: "New"},
		),
	}

	cfg.Render(w, r, "templates/pages/admin_shipping_form.html", data)
}

func (cfg *apiConfig) handleAdminShippingOptionUpdate(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		cfg.RenderError(w, r, http.StatusBadRequest, "Invalid shipping option ID")
		log.Printf("invalid shipping option ID: %v", err)
		return
	}

	if err := r.ParseForm(); err != nil {
		cfg.RenderError(w, r, http.StatusBadRequest, "Invalid form")
		log.Printf("error parsing form: %v", err)
		return
	}

	name := r.FormValue("name")
	description := r.FormValue("description")
	price := r.FormValue("price")
	estimatedDays := r.FormValue("estimated_days")
	sortOrder := r.FormValue("sort_order")
	isActive := r.PostFormValue("is_active") != ""

	sortOrderInt64, err := strconv.ParseInt(sortOrder, 10, 32)
	if err != nil {
		cfg.RenderError(w, r, http.StatusBadRequest, "Invalid sort order")
		log.Printf("invalid sort order: %v", err)
		return
	}

	sortOrderInt32 := int32(sortOrderInt64)
	err = cfg.db.UpdateShippingOption(r.Context(), database.UpdateShippingOptionParams{
		ID:   id,
		Name: name,
		Description: sql.NullString{
			String: description,
			Valid:  true,
		},
		Price:         price,
		EstimatedDays: estimatedDays,
		IsActive:      isActive,
		SortOrder:     sortOrderInt32,
	},
	)

	if err != nil {
		cfg.RenderError(w, r, http.StatusInternalServerError, "Failed to update shipping option")
		log.Printf("failed to update shipping option: %v", err)
		return
	}
	http.Redirect(w, r, "/admin/shipping", http.StatusSeeOther)
}

func (cfg *apiConfig) handleAdminShippingOptionDelete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		cfg.RenderError(w, r, http.StatusBadRequest, "Invalid shipping option ID")
		log.Printf("invalid shipping option ID: %v", err)
		return
	}

	err = cfg.db.DeleteShippingOption(r.Context(), id)
	if err != nil {
		cfg.RenderError(w, r, http.StatusInternalServerError, "Failed to delete shipping option")
		log.Printf("failed to delete shipping option: %v", err)
		return
	}
	http.Redirect(w, r, "/admin/shipping", http.StatusSeeOther)
}
