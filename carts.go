package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/bzelaznicki/bzCommerce/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handleAddToCart(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		cfg.RenderError(w, r, http.StatusBadRequest, "Invalid form")
		return
	}

	variantID, err := uuid.Parse(r.FormValue("variant_id"))
	if err != nil {
		cfg.RenderError(w, r, http.StatusBadRequest, "Invalid variant ID")
		return
	}

	quantity, err := strconv.Atoi(r.FormValue("quantity"))
	if err != nil || quantity <= 0 {
		cfg.RenderError(w, r, http.StatusBadRequest, "Invalid quantity")
		return
	}

	// Get or create cart only when adding item
	cartID, ok := getCartIDFromCookie(r)
	if !ok || cartID == uuid.Nil {

		// Create new cart
		userID := getUserIDFromContext(r.Context()) // optional
		var nullableUser uuid.NullUUID
		if userID != uuid.Nil {
			nullableUser = uuid.NullUUID{UUID: userID, Valid: true}
		}
		newCart, err := cfg.db.CreateCart(r.Context(), nullableUser)
		if err != nil {
			cfg.RenderError(w, r, http.StatusInternalServerError, "Could not create cart")
			return
		}
		cartID = newCart.ID
		// Set cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "cart_id",
			Value:    cartID.String(),
			Path:     "/",
			MaxAge:   30 * 24 * 60 * 60,
			HttpOnly: true,
			Secure:   false, // set to true in prod
		})
	}
	v, err := cfg.db.GetVariantByID(r.Context(), variantID)
	if err != nil {
		cfg.RenderError(w, r, http.StatusNotFound, "Variant not found")
		return
	}

	pricePerItem := v.Price

	if err != nil {
		cfg.RenderError(w, r, http.StatusInternalServerError, "Error parsing price")
	}
	// Add variant to cart
	variant := database.UpsertVariantToCartParams{
		CartID:           cartID,
		ProductVariantID: variantID,
		Quantity:         int32(quantity),
		PricePerItem:     pricePerItem,
	}
	_, err = cfg.db.UpsertVariantToCart(r.Context(), variant)
	if err != nil {
		cfg.RenderError(w, r, http.StatusInternalServerError, "Could not add to cart")
		return
	}

	http.Redirect(w, r, "/cart", http.StatusSeeOther)
}

func getCartIDFromCookie(r *http.Request) (uuid.UUID, bool) {
	cookie, err := r.Cookie("cart_id")
	if err != nil {
		return uuid.Nil, false
	}
	id, err := uuid.Parse(cookie.Value)
	if err != nil {
		return uuid.Nil, false
	}
	return id, true
}

func setCartIDCookie(w http.ResponseWriter, cartID uuid.UUID) {
	http.SetCookie(w, &http.Cookie{
		Name:     "cart_id",
		Value:    cartID.String(),
		Path:     "/",
		HttpOnly: true,
		MaxAge:   60 * 60 * 24 * 30, // 30 days
	})
}

func (cfg *apiConfig) getOrCreateCartID(w http.ResponseWriter, r *http.Request) (uuid.UUID, error) {
	ctx := r.Context()
	userID := getUserIDFromContext(ctx)

	// Check if cart_id cookie exists
	cartID, ok := getCartIDFromCookie(r)
	if ok {
		cart, err := cfg.db.GetCartById(ctx, cartID)
		if err == nil {
			// Cart exists and is valid
			return cart.ID, nil
		}
	}

	// If user is logged in, check their most recent active cart
	if userID != uuid.Nil {
		carts, err := cfg.db.GetCartsByUserId(ctx, uuid.NullUUID{
			UUID:  userID,
			Valid: true,
		})
		if err == nil {
			for _, cart := range carts {
				if cart.Status == "new" {
					// Found active cart
					setCartIDCookie(w, cart.ID)
					return cart.ID, nil
				}
			}
		}
	}

	// Create new cart
	newCart, err := cfg.db.CreateCart(ctx, uuid.NullUUID{
		UUID:  userID,
		Valid: userID != uuid.Nil,
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create cart: %w", err)
	}

	setCartIDCookie(w, newCart.ID)
	return newCart.ID, nil
}

func (cfg *apiConfig) handleViewCart(w http.ResponseWriter, r *http.Request) {
	cartID, ok := getCartIDFromCookie(r)
	if !ok {
		cfg.Render(w, r, "templates/pages/cart.html", nil)
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

	data := struct {
		Items []database.GetCartDetailsWithSnapshotPriceRow
		Total float64
	}{
		Items: items,
		Total: total,
	}

	cfg.Render(w, r, "templates/pages/cart.html", data)
}
