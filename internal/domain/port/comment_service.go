// Work on this later
package port

import (
	"1337b04rd/internal/domain/model"
	"1337b04rd/pkg/utils"
	"context"
	"io"
)

type CommentService interface {
	CreateComment(ctx context.Context, comment *model.Comment, imageData map[string]io.Reader) error
	GetCommentsByPostID(ctx context.Context, postID utils.UUID, includeArchived bool) ([]*model.Comment, error)
	GetCommentByID(ctx context.Context, commentID utils.UUID) (*model.Comment, error)
	ArchiveCommentsByPostIDTx(ctx context.Context, postID utils.UUID) error
}
