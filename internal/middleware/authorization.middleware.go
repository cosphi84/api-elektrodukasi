package middleware

import (
	"encoding/json"
	"elektrod/internal/auth"
	"net/http"
)

// ResourceType defines which resource is being accessed
type ResourceType string

const (
	ResourceUser     ResourceType = "user"
	ResourceCategory ResourceType = "category"
	ResourceTag      ResourceType = "tag"
	ResourceArticle  ResourceType = "article"
	ResourceComment  ResourceType = "comment"
	ResourceProject  ResourceType = "project"
)

// ActionType defines the action being performed
type ActionType string

const (
	ActionCreate ActionType = "create"
	ActionRead   ActionType = "read"
	ActionUpdate ActionType = "update"
	ActionDelete ActionType = "delete"
)

// AuthorizationMiddleware checks if user has permission for specific resource and action
func AuthorizationMiddleware(resource ResourceType, action ActionType) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			claims := GetUserClaims(r)
			if claims == nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]string{"error": "authentication required"})
				return
			}

			// Admins can do everything
			if auth.IsAdmin(claims.Role) {
				next(w, r)
				return
			}

			// Regular users can only POST/PUT/PATCH comments
			if auth.IsUser(claims.Role) {
				// Users can only write (POST, PUT, PATCH, DELETE) to comments
				if resource == ResourceComment && (action == ActionCreate || action == ActionUpdate || action == ActionDelete) {
					next(w, r)
					return
				}

				// Users can read all resources
				if action == ActionRead {
					next(w, r)
					return
				}
			}

			// Otherwise deny
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{"error": "insufficient permissions"})
		}
	}
}

// InferAuthorizationFromMethod creates authorization middleware based on HTTP method
func InferAuthorizationFromMethod(resource ResourceType) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			claims := GetUserClaims(r)
			if claims == nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]string{"error": "authentication required"})
				return
			}

			// Admins can do everything
			if auth.IsAdmin(claims.Role) {
				next(w, r)
				return
			}

			// Regular users rules
			if auth.IsUser(claims.Role) {
				method := r.Method

				// Users can READ (GET)
				if method == http.MethodGet {
					next(w, r)
					return
				}

				// Users can WRITE (POST, PUT, PATCH, DELETE) only to COMMENTS
				if resource == ResourceComment && (method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch || method == http.MethodDelete) {
					next(w, r)
					return
				}
			}

			// Otherwise deny
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{"error": "insufficient permissions for this operation"})
		}
	}
}
