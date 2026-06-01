package middleware

import (
	"context"
	"encoding/json"
	"elektrod/internal/auth"
	"net/http"
)

// ContextKey for storing user claims in request context
type ContextKey string

const UserClaimsKey ContextKey = "user_claims"

// AuthMiddleware validates JWT token from Authorization header
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		token, err := auth.ExtractTokenFromHeader(authHeader)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		claims, err := auth.ValidateToken(token)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid or expired token"})
			return
		}

		// Store claims in request context
		ctx := context.WithValue(r.Context(), UserClaimsKey, claims)
		next(w, r.WithContext(ctx))
	}
}

// GetUserClaims extracts user claims from request context
func GetUserClaims(r *http.Request) *auth.Claims {
	claims, ok := r.Context().Value(UserClaimsKey).(*auth.Claims)
	if !ok {
		return nil
	}
	return claims
}

// RequireAdmin middleware ensures only admin users can access the endpoint
func RequireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := GetUserClaims(r)
		if claims == nil || !auth.IsAdmin(claims.Role) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{"error": "admin access required"})
			return
		}
		next(w, r)
	}
}

// RequireUser middleware ensures authenticated user (admin or user role)
func RequireUser(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := GetUserClaims(r)
		if claims == nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "authentication required"})
			return
		}
		next(w, r)
	}
}

// Chain multiple middleware functions
func Chain(handler http.HandlerFunc, middleware ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for i := len(middleware) - 1; i >= 0; i-- {
		handler = middleware[i](handler)
	}
	return handler
}
