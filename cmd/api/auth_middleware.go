package main

import (
	"net/http"
	"strings"

	"raven.go.invoice-builder/internal/contextutil"
	"raven.go.invoice-builder/internal/models"
)

func (h *ApiHandler) ContextGetUser(r *http.Request) *models.User {
	return contextutil.ContextGetUser(r)
}

func (h *ApiHandler) RequireJWTAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			h.AuthenticationRequiredResponse(w, r)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			h.AuthenticationRequiredResponse(w, r)
			return
		}

		tokenString := parts[1]
		claims, err := h.JWTManager.ValidateToken(tokenString)
		if err != nil {
			h.AuthenticationRequiredResponse(w, r)
			return
		}

		user, err := h.Services.Auth.GetUserByID(r.Context(), claims.UserID)
		if err != nil || user == nil || !user.IsActivated {
			h.AuthenticationRequiredResponse(w, r)
			return
		}

		// Verify token version matches — rejects tokens issued before
		// a password change, role change, or manual revocation.
		if claims.TokenVersion != user.TokenVersion {
			h.AuthenticationRequiredResponse(w, r)
			return
		}

		r = contextutil.ContextSetUser(r, user)
		next.ServeHTTP(w, r)
	})
}

// RequireRole returns middleware that restricts access to users with the given role.
func (h *ApiHandler) RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := h.ContextGetUser(r)
			if user == nil {
				h.AuthenticationRequiredResponse(w, r)
				return
			}
			if user.Role.Name != role && !user.IsAdmin() {
				h.ForbiddenResponse(w, r)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
