package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/bzelaznicki/bzCommerce/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handleAdminPaymentOptionsList(w http.ResponseWriter, r *http.Request) {
	rows, err := cfg.db.GetPaymentOptions(r.Context())

	if err != nil {
		cfg.RenderError(w, r, http.StatusInternalServerError, "Could not load shipping options")
		log.Printf("could not load shipping options: %v", err)
		return
	}

	type AdminPaymentListPageData struct {
		PaymentOptions []database.PaymentOption
		Breadcrumbs    []Breadcrumb
	}

	data := AdminPaymentListPageData{
		PaymentOptions: rows,
		Breadcrumbs: NewBreadcrumbTrail(
			Breadcrumb{Label: "Payment Options"},
		),
	}

	cfg.Render(w, r, "templates/pages/admin_payment_list.html", data)
}

type AdminPaymentFormPageData struct {
	PaymentOption database.PaymentOption
	IsEdit        bool
	Breadcrumbs   []Breadcrumb
}

func (cfg *apiConfig) handleAdminPaymentOptionsNew(w http.ResponseWriter, r *http.Request) {
	data := AdminPaymentFormPageData{
		PaymentOption: database.PaymentOption{},
		IsEdit:        false,
		Breadcrumbs: NewBreadcrumbTrail(
			Breadcrumb{Label: "Payment Options", URL: "/admin/payment"},
			Breadcrumb{Label: "New"},
		),
	}

	cfg.Render(w, r, "templates/pages/admin_payment_form.html", data)
}

func (cfg *apiConfig) handleAdminPaymentOptionsCreate(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		cfg.RenderError(w, r, http.StatusBadRequest, "Invalid form")
		log.Printf("error parsing form: %v", err)
		return
	}

	name := r.FormValue("name")
	description := r.FormValue("description")
	isActive := r.FormValue("is_active") != ""
	sortOrderStr := r.FormValue("sort_order")
	var sortOrder sql.NullInt32
	if sortOrderStr != "" {
		sortOrderInt64, err := strconv.ParseInt(sortOrderStr, 10, 32)
		if err != nil {
			cfg.RenderError(w, r, http.StatusBadRequest, "Invalid sort order")
			log.Printf("invalid sort order %v: %v", sortOrderStr, err)
			return
		}
		sortOrder = sql.NullInt32{Int32: int32(sortOrderInt64), Valid: true}
	}

	_, err := cfg.db.CreatePaymentOption(r.Context(), database.CreatePaymentOptionParams{
		Name: name,
		Description: sql.NullString{
			String: description,
			Valid:  description != "",
		},
		SortOrder: sortOrder,
		IsActive:  isActive,
	})

	if err != nil {
		cfg.RenderError(w, r, http.StatusInternalServerError, "Failed to create payment option")
		log.Printf("failed to create payment option: %v", err)
		return
	}
	http.Redirect(w, r, "/admin/payment", http.StatusSeeOther)
}

func (cfg *apiConfig) handleAdminPaymentOptionEdit(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))

	if err != nil {
		cfg.RenderError(w, r, http.StatusBadRequest, "Invalid payment option ID")
		log.Printf("invalid payment option ID: %v", err)
		return
	}

	paymentOption, err := cfg.db.GetPaymentOptionById(r.Context(), id)

	if err != nil {
		cfg.RenderError(w, r, http.StatusBadRequest, "Failed to load payment option")
		log.Printf("failed to load payment option ID %v: %v", id, err)
		return
	}

	data := AdminPaymentFormPageData{
		PaymentOption: paymentOption,
		IsEdit:        true,
		Breadcrumbs: NewBreadcrumbTrail(
			Breadcrumb{Label: "Payment Options", URL: "/admin/payment"},
			Breadcrumb{Label: "New"},
		),
	}

	cfg.Render(w, r, "templates/pages/admin_payment_form.html", data)
}

func (cfg *apiConfig) handleAdminPaymentOptionUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))

	if err != nil {
		cfg.RenderError(w, r, http.StatusBadRequest, "Invalid payment option ID")
		log.Printf("invalid payment option ID: %v", err)
		return
	}

	if err := r.ParseForm(); err != nil {
		cfg.RenderError(w, r, http.StatusBadRequest, "Invalid form")
		log.Printf("error parsing form: %v", err)
		return
	}

	name := r.FormValue("name")
	description := r.FormValue("description")
	isActive := r.FormValue("is_active") != ""
	sortOrderStr := r.FormValue("sort_order")
	var sortOrder sql.NullInt32
	if sortOrderStr != "" {
		sortOrderInt64, err := strconv.ParseInt(sortOrderStr, 10, 32)
		if err != nil {
			cfg.RenderError(w, r, http.StatusBadRequest, "Invalid sort order")
			log.Printf("invalid sort order %v: %v", sortOrderStr, err)
			return
		}
		sortOrder = sql.NullInt32{Int32: int32(sortOrderInt64), Valid: true}
	}

	err = cfg.db.UpdatePaymentOption(r.Context(), database.UpdatePaymentOptionParams{
		ID:   id,
		Name: name,
		Description: sql.NullString{
			String: description,
			Valid:  description != "",
		},
		SortOrder: sortOrder,
		IsActive:  isActive,
	})

	if err != nil {
		cfg.RenderError(w, r, http.StatusBadRequest, "Failed to update payment option")
		log.Printf("failed to update payment option: %v", err)
		return
	}

	http.Redirect(w, r, "/admin/payment", http.StatusSeeOther)
}

func (cfg *apiConfig) handleAdminPaymentOptionDelete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))

	if err != nil {
		cfg.RenderError(w, r, http.StatusBadRequest, "Invalid payment option ID")
		log.Printf("invalid payment option ID: %v", err)
		return
	}

	err = cfg.db.DeletePaymentOption(r.Context(), id)
	if err != nil {
		cfg.RenderError(w, r, http.StatusInternalServerError, "Failed to delete payment option")
		log.Printf("failed to delete payment option: %v", err)
		return
	}
	http.Redirect(w, r, "/admin/payment", http.StatusSeeOther)

}
