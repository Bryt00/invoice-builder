package contextutil

import (
	"context"
	"net/http"

	"raven.go.invoice-builder/internal/models"
)

// contextKey is an unexported type for context keys to avoid collisions.
type contextKey string

// AuthenticatedUserKey is the context key used to store the authenticated user
// across both web (session-based) and API (JWT-based) middleware layers.
const AuthenticatedUserKey = contextKey("authenticatedUser")

// ContextSetUser sets the authenticated user in the request context.
func ContextSetUser(r *http.Request, user *models.User) *http.Request {
	ctx := context.WithValue(r.Context(), AuthenticatedUserKey, user)
	return r.WithContext(ctx)
}

// ContextGetUser retrieves the authenticated user from the request context.
func ContextGetUser(r *http.Request) *models.User {
	user, ok := r.Context().Value(AuthenticatedUserKey).(*models.User)
	if !ok {
		return nil
	}
	return user
}
