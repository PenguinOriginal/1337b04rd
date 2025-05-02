// Reviewed and double-checked
package postgresql

import (
	"1337b04rd/internal/domain/model"
	"1337b04rd/pkg/logger"
	"1337b04rd/pkg/utils"
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/lib/pq"
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
			comment_id, post_id, session_id, user_name, comment_content, parent_comment_id, image_urls, created_at, is_archived
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		comment.CommentID,
		comment.PostID,
		comment.SessionID,
		comment.UserName,
		comment.Content,
		comment.ParentCommentID,
		pq.Array(comment.ImageURLs),
		comment.CreatedAt,
		comment.IsArchived,
	)

	if err != nil {
		return logger.ErrorWrapper("repository", "CreateComment", "insert into comments", model.ErrDatabase)
	}

	return nil
}

func (r *PostgresCommentRepo) GetCommentsByPostID(ctx context.Context, postID utils.UUID, includeArchived bool) ([]*model.Comment, error) {
	query := `
		SELECT comment_id, post_id, session_id, user_name, comment_content, parent_comment_id, image_urls, created_at, is_archived
		FROM comments
		WHERE post_id = $1
	`

	// Return only active comments if it is the main page
	if !includeArchived {
		query += " AND is_archived = false"
	}
	query += " ORDER BY created_at DESC"

	rows, err := r.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, logger.ErrorWrapper("repository", "GetCommentsByPostID", "select from comments", model.ErrDatabase)
	}
	defer rows.Close()

	var comments []*model.Comment

	for rows.Next() {
		var comment model.Comment
		err := rows.Scan(
			&comment.CommentID,
			&comment.PostID,
			&comment.SessionID,
			&comment.UserName,
			&comment.Content,
			&comment.ParentCommentID,
			pq.Array(comment.ImageURLs),
			&comment.CreatedAt,
			&comment.IsArchived,
		)
		if err != nil {
			return nil, logger.ErrorWrapper("repository", "GetCommentsByPostID", "row scan", model.ErrDatabase)
		}
		comments = append(comments, &comment)
	}

	if err = rows.Err(); err != nil {
		return nil, logger.ErrorWrapper("repository", "GetCommentsByPostID", "row iteration", model.ErrDatabase)
	}

	if len(comments) == 0 {
		return nil, model.ErrCommentNotFound
	}

	return comments, nil
}

func (r *PostgresCommentRepo) GetCommentByID(ctx context.Context, commentID utils.UUID) (*model.Comment, error) {
	query := `
		SELECT comment_id, post_id, session_id, user_name, comment_content, parent_comment_id, image_urls, created_at, is_archived
		FROM comments
		WHERE comment_id = $1
	`

	var c model.Comment
	var parentCommentID sql.NullString
	var imageURLs []string

	err := r.db.QueryRowContext(ctx, query, commentID).Scan(
		&c.CommentID,
		&c.PostID,
		&c.SessionID,
		&c.UserName,
		&c.Content,
		&parentCommentID,
		pq.Array(&imageURLs),
		&c.CreatedAt,
		&c.IsArchived,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrCommentNotFound
		}
		return nil, logger.ErrorWrapper("repository", "GetCommentByID", "scanning result", err)
	}

	if parentCommentID.Valid {
		id := utils.UUID(parentCommentID.String)
		c.ParentCommentID = id
	} else {
		c.ParentCommentID = ""
	}
	c.ImageURLs = imageURLs

	return &c, nil
}

// Fetch the most recent comment's created_at
// Need it for archiving logic
func (r *PostgresCommentRepo) GetLatestCommentTime(ctx context.Context, postID utils.UUID) (*time.Time, error) {
	query := `
		SELECT MAX(created_at)
		FROM comments
		WHERE post_id = $1 AND is_archived = false
	`

	var latestTime sql.NullTime
	err := r.db.QueryRowContext(ctx, query, postID).Scan(&latestTime)
	if err != nil {
		return nil, logger.ErrorWrapper("repository", "GetLatestCommentTime", "select MAX(created_at)", model.ErrDatabase)
	}

	if !latestTime.Valid {
		return nil, model.ErrCommentNotFound // No comments
	}

	return &latestTime.Time, nil
}

// Make all comments of a specific post as is_archived = true
func (r *PostgresCommentRepo) ArchiveCommentByPostIDTx(ctx context.Context, tx *sql.Tx, postID utils.UUID) error {
	query := `
		UPDATE comments
		SET is_archived = true
		WHERE post_id = $1
	`

	result, err := tx.ExecContext(ctx, query, postID)
	if err != nil {
		return logger.ErrorWrapper("repository", "ArchiveCommentByPostID", "update comments", model.ErrDatabase)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return model.ErrCommentNotFound
	}

	return nil
}

func (r *PostgresCommentRepo) UpdateUserNameForSession(ctx context.Context, sessionID utils.UUID, newName string) error {
	query := `
	UPDATE comments 
	SET user_name = $1 
	WHERE session_id = $2
	`

	result, err := r.db.ExecContext(ctx, query, newName, sessionID)
	if err != nil {
		return logger.ErrorWrapper("repository", "UpdateUserNameForSession", "updating user name for comment", model.ErrDatabase)
	}

	// To check whether we changed rows
	affected, err := result.RowsAffected()
	if err != nil {
		return logger.ErrorWrapper("repository", "UpdateUserNameForSession", "checking rows affected for comment", err)
	}

	if affected == 0 {
		// If we did not update any row, it means there is no such comment
		return model.ErrCommentNotFound
	}
	return nil
}
