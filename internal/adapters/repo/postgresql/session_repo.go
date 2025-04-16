// Revise this later
package postgresql

import (
	"1337b04rd/internal/domain/model"
	"1337b04rd/logger"
	"1337b04rd/pkg/utils"
	"context"
	"database/sql"
	"log/slog"
)

type PostgresSessionRepo struct {
	db     *sql.DB
	logger *slog.Logger
}

// Constructor
func NewPostgresSessionRepo(db *sql.DB, logger *slog.Logger) *PostgresSessionRepo {
	return &PostgresSessionRepo{db: db, logger: logger}
}

// Use it on the first POST request to create a post
func (r *PostgresSessionRepo) CreateSession(ctx context.Context, session *model.Session) error {
	const query = `
		INSERT INTO sessions (session_id, avatar_url, created_at, expires_at)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.ExecContext(ctx, query,
		session.SessionID,
		session.AvatarURL,
		session.CreatedAt,
		session.ExpiresAt,
	)

	if err != nil {
		r.logger.Error("failed to create session", slog.Any("error", err))
		return logger.ErrorWrapper("repository", "CreateSession", "insert into sessions", model.ErrDatabase)
	}

	return nil
}

// To identify returning user by their session ID
func (r *PostgresSessionRepo) GetSessionByID(ctx context.Context, id utils.UUID) (*model.Session, error) {
	const query = `
		SELECT session_id, avatar_url, created_at, expires_at
		FROM sessions
		WHERE session_id = $1
	`

	row := r.db.QueryRowContext(ctx, query, id)

	var session model.Session
	err := row.Scan(
		&session.SessionID,
		&session.AvatarURL,
		&session.CreatedAt,
		&session.ExpiresAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Warn("session not found", slog.String("session_id", string(id)))
			return nil, model.ErrSessionNotFound
		}
		r.logger.Error("failed to fetch session", slog.Any("error", err))
		return nil, logger.ErrorWrapper("repository", "GetSessionByID", "select from sessions", model.ErrDatabase)
	}

	return &session, nil
}

func (r *PostgresSessionRepo) DeleteExpiredSession(ctx context.Context) error {
	const query = `
		DELETE FROM sessions
		WHERE expires_at < CURRENT_TIMESTAMP
	`

	_, err := r.db.ExecContext(ctx, query)

	if err != nil {
		r.logger.Error("failed to delete expired sessions", slog.Any("error", err))
		return logger.ErrorWrapper("repository", "DeleteExpiredSession", "delete expired sessions", model.ErrDatabase)
	}

	return nil
}
