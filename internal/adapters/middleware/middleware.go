// Ensures that every visitor has valid session, and that session stored in the context
// Runs before handlers
package middleware

import (
	"1337b04rd/internal/domain/model"
	"1337b04rd/internal/domain/port"
	"context"
	"net/http"

	"1337b04rd/pkg/utils"
)

const sessionCookieName = "session_id"

// Custom key type used for storing values in context.Context
// To avoid name collision
type sessionKeyType string

const sessionKey sessionKeyType = "user-session"

// Struct because I can add more data later (e.g. UserName, Role, etc)
type SessionData struct {
	AvatarURL string
}

// Why do I need this and why can't I pass just avaratURL var for example?

// This returns middleware function that takes and wraps another http.Handler
func SessionMiddleware(sessionService port.SessionService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler { // standard go middleware signature
		// my logic BEFORE the handler
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Try to read existing cookie
			cookie, err := r.Cookie(sessionCookieName)

			var session *model.Session

			// CASE 1: Cookie exists → fetch session from DB
			if err == nil {
				// If cookie is present, convert it to UUID type
				sessionID := utils.UUID(cookie.Value)

				// Fetch session from the db
				s, err := sessionService.GetSessionByID(r.Context(), sessionID)
				if err == nil {
					session = s
				}
			}

			// CASE 2: No valid session → create one
			if session == nil {
				newSession := &model.Session{}
				err := sessionService.CreateSession(r.Context(), newSession)
				if err != nil {
					http.Error(w, "Failed to initialize session", http.StatusInternalServerError)
					return
				}
				session = newSession

				// Set cookie
				http.SetCookie(w, &http.Cookie{
					Name:     sessionCookieName,         // Cookie key
					Value:    string(session.SessionID), // Cookie value
					Path:     "/",
					HttpOnly: true,
					SameSite: http.SameSiteStrictMode,
					Secure:   false,
					Expires:  session.ExpiresAt,
				})
			}

			// Add session to context
			ctx := context.WithValue(r.Context(), sessionKey, session)
			next.ServeHTTP(w, r.WithContext(ctx))
			// my logic AFTER the handler
		})
	}
}

// Allows handlers to retrieve session
func GetSessionFromContext(ctx context.Context) *model.Session {
	val := ctx.Value(sessionKey)
	if session, ok := val.(*model.Session); ok {
		return session
	}
	return nil
}
