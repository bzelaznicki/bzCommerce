package main

import (
	"log"
	"net/http"

	"github.com/golang-jwt/jwt"
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

		userID, ok := claims["user_id"].(string)
		if !ok {
			cfg.RenderError(w, r, http.StatusUnauthorized, "Invalid user ID")
			log.Printf("invalid user ID: %v", err)
			return
		}

		ctx := context.WithValue(r.Context(), userIDContextKey, userID)
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
					if userID, ok := claims["user_id"].(string); ok {
						ctx := context.WithValue(r.Context(), userIDContextKey, userID)
						r = r.WithContext(ctx)
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
				// Already logged in — redirect away
				http.Redirect(w, r, "/account", http.StatusSeeOther)
				return
			}
		}

		// Continue to login/register handler
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

		if userID, ok := claims["user_id"].(string); ok {
			ctx := context.WithValue(r.Context(), userIDContextKey, userID)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}
