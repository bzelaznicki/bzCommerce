package main

import (
	"context"
	"database/sql"
	"log"
	"time"
)

func (cfg *apiConfig) startCartExpirationWorker() {
	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		for range ticker.C {
			cfg.expireOldCarts()
		}
	}()
}

func (cfg *apiConfig) expireOldCarts() {
	threshold := time.Now().Add(-time.Duration(cfg.cartTimeoutMinutes) * time.Minute)
	err := cfg.db.MarkCartsAsAbandoned(context.Background(), sql.NullTime{Time: threshold, Valid: true})
	if err != nil {
		log.Printf("failed to mark abandoned carts: %v", err)
	}
	log.Printf("clearing abandoned carts completed")
}
