package main

import (
	"net/http"

	"github.com/bzelaznicki/bzCommerce/internal/database"
)

type CheckoutPageData struct {
	Items []database.GetCartDetailsWithSnapshotPriceRow
	Total float64
}

func (cfg *apiConfig) handleViewCheckout(w http.ResponseWriter, r *http.Request) {
	cartID, err := cfg.getOrCreateCartID(w, r)

	if err != nil {
		cfg.Render(w, r, "templates/pages/cart.html", struct {
			Items []database.GetCartDetailsWithSnapshotPriceRow
			Total float64
		}{})
		return
	}

	items, err := cfg.db.GetCartDetailsWithSnapshotPrice(r.Context(), cartID)
	if err != nil {
		cfg.RenderError(w, r, http.StatusInternalServerError, "Could not load cart")
		return
	}

	var total float64
	for _, item := range items {
		total += float64(item.Quantity) * item.PricePerItem
	}

	data := CheckoutPageData{
		Items: items,
		Total: total,
	}
	cfg.Render(w, r, "templates/pages/checkout.html", data)

}
