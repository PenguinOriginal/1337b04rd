// Work on this later
package port

import (
	"1337b04rd/internal/domain/model"
	"1337b04rd/pkg/utils"
	"context"
	"database/sql"
	"time"
)

type CommentRepo interface {
	CreateComment(ctx context.Context, comment *model.Comment) error
	GetCommentByPostID(ctx context.Context, postID utils.UUID) ([]*model.Comment, error)
	GetLatestCommentTime(ctx context.Context, postID utils.UUID) (*time.Time, error)
	ArchiveCommentByPostIDTx(ctx context.Context, tx *sql.Tx, postID utils.UUID) error
}
