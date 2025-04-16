package port

import "1337b04rd/internal/domain/model"

type PostService interface {
	CreatePost(post *model.Post) error
	// others later
}
