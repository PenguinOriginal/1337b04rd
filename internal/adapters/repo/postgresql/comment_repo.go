package postgresql

import (
	"1337b04rd/internal/domain/model"
	"1337b04rd/logger"
	"1337b04rd/pkg/utils"
	"context"
	"database/sql"
	"log/slog"
	"time"
)

// Injecting PostgreSQL
type PostgresCommentRepo struct {
	db     *sql.DB
	logger *slog.Logger
}

// Constructor
func NewPostgresCommentRepo(db *sql.DB, logger *slog.Logger) *PostgresCommentRepo {
	return &PostgresCommentRepo{db: db, logger: logger}
}

func (r *PostgresCommentRepo) CreateComment(ctx context.Context, comment *model.Comment) error {
	query := `
		INSERT INTO comments (
			comment_id, post_id, session_id, comment_content, parent_comment_id, created_at, is_archived
		) VALUES ($1, $2, $3, $4, $5, $6, $7);
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		comment.CommentID,
		comment.PostID,
		comment.SessionID,
		comment.Content,
		comment.ParentCommentID,
		comment.CreatedAt,
		comment.IsArchived,
	)

	if err != nil {
		r.logger.Error("failed to create comment", slog.Any("error", err))
		return logger.ErrorWrapper("repository", "CreateComment", "insert into comments", model.ErrDatabase)
	}

	return nil
}

func (r *PostgresCommentRepo) GetCommentByPostID(ctx context.Context, postID utils.UUID) ([]*model.Comment, error) {
	query := `
		SELECT comment_id, post_id, session_id, comment_content, parent_comment_id, created_at, is_archived
		FROM comments
		WHERE post_id = $1 AND is_archived = false
		ORDER BY created_at ASC;
	`

	rows, err := r.db.QueryContext(ctx, query, postID)
	if err != nil {
		r.logger.Error("failed to get comments", slog.Any("error", err))
		return nil, logger.ErrorWrapper("repository", "GetCommentByPostID", "select from comments", model.ErrDatabase)
	}
	defer rows.Close()

	var comments []*model.Comment

	for rows.Next() {
		var comment model.Comment
		err := rows.Scan(
			&comment.CommentID,
			&comment.PostID,
			&comment.SessionID,
			&comment.Content,
			&comment.ParentCommentID,
			&comment.CreatedAt,
			&comment.IsArchived,
		)
		if err != nil {
			r.logger.Error("failed to scan comment row", slog.Any("error", err))
			return nil, logger.ErrorWrapper("repository", "GetCommentByPostID", "row scan", model.ErrDatabase)
		}
		comments = append(comments, &comment)
	}

	if err = rows.Err(); err != nil {
		r.logger.Error("rows iteration error", slog.Any("error", err))
		return nil, logger.ErrorWrapper("repository", "GetCommentByPostID", "row iteration", model.ErrDatabase)
	}

	if len(comments) == 0 {
		return nil, model.ErrCommentNotFound
	}

	return comments, nil
}

// Why do we need this function?
func (r *PostgresCommentRepo) GetLatestCommentTime(ctx context.Context, postID utils.UUID) (*time.Time, error) {
	query := `
		SELECT MAX(created_at)
		FROM comments
		WHERE post_id = $1 AND is_archived = false;
	`

	var latestTime sql.NullTime
	err := r.db.QueryRowContext(ctx, query, postID).Scan(&latestTime)
	if err != nil {
		r.logger.Error("failed to get latest comment time", slog.Any("error", err))
		return nil, logger.ErrorWrapper("repository", "GetLatestCommentTime", "select MAX(created_at)", model.ErrDatabase)
	}

	if !latestTime.Valid {
		return nil, model.ErrCommentNotFound // No comments
	}

	return &latestTime.Time, nil
}

// ArchiveByPostID sets all comments on a post to archived
func (r *PostgresCommentRepo) ArchiveByPostID(ctx context.Context, postID utils.UUID) error {
	query := `
		UPDATE comments
		SET is_archived = true
		WHERE post_id = $1;
	`

	result, err := r.db.ExecContext(ctx, query, postID)
	if err != nil {
		r.logger.Error("failed to archive comments", slog.Any("error", err))
		return logger.ErrorWrapper("repository", "ArchiveByPostID", "update comments", model.ErrDatabase)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return model.ErrCommentNotFound
	}

	return nil
}
