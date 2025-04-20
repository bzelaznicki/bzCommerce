package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/bzelaznicki/bzCommerce/internal/auth"
	"github.com/bzelaznicki/bzCommerce/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handleAddToCart(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		cfg.RenderError(w, r, http.StatusBadRequest, "Invalid form")
		return
	}

	variantIDStr := r.FormValue("variant_id")
	if variantIDStr == "" {
		cfg.RenderError(w, r, http.StatusBadRequest, "Please select a variant before adding to your cart.")
		return
	}

	variantID, err := uuid.Parse(variantIDStr)
	if err != nil {
		cfg.RenderError(w, r, http.StatusBadRequest, "Invalid variant ID.")
		return
	}

	quantity, err := strconv.Atoi(r.FormValue("quantity"))
	if err != nil || quantity <= 0 {
		cfg.RenderError(w, r, http.StatusBadRequest, "Invalid quantity")
		return
	}

	cartID, err := cfg.getOrCreateCartID(w, r)
	if err != nil {
		cfg.RenderError(w, r, http.StatusInternalServerError, "Could not get or create cart")
		return
	}

	v, err := cfg.db.GetVariantByID(r.Context(), variantID)
	if err != nil {
		cfg.RenderError(w, r, http.StatusNotFound, "Variant not found")
		return
	}

	variant := database.UpsertVariantToCartParams{
		CartID:           cartID,
		ProductVariantID: variantID,
		Quantity:         int32(quantity),
		PricePerItem:     v.Price,
	}
	_, err = cfg.db.UpsertVariantToCart(r.Context(), variant)
	if err != nil {
		cfg.RenderError(w, r, http.StatusInternalServerError, "Could not add to cart")
		return
	}

	http.Redirect(w, r, "/cart", http.StatusSeeOther)
}

func getCartIDFromCookie(r *http.Request, secret []byte) (uuid.UUID, bool) {
	cookie, err := r.Cookie("cart_id")
	if err != nil {
		return uuid.Nil, false
	}

	parts := strings.Split(cookie.Value, "|")
	if len(parts) != 2 {
		return uuid.Nil, false
	}

	cartIDStr, sig := parts[0], parts[1]

	if !auth.VerifyCartID(cartIDStr, sig, secret) {
		return uuid.Nil, false
	}

	cartID, err := uuid.Parse(cartIDStr)
	if err != nil {
		return uuid.Nil, false
	}

	return cartID, true
}

func setCartIDCookie(w http.ResponseWriter, cartID uuid.UUID, secret []byte) {
	cartIDStr := cartID.String()
	sig := auth.SignCartID(cartIDStr, secret)

	http.SetCookie(w, &http.Cookie{
		Name:     "cart_id",
		Value:    fmt.Sprintf("%s|%s", cartIDStr, sig),
		Path:     "/",
		MaxAge:   60 * 60 * 24 * 30, // 30 days
		HttpOnly: true,
		Secure:   false, // change to true in prod
	})
}

func (cfg *apiConfig) getOrCreateCartID(w http.ResponseWriter, r *http.Request) (uuid.UUID, error) {
	ctx := r.Context()
	userID := getUserIDFromContext(ctx)

	cartID, ok := getCartIDFromCookie(r, cfg.cartCookieKey)
	if ok && cartID != uuid.Nil {
		cart, err := cfg.db.GetCartById(ctx, cartID)
		if err == nil {
			if cart.UserID.Valid && userID != uuid.Nil && cart.UserID.UUID != userID {
				clearCartIDCookie(w)
			} else {
				return cart.ID, nil
			}
		}
	}

	if userID != uuid.Nil {
		carts, err := cfg.db.GetCartsByUserId(ctx, uuid.NullUUID{
			UUID:  userID,
			Valid: true,
		})
		if err == nil {
			for _, cart := range carts {
				if cart.Status == "new" {
					setCartIDCookie(w, cart.ID, cfg.cartCookieKey)
					return cart.ID, nil
				}
			}
		}
	}

	newCart, err := cfg.db.CreateCart(ctx, uuid.NullUUID{
		UUID:  userID,
		Valid: userID != uuid.Nil,
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create cart: %w", err)
	}

	setCartIDCookie(w, newCart.ID, cfg.cartCookieKey)
	return newCart.ID, nil
}

func clearCartIDCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "cart_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
}

func (cfg *apiConfig) mergeAnonymousCart(ctx context.Context, anonCartID, userID uuid.UUID) error {

	userCarts, err := cfg.db.GetCartsByUserId(ctx, uuid.NullUUID{UUID: userID, Valid: true})
	if err != nil {
		return fmt.Errorf("failed to fetch user carts: %w", err)
	}

	var userCartID uuid.UUID
	for _, cart := range userCarts {
		if cart.Status == "new" {
			userCartID = cart.ID
			break
		}
	}

	if userCartID == uuid.Nil {
		_, err := cfg.db.UpdateCartOwner(ctx, database.UpdateCartOwnerParams{
			ID:     anonCartID,
			UserID: uuid.NullUUID{UUID: userID, Valid: true},
		})
		if err != nil {
			return fmt.Errorf("failed to assign cart to user: %w", err)
		}
		return nil
	}

	items, err := cfg.db.GetCartProductsByCartId(ctx, anonCartID)
	if err != nil {
		return fmt.Errorf("failed to get anonymous cart items: %w", err)
	}

	for _, item := range items {
		_, err := cfg.db.UpsertVariantToCart(ctx, database.UpsertVariantToCartParams{
			CartID:           userCartID,
			ProductVariantID: item.ProductVariantID,
			Quantity:         item.Quantity,
			PricePerItem:     item.PricePerItem,
		})
		if err != nil {
			return fmt.Errorf("failed to upsert item during cart merge: %w", err)
		}
	}

	err = cfg.db.ClearCartItems(ctx, anonCartID)
	if err != nil {
		return fmt.Errorf("failed to clear anonymous cart after merge: %w", err)
	}
	err = cfg.db.DeleteCart(ctx, anonCartID)
	if err != nil {
		log.Printf("error deleting merged anonymous cart: %v", err)
	}

	return nil
}

func (cfg *apiConfig) getOrCreateCartIDForUser(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	if userID == uuid.Nil {
		return uuid.Nil, fmt.Errorf("invalid user ID")
	}

	carts, err := cfg.db.GetCartsByUserId(ctx, uuid.NullUUID{
		UUID:  userID,
		Valid: true,
	})
	if err == nil {
		for _, cart := range carts {
			if cart.Status == "new" {
				return cart.ID, nil
			}
		}
	}

	newCart, err := cfg.db.CreateCart(ctx, uuid.NullUUID{
		UUID:  userID,
		Valid: true,
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create cart: %w", err)
	}

	return newCart.ID, nil
}

func (cfg *apiConfig) handleViewCart(w http.ResponseWriter, r *http.Request) {
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

	data := struct {
		Items []database.GetCartDetailsWithSnapshotPriceRow
		Total float64
	}{
		Items: items,
		Total: total,
	}

	cfg.Render(w, r, "templates/pages/cart.html", data)
}
