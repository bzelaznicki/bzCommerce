package main

import (
	"encoding/json"
	"net/http"

	"github.com/bzelaznicki/bzCommerce/internal/database"
	"github.com/google/uuid"
)

type CartResponse struct {
	CartID      uuid.UUID                                     `json:"cart_id"`
	ItemCount   int                                           `json:"item_count"`
	Items       []database.GetCartDetailsWithSnapshotPriceRow `json:"items"`
	Subtotal    float64                                       `json:"subtotal"`
	ShippingFee float64                                       `json:"shipping"`
	Total       float64                                       `json:"total"`
}

func (cfg *apiConfig) handleApiAddToCart(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		VariantId string `json:"variant_id"`
		Quantity  int32  `json:"quantity"`
	}

	var params parameters
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not decode parameters")
		return
	}

	variantID, err := uuid.Parse(params.VariantId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid variant ID")
		return
	}

	if params.Quantity <= 0 {
		respondWithError(w, http.StatusBadRequest, "Invalid quantity")
		return
	}

	cartID, err := cfg.getOrCreateCartID(w, r)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not get or create cart")
		return
	}

	dbVariant, err := cfg.db.GetVariantByID(r.Context(), variantID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Variant not found")
		return
	}

	if dbVariant.StockQuantity < params.Quantity {
		respondWithError(w, http.StatusBadRequest, "Not enough stock available")
		return
	}

	_, err = cfg.db.UpsertVariantToCart(r.Context(), database.UpsertVariantToCartParams{
		CartID:           cartID,
		ProductVariantID: dbVariant.ID,
		Quantity:         params.Quantity,
		PricePerItem:     dbVariant.Price,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not add to cart")
		return
	}

	cartItems, err := cfg.db.GetCartDetailsWithSnapshotPrice(r.Context(), cartID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get cart items")
		return
	}

	resp := calculateCartTotal(cartID, cartItems, 0)

	respondWithJSON(w, http.StatusOK, resp)
}

func (cfg *apiConfig) handleApiGetCart(w http.ResponseWriter, r *http.Request) {
	cartID, err := cfg.getOrCreateCartID(w, r)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not get or create cart")
		return
	}
	userID := getUserIDFromContext(r.Context())

	cart, err := cfg.db.GetCartById(r.Context(), cartID)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not load cart")
		return
	}
	if cart.UserID.UUID != userID {
		clearCartIDCookie(w)
		response := CartResponse{
			ItemCount: 0,
			Items:     []database.GetCartDetailsWithSnapshotPriceRow{},
			Subtotal:  0.0,
		}
		respondWithJSON(w, http.StatusOK, response)
		return
	}

	items, err := cfg.db.GetCartDetailsWithSnapshotPrice(r.Context(), cartID)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not load cart")
		return
	}

	resp := calculateCartTotal(cartID, items, 0)

	respondWithJSON(w, http.StatusOK, resp)
}

func (cfg *apiConfig) handleApiDeleteFromCart(w http.ResponseWriter, r *http.Request) {
	cartID, err := cfg.getOrCreateCartID(w, r)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting cart")
		return
	}

	variantId, err := uuid.Parse(r.PathValue("id"))

	if err != nil {
		respondWithError(w, http.StatusNotFound, "Variant not found")
		return
	}

	err = cfg.db.DeleteCartVariant(r.Context(), database.DeleteCartVariantParams{
		CartID:           cartID,
		ProductVariantID: variantId,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error deleting cart item")
		return
	}

	cartItems, err := cfg.db.GetCartDetailsWithSnapshotPrice(r.Context(), cartID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get cart items")
		return
	}

	resp := calculateCartTotal(cartID, cartItems, 0)

	respondWithJSON(w, http.StatusOK, resp)
}
func (cfg *apiConfig) handleApiUpdateCartVariant(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		VariantId string `json:"variant_id"`
		Quantity  int32  `json:"quantity"`
	}

	var params parameters
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not decode parameters")
		return
	}

	variantID, err := uuid.Parse(params.VariantId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid variant ID")
		return
	}

	cartID, err := cfg.getOrCreateCartID(w, r)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not get or create cart")
		return
	}

	ctx := r.Context()

	if params.Quantity <= 0 {
		err = cfg.db.DeleteCartVariant(ctx, database.DeleteCartVariantParams{
			CartID:           cartID,
			ProductVariantID: variantID,
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error deleting cart item")
			return
		}
	} else {
		dbVariant, err := cfg.db.GetVariantByID(ctx, variantID)
		if err != nil {
			respondWithError(w, http.StatusNotFound, "Variant not found")
			return
		}

		_, err = cfg.db.UpdateCartVariant(ctx, database.UpdateCartVariantParams{
			CartID:           cartID,
			ProductVariantID: variantID,
			Quantity:         params.Quantity,
			PricePerItem:     dbVariant.Price,
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Could not update cart item")
			return
		}
	}

	cartItems, err := cfg.db.GetCartDetailsWithSnapshotPrice(ctx, cartID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get cart items")
		return
	}

	resp := calculateCartTotal(cartID, cartItems, 0)
	respondWithJSON(w, http.StatusOK, resp)
}
