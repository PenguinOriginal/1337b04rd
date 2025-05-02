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

	"github.com/lib/pq"
)

// Injecting PostgreSQL
type PostgresPostRepo struct {
	db     *sql.DB
	logger *slog.Logger
}

// Constructor
func NewPostgresPostRepo(db *sql.DB, logger *slog.Logger) *PostgresPostRepo {
	return &PostgresPostRepo{db: db, logger: logger}
}

func (r *PostgresPostRepo) CreatePost(ctx context.Context, post *model.Post) error {
	query := `
	INSERT INTO posts (post_id, session_id, user_name, post_title, post_content, image_urls, created_at, is_archived)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.ExecContext(ctx, query,
		post.PostID,
		post.SessionID,
		post.UserName,
		post.Title,
		post.Content,
		pq.Array(post.ImageURLs),
		post.CreatedAt,
		post.IsArchived,
	)

	if err != nil {
		return logger.ErrorWrapper("repository", "CreatePost", "insert into posts", err)
	}
	return nil
}

func (r *PostgresPostRepo) GetPostByID(ctx context.Context, id utils.UUID) (*model.Post, error) {

	var post model.Post
	query := `
	SELECT post_id, session_id, user_name, post_title, post_content, image_urls, created_at, is_archived
	FROM posts 
	WHERE post_id = $1
	`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&post.PostID,
		&post.SessionID,
		&post.UserName,
		&post.Title,
		&post.Content,
		&post.ImageURLs,
		&post.CreatedAt,
		&post.IsArchived,
	)

	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			// Post not found case
			return nil, model.ErrPostNotFound
		}

		return nil, logger.ErrorWrapper("repository", "GetPostByID", "select post by ID", err)
	}
	return &post, nil
}

// Pass "archived" value to retrieve either active or archived posts
func (r *PostgresPostRepo) GetAllPosts(ctx context.Context, archived bool) ([]*model.Post, error) {
	query := `
	SELECT post_id, session_id, user_name, post_title, post_content, image_urls, created_at, is_archived
	FROM posts 
	WHERE is_archived = $1
	ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, archived)
	if err != nil {
		return nil, logger.ErrorWrapper("repository", "GetAllPosts", "query all posts", err)
	}
	defer rows.Close()

	var posts []*model.Post
	// rows.Next moves cursos
	for rows.Next() {
		var post model.Post

		if err := rows.Scan(
			&post.PostID,
			&post.SessionID,
			&post.UserName,
			&post.Title,
			&post.Content,
			pq.Array((&post.ImageURLs)),
			&post.CreatedAt,
			&post.IsArchived,
		); err != nil {
			return nil, logger.ErrorWrapper("repository", "GetAllPosts", "scan post row", err)
		}
		posts = append(posts, &post)
	}

	if err = rows.Err(); err != nil {
		return nil, logger.ErrorWrapper("repository", "GetAllPosts", "rows iteration", err)
	}
	return posts, nil
}

func (r *PostgresPostRepo) ArchivePostTx(ctx context.Context, tx *sql.Tx, postID utils.UUID) error {
	query := `
	UPDATE posts 
	SET is_archived = true 
	WHERE post_id = $1
	`
	// Execution of the query
	result, err := tx.ExecContext(ctx, query, postID)
	if err != nil {
		return logger.ErrorWrapper("repository", "ArchivePosts", "update is_archived", err)
	}

	// To check whether we changed rows
	affected, err := result.RowsAffected()
	if err != nil {
		return logger.ErrorWrapper("repository", "ArchivePosts", "rows affected", err)
	}

	if affected == 0 {
		// If we did not update any row, it means there is no such post
		return model.ErrPostNotFound
	}

	return nil
}

// Need this to update username during current session
func (r *PostgresPostRepo) UpdateUserNameForSession(ctx context.Context, sessionID utils.UUID, newName string) error {
	query := `
	UPDATE posts 
	SET user_name = $1 
	WHERE session_id = $2
	`

	// Execution of the query
	result, err := r.db.ExecContext(ctx, query, newName, sessionID)
	if err != nil {
		return logger.ErrorWrapper("repository", "UpdateUserNameForSession", "updating new name in post", err)
	}

	// To check whether we changed rows
	affected, err := result.RowsAffected()
	if err != nil {
		return logger.ErrorWrapper("repository", "UpdateUserNameForSession", "checking rows affected for post", err)
	}

	if affected == 0 {
		// If we did not update any row, it means there is no such post
		return model.ErrPostNotFound
	}
	return nil
}
