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

	quantityInt, err := strconv.ParseInt(r.FormValue("quantity"), 10, 32)
	if err != nil || quantityInt <= 0 {
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
		Quantity:         int32(quantityInt),
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

func (cfg *apiConfig) setCartIDCookie(w http.ResponseWriter, cartID uuid.UUID, secret []byte) {
	cartIDStr := cartID.String()
	sig := auth.SignCartID(cartIDStr, secret)

	http.SetCookie(w, &http.Cookie{
		Name:     "cart_id",
		Value:    fmt.Sprintf("%s|%s", cartIDStr, sig),
		Path:     "/",
		MaxAge:   60 * 60 * 24 * 30,
		HttpOnly: true,
		Secure:   cfg.platform != "dev",
	})
}

func (cfg *apiConfig) getOrCreateCartID(w http.ResponseWriter, r *http.Request) (uuid.UUID, error) {
	ctx := r.Context()
	userID := getUserIDFromContext(ctx)

	cartID, ok := getCartIDFromCookie(r, cfg.cartCookieKey)
	if ok && cartID != uuid.Nil {
		cart, err := cfg.db.GetCartById(ctx, cartID)
		if err == nil && cart.UserID.UUID == userID {
			return cart.ID, nil
		}

		clearCartIDCookie(w)
	}

	if userID != uuid.Nil {
		carts, err := cfg.db.GetCartsByUserId(ctx, uuid.NullUUID{
			UUID:  userID,
			Valid: true,
		})
		if err == nil {
			for _, cart := range carts {
				if cart.Status == "new" {
					cfg.setCartIDCookie(w, cart.ID, cfg.cartCookieKey)
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

	cfg.setCartIDCookie(w, newCart.ID, cfg.cartCookieKey)
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

func (cfg *apiConfig) handleDeleteCartItem(w http.ResponseWriter, r *http.Request) {
	cartID, err := cfg.getOrCreateCartID(w, r)
	if err != nil {
		cfg.RenderError(w, r, http.StatusInternalServerError, "Internal server error")
		log.Printf("error getting cart for user: %v", err)
		return
	}

	variantIdString := r.PathValue("id")

	variantId, err := uuid.Parse(variantIdString)

	if err != nil {
		cfg.RenderError(w, r, http.StatusNotFound, "invalid variant ID")
		log.Printf("error parsing ID %s on delete: %v", variantIdString, err)
		return
	}

	err = cfg.db.DeleteCartVariant(r.Context(), database.DeleteCartVariantParams{
		CartID:           cartID,
		ProductVariantID: variantId,
	})

	if err != nil {
		cfg.RenderError(w, r, http.StatusInternalServerError, "Internal server error")
		log.Printf("error deleting variant ID %v from cart ID %v: %v", variantId, cartID, err)
		return
	}

	http.Redirect(w, r, "/cart", http.StatusSeeOther)
}

func (cfg *apiConfig) handleCartUpdate(w http.ResponseWriter, r *http.Request) {
	cartID, err := cfg.getOrCreateCartID(w, r)
	if err != nil {
		cfg.RenderError(w, r, http.StatusInternalServerError, "Internal server error")
		log.Printf("error getting cart for user: %v", err)
		return
	}

	if err := r.ParseForm(); err != nil {
		cfg.RenderError(w, r, http.StatusBadRequest, "Invalid form data")
		return
	}

	for key, values := range r.PostForm {
		if strings.HasPrefix(key, "quantities[") {
			variantIDStr := strings.TrimSuffix(strings.TrimPrefix(key, "quantities["), "]")
			variantID, err := uuid.Parse(variantIDStr)
			if err != nil {
				log.Printf("invalid variant ID %s: %v", variantIDStr, err)
				continue
			}

			quantity64, err := strconv.ParseInt(values[0], 10, 32)
			if err != nil {
				log.Printf("invalid quantity for %s: %v", variantIDStr, err)
				continue
			}
			quantity := int32(quantity64)

			if quantity <= 0 {
				if err := cfg.db.DeleteCartVariant(r.Context(), database.DeleteCartVariantParams{
					CartID:           cartID,
					ProductVariantID: variantID,
				}); err != nil {
					log.Printf("error deleting variant ID %v from cart: %v", variantID, err)
				}
			} else {
				if int(quantity) > cfg.maxCartQuantity {
					cfg.RenderError(w, r, http.StatusBadRequest, "Invalid number of variants")
					return
				}
				_, err := cfg.db.UpdateCartVariantQuantity(r.Context(), database.UpdateCartVariantQuantityParams{
					CartID:           cartID,
					ProductVariantID: variantID,
					Quantity:         quantity,
				})
				if err != nil {
					log.Printf("error updating variant ID %v in cart: %v", variantID, err)
				}
			}
		}
	}

	http.Redirect(w, r, "/cart", http.StatusSeeOther)
}
