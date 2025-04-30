// Work on this later
package port

import (
	"1337b04rd/internal/domain/model"
	"1337b04rd/pkg/utils"
	"context"
)

type CommentService interface {
	CreateComment(ctx context.Context, comment *model.Comment) error

	// GetCommentsByPostID retrieves all comments for a given post.
	// Needed to display the comment thread for a post.
	GetCommentsByPostID(ctx context.Context, postID utils.UUID) ([]*model.Comment, error)

	// // Maybe don't need this
	// GetCommentByID(ctx context.Context, commentID utils.UUID) (*model.Comment, error)

	// ArchiveCommentsByPostID marks all comments of a post as archived.
	// Used when a post is archived/deleted and comments shouldn't appear anymore.
	ArchiveCommentsByPostIDTx(ctx context.Context, postID utils.UUID) error
}
