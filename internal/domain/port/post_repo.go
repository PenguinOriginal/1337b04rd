// Interfaces implemented by adapters
package port

import (
	"1337b04rd/internal/domain/model"
	"1337b04rd/pkg/utils"
	"context"
	"database/sql"
)

type PostRepo interface {
	CreatePost(ctx context.Context, post *model.Post) error
	GetPostByID(ctx context.Context, id utils.UUID) (*model.Post, error)
	GetAllPosts(ctx context.Context, archived bool) ([]*model.Post, error)
	ArchivePostTx(ctx context.Context, tx *sql.Tx, postID utils.UUID) error
	UpdateUserNameForSession(ctx context.Context, sessionID utils.UUID, newName string) error
}
