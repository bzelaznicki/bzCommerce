package main

import (
	"log"
	"net/http"

	"github.com/google/uuid"
)

type UserAddressListPageData struct {
	Addresses   []UserAddressRow
	Breadcrumbs []Breadcrumb
}

type UserAddressRow struct {
	ID                uuid.UUID
	Name              string
	Address           string
	City              string
	PostalCode        string
	Country           string
	Phone             string
	IsShipping        bool
	IsShippingDefault bool
	IsBilling         bool
	IsBillingDefault  bool
}

func (cfg *apiConfig) handleAccountPage(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDContextKey).(string)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	user, err := cfg.db.GetUserById(r.Context(), uuid.MustParse(userID))
	if err != nil {
		cfg.RenderError(w, r, http.StatusNotFound, "User not found")
		log.Printf("user not found: %v", err)
		return
	}

	cfg.Render(w, r, "templates/pages/account.html", user)
}

func (cfg *apiConfig) handleManageAddressesList(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDContextKey).(string)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
	userUUID, err := uuid.Parse(userID)

	if err != nil {
		cfg.RenderError(w, r, http.StatusInternalServerError, "Failed to parse UUID")
		log.Printf("Failed to parse UUID: %v", err)
	}
	rows, err := cfg.db.GetUserAddresses(r.Context(), userUUID)

	if err != nil {
		cfg.RenderError(w, r, http.StatusInternalServerError, "Failed to get addresses")
	}

	addresses := make([]UserAddressRow, 0, len(rows))

	for _, row := range rows {
		addresses = append(addresses, UserAddressRow{
			ID:                row.ID,
			Name:              row.Name,
			Address:           row.Address.String,
			City:              row.City.String,
			PostalCode:        row.PostalCode.String,
			Country:           row.Country,
			IsShipping:        row.IsShipping,
			IsShippingDefault: row.IsShippingDefault,
			IsBilling:         row.IsBilling,
			IsBillingDefault:  row.IsBillingDefault,
		})
	}

	data := UserAddressListPageData{
		Addresses: addresses,
		Breadcrumbs: NewBreadcrumbTrail(
			Breadcrumb{Label: "Account", URL: "/account"},
		),
	}

	cfg.Render(w, r, "templates/pages/account_addresses_list.html", data)
}
