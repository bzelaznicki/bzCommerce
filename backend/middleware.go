package main

import (
	"log"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/net/context"
)

type contextKey string

const userIDContextKey contextKey = "userID"

func (cfg *apiConfig) withAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("auth_token")
		if err != nil {
			log.Println("No auth_token cookie found")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		token, err := jwt.Parse(cookie.Value, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(cfg.jwtSecret), nil
		})

		if err != nil || !token.Valid {
			cfg.RenderError(w, r, http.StatusUnauthorized, "Invalid token")
			log.Printf("invalid token: %v", err)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			cfg.RenderError(w, r, http.StatusUnauthorized, "Invalid token claims")
			log.Printf("invalid token claims: %v", err)
			return
		}

		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			cfg.RenderError(w, r, http.StatusUnauthorized, "Invalid user ID")
			log.Printf("invalid user ID: %v", err)
			return
		}

		userID, parseErr := uuid.Parse(userIDStr)
		if parseErr != nil {
			cfg.RenderError(w, r, http.StatusUnauthorized, "Invalid user ID")
			return
		}

		user, dbErr := cfg.db.GetUserById(r.Context(), userID)
		if dbErr != nil || !user.IsActive {
			cfg.RenderError(w, r, http.StatusForbidden, "Your account has been disabled")
			return
		}

		ctx := context.WithValue(r.Context(), userIDContextKey, userIDStr)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (cfg *apiConfig) maybeWithAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("auth_token")
		if err == nil {
			token, err := jwt.Parse(cookie.Value, func(t *jwt.Token) (any, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(cfg.jwtSecret), nil
			})

			if err == nil && token.Valid {
				if claims, ok := token.Claims.(jwt.MapClaims); ok {
					if userIDStr, ok := claims["user_id"].(string); ok {
						userID, parseErr := uuid.Parse(userIDStr)
						if parseErr == nil {
							user, dbErr := cfg.db.GetUserById(r.Context(), userID)
							if dbErr == nil && user.IsActive {
								ctx := context.WithValue(r.Context(), userIDContextKey, userIDStr)
								r = r.WithContext(ctx)
							}
						}
					}
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) redirectIfAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("auth_token")
		if err == nil {
			token, err := jwt.Parse(cookie.Value, func(t *jwt.Token) (any, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(cfg.jwtSecret), nil
			})

			if err == nil && token.Valid {
				http.Redirect(w, r, "/account", http.StatusSeeOther)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) requireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("auth_token")
		if err != nil {
			cfg.RenderError(w, r, http.StatusUnauthorized, "Access denied")
			return
		}

		token, err := jwt.Parse(cookie.Value, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(cfg.jwtSecret), nil
		})
		if err != nil || !token.Valid {
			cfg.RenderError(w, r, http.StatusUnauthorized, "Access denied")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || claims["is_admin"] != true {
			cfg.RenderError(w, r, http.StatusForbidden, "Access denied")
			return
		}

		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			cfg.RenderError(w, r, http.StatusUnauthorized, "Invalid user ID")
			return
		}

		userID, parseErr := uuid.Parse(userIDStr)
		if parseErr != nil {
			cfg.RenderError(w, r, http.StatusUnauthorized, "Invalid user ID")
			return
		}

		user, dbErr := cfg.db.GetUserById(r.Context(), userID)
		if dbErr != nil || !user.IsActive {
			cfg.RenderError(w, r, http.StatusForbidden, "Your account has been disabled")
			return
		}

		ctx := context.WithValue(r.Context(), userIDContextKey, userIDStr)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
