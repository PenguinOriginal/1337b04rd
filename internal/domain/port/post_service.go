package port

import (
	"1337b04rd/internal/domain/model"
	"1337b04rd/pkg/utils"
	"context"
	"io"
)

type PostService interface {
	CreatePost(ctx context.Context, post *model.Post, imageData map[string]io.Reader) error
	GetAllPosts(ctx context.Context, archived bool) ([]*model.Post, error)
	GetPostByID(ctx context.Context, postID utils.UUID) (*model.Post, error)
	ArchivePost(ctx context.Context, postID utils.UUID) error
}
