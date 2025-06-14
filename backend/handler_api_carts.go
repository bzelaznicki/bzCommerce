package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bzelaznicki/bzCommerce/internal/database"
	"github.com/google/uuid"
)

type CartResponse struct {
	ItemCount int                                           `json:"item_count"`
	Items     []database.GetCartDetailsWithSnapshotPriceRow `json:"items"`
	Total     float64                                       `json:"total"`
}

func (cfg *apiConfig) handleApiAddToCart(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		VariantId string `json:"variant_id"`
		Quantity  int32  `json:"quantity"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not decode parameters")
		return
	}

	variantId, err := uuid.Parse(params.VariantId)

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
	userID := getUserIDFromContext(r.Context())
	fmt.Printf("User of the current user is: %v\n", userID)

	cart, err := cfg.db.GetCartById(r.Context(), cartID)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not load cart")
		return
	}

	//checking if the current user owns the cart - edge case
	if cart.UserID.Valid && cart.UserID.UUID != userID || cart.UserID.Valid && userID == uuid.Nil {
		fmt.Printf("This is not this user's cart, deleting cookie\n")
		clearCartIDCookie(w)

		// Create a new cart for the current user - subsequent operations will use this new cart ID
		cartID, err = cfg.getOrCreateCartID(w, r)
		fmt.Printf("new cart ID: %v", cartID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Could not load cart")
			return
		}
	}

	dbVariant, err := cfg.db.GetVariantByID(r.Context(), variantId)

	if err != nil {
		respondWithError(w, http.StatusNotFound, "Variant not found")
		return
	}

	if dbVariant.StockQuantity < params.Quantity {
		respondWithError(w, http.StatusBadRequest, "Not enough stock available")
		return
	}

	variant := database.UpsertVariantToCartParams{
		CartID:           cartID,
		ProductVariantID: dbVariant.ID,
		Quantity:         params.Quantity,
		PricePerItem:     dbVariant.Price,
	}

	_, err = cfg.db.UpsertVariantToCart(r.Context(), variant)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not add to cart")
		return
	}

	cartItems, err := cfg.db.GetCartDetailsWithSnapshotPrice(r.Context(), cartID)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get cart items")
		return
	}

	total := calculateCartTotal(cartItems)
	resp := CartResponse{
		ItemCount: len(cartItems),
		Items:     cartItems,
		Total:     total,
	}
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
	if cart.UserID.Valid && cart.UserID.UUID != userID {
		clearCartIDCookie(w)
		response := CartResponse{
			ItemCount: 0,
			Items:     []database.GetCartDetailsWithSnapshotPriceRow{},
			Total:     0.0,
		}
		respondWithJSON(w, http.StatusOK, response)
		return
	}

	items, err := cfg.db.GetCartDetailsWithSnapshotPrice(r.Context(), cartID)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not load cart")
		return
	}

	total := calculateCartTotal(items)

	response := CartResponse{
		ItemCount: len(items),
		Items:     items,
		Total:     total,
	}

	respondWithJSON(w, http.StatusOK, response)
}
