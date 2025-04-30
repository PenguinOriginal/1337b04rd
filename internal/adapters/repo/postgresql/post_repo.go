// DB adapter
// LATER: check if the mistakes are not dublicated across layers
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
		r.logger.Error("failed to create post", slog.Any("error", err))
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
			r.logger.Warn("post not found", slog.String("post_id", string(id)))
			return nil, model.ErrPostNotFound
		}

		r.logger.Error("failed to get post by ID", slog.Any("error", err), slog.Any("post_id", id))
		return nil, logger.ErrorWrapper("repository", "GetPostByID", "select post by ID", err)
	}
	return &post, nil
}

func (r *PostgresPostRepo) GetAllPosts(ctx context.Context) ([]*model.Post, error) {
	query := `
	SELECT post_id, session_id, user_name, post_title, post_content, image_urls, created_at, is_archived
	FROM posts 
	WHERE is_archived = false
	ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		r.logger.Error("failed to get all posts", slog.Any("error", err))
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
			r.logger.Error("failed to scan post", slog.Any("error", err))
			return nil, logger.ErrorWrapper("repository", "GetAllPosts", "scan post row", err)
		}
		posts = append(posts, &post)
	}
	// Logging errors during reading the query results
	if err = rows.Err(); err != nil {
		r.logger.Error("row iteration error in GetAllPosts", slog.Any("error", err))
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
		// do I need to log this one here?
		r.logger.Error("failed to archive post", slog.Any("error", err))
		return logger.ErrorWrapper("repository", "ArchivePosts", "update is_archived", err)
	}

	// To check whether we changed rows
	affected, err := result.RowsAffected()
	if err != nil {
		// Same questions
		r.logger.Error("failed to get rows affected for archiving", slog.Any("error", err))
		return logger.ErrorWrapper("repository", "ArchivePosts", "rows affected", err)
	}

	if affected == 0 {
		// If we did not update any row, it means there is no such post
		r.logger.Warn("no post archived", slog.Any("post_id", postID))
		return model.ErrPostNotFound
	}

	return nil
}
