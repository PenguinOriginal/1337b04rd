package port

import (
	"1337b04rd/internal/domain/model"
	"1337b04rd/pkg/utils"
	"context"
)

type PostService interface {
	// CreatePost inserts a new post into the database.
	// Used by frontend for creating new threads.
	CreatePost(ctx context.Context, post *model.Post) error

	// GetAllPosts retrieves all non-archived posts from the database.
	// Used to display the post catalog.
	GetAllPosts(ctx context.Context) ([]*model.Post, error)

	// GetPostByID retrieves a single post by its ID.
	// Used to view the full thread along with comments.
	GetPostByID(ctx context.Context, id utils.UUID) (*model.Post, error)

	// ArchivePost marks a post (and its comments) as archived.
	// Used for removing content from the board (e.g., moderation, TTL).
	ArchivePost(ctx context.Context, id utils.UUID) error
}
