package main

import (
	"github.com/bzelaznicki/bzCommerce/internal/database"
	"github.com/google/uuid"
)

func calculateCartTotal(cartId uuid.UUID, cartItems []database.GetCartDetailsWithSnapshotPriceRow, shippingFee float64) CartResponse {
	subtotal := 0.0
	for _, item := range cartItems {
		subtotal += float64(item.Quantity) * item.PricePerItem
	}
	itemCount := len(cartItems)
	total := subtotal + shippingFee

	return CartResponse{
		CartID:      cartId,
		ItemCount:   itemCount,
		Items:       cartItems,
		Subtotal:    subtotal,
		ShippingFee: shippingFee,
		Total:       total,
	}
}
