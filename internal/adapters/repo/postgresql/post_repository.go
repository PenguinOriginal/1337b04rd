// DB adapter
package postgresql

import (
	"1337b04rd/internal/domain/model"
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type PostgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (r *PostgresRepo) CreatePost(ctx context.Context, post *model.Post) error {

	// ExecContext, QueryContext, QueryRowContext

	query := `
	INSERT INTO posts (post_id, session_id, user_name, post_content, post_title, image_urls, created_at, is_archived)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.ExecContext(ctx, query,
		post.PostID,
		post.SessionID,
		post.UserName,
		post.Content,
		post.Title,
		pq.Array(post.ImageURLs),
		post.CreatedAt,
		post.IsArchived,
	)
	return err
}
