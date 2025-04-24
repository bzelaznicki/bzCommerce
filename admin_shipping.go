package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/bzelaznicki/bzCommerce/internal/database"
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
	}

	data := AdminShippingListPageData{
		ShippingOptions: rows,
	}
	cfg.Render(w, r, "templates/pages/admin_shipping_list.html", data)
}

type AdminShippingFormPageData struct {
	ShippingOption database.ShippingOption
	IsEdit         bool
}

func (cfg *apiConfig) handleAdminShippingOptionsNew(w http.ResponseWriter, r *http.Request) {
	data := AdminShippingFormPageData{
		ShippingOption: database.ShippingOption{},
		IsEdit:         false,
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
	isActive := r.FormValue("is_active") == "true"

	sortOrderInt, err := strconv.Atoi(sortOrder)

	if err != nil {
		cfg.RenderError(w, r, http.StatusBadRequest, "Invalid sort order")
		log.Printf("invalid sort order: %v", err)
		return
	}

	_, err = cfg.db.CreateShippingOption(r.Context(), database.CreateShippingOptionParams{
		Name: name,
		Description: sql.NullString{
			String: description,
			Valid:  true,
		},
		Price: price,
		EstimatedDays: sql.NullString{
			String: estimatedDays,
			Valid:  true,
		},
		SortOrder: sql.NullInt32{
			Int32: int32(sortOrderInt),
			Valid: true,
		},
		IsActive: isActive,
	})

	if err != nil {
		cfg.RenderError(w, r, http.StatusInternalServerError, "Failed to create shipping option")
		log.Printf("failed to create shipping option: %v", err)
		return
	}
	http.Redirect(w, r, "/admin/shipping", http.StatusSeeOther)
}
