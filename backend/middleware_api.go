package main

import (
	"context"
	"net/http"

	"github.com/bzelaznicki/bzCommerce/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == cfg.frontendUrl {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

const (
	contextKeyUserID  contextKey = "userID"
	contextKeyIsAdmin contextKey = "isAdmin"
)

func (cfg *apiConfig) checkAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearerToken, err := auth.GetBearerToken(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Bearer Token missing or invalid")
			return
		}

		userID, isAdmin, err := auth.ValidateJWT(bearerToken, cfg.jwtSecret)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		userUUID, parseErr := uuid.Parse(userID)
		if parseErr != nil {
			respondWithError(w, http.StatusUnauthorized, "Invalid user ID")
			return
		}

		user, dbErr := cfg.db.GetUserById(r.Context(), userUUID)
		if dbErr != nil {
			respondWithError(w, http.StatusUnauthorized, "User not found")
			return
		}

		if !user.IsActive {
			respondWithError(w, http.StatusForbidden, "Your account has been disabled")
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyUserID, userID)
		ctx = context.WithValue(ctx, contextKeyIsAdmin, isAdmin)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (cfg *apiConfig) optionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearerToken, err := auth.GetBearerToken(r.Header)
		if err == nil && bearerToken != "" {
			userID, isAdmin, err := auth.ValidateJWT(bearerToken, cfg.jwtSecret)
			if err == nil {
				userUUID, parseErr := uuid.Parse(userID)
				if parseErr == nil {
					user, dbErr := cfg.db.GetUserById(r.Context(), userUUID)
					if dbErr == nil && user.IsActive {
						ctx := context.WithValue(r.Context(), contextKeyUserID, userID)
						ctx = context.WithValue(ctx, contextKeyIsAdmin, isAdmin)
						r = r.WithContext(ctx)
					}
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) checkAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isAdmin, ok := r.Context().Value(contextKeyIsAdmin).(bool)

		if !ok || !isAdmin {
			respondWithError(w, http.StatusUnauthorized, "Admin access required")
			return
		}
		next.ServeHTTP(w, r)
	})
}
