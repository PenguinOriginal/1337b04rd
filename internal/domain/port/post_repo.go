// Interfaces implemented by adapters
package port

import (
	"1337b04rd/internal/domain/model"
	"1337b04rd/pkg/utils"
	"context"
)

type PostRepo interface {
	CreatePost(ctx context.Context, post *model.Post) error
	GetPostByID(ctx context.Context, id utils.UUID) (*model.Post, error)
	GetAllPosts(ctx context.Context) ([]*model.Post, error)
	ArchivePosts(ctx context.Context, postID utils.UUID) error
	DeleteExpiredPosts(ctx context.Context) error
}
