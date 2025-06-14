package main

import "github.com/bzelaznicki/bzCommerce/internal/database"

func calculateCartTotal(cartItems []database.GetCartDetailsWithSnapshotPriceRow) float64 {
	total := 0.0
	for _, item := range cartItems {
		total += float64(item.Quantity) * item.PricePerItem
	}

	return total
}
