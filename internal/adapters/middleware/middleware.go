// Разобрать позже
package middleware

import (
	"1337b04rd/internal/domain/model"
	"1337b04rd/internal/domain/port"
	"context"
	"net/http"

	"1337b04rd/pkg/utils"
)

const sessionCookieName = "session_id"

type sessionKeyType string

const sessionKey sessionKeyType = "user-session"

func SessionMiddleware(sessionService port.SessionService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(sessionCookieName)

			var session *model.Session

			// CASE 1: Cookie exists → fetch session from DB
			if err == nil {
				sessionID := utils.UUID(cookie.Value)

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
					Name:     sessionCookieName,
					Value:    string(session.SessionID),
					Path:     "/",
					HttpOnly: true,
					SameSite: http.SameSiteStrictMode,
					Secure:   false, // true in prod (HTTPS)
					Expires:  session.ExpiresAt,
				})
			}

			// Add session to context
			ctx := context.WithValue(r.Context(), sessionKey, session)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
