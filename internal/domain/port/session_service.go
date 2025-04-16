package port

import (
	"1337b04rd/internal/domain/model"
	"1337b04rd/pkg/utils"
	"context"
)

type SessionService interface {
	// CreateSession stores a new session (based on cookie logic).
	// Called when a new user lands and doesn't have a session yet.
	CreateSession(ctx context.Context, session *model.Session) error

	// GetSessionByID retrieves a session by its UUID.
	// Used to check who is making a request and fetch avatar info.
	GetSessionByID(ctx context.Context, id utils.UUID) (*model.Session, error)

	// DeleteExpiredSessions cleans up all expired sessions.
	// Useful for avoiding clutter or memory bloat; maybe run as a cron job.
	DeleteExpiredSessions(ctx context.Context) error
}
