// File: internal/service/mocks_test.go
package service

import (
	"1337b04rd/internal/domain/model"
	"1337b04rd/pkg/utils"
	"context"
	"database/sql"
	"io"
	"time"
)

// ========== Mock SessionRepo ==========
type MockSessionRepo struct {
	Sessions map[utils.UUID]*model.Session
	Saved    *model.Session
	Expired  bool
}

func (m *MockSessionRepo) CreateSession(ctx context.Context, session *model.Session) error {
	m.Saved = session
	m.Sessions[session.SessionID] = session
	return nil
}

func (m *MockSessionRepo) GetSessionByID(ctx context.Context, id utils.UUID) (*model.Session, error) {
	sess, ok := m.Sessions[id]
	if !ok {
		return nil, model.ErrSessionNotFound
	}
	return sess, nil
}

func (m *MockSessionRepo) DeleteExpiredSession(ctx context.Context) error {
	m.Expired = true
	return nil
}

// ========== Mock PostRepo ==========
type MockPostRepo struct {
	Posts       map[utils.UUID]*model.Post
	CreatedPost *model.Post
	ArchivedID  utils.UUID
	UpdatedName bool
}

func (m *MockPostRepo) CreatePost(ctx context.Context, post *model.Post) error {
	m.CreatedPost = post
	m.Posts[post.PostID] = post
	return nil
}

func (m *MockPostRepo) GetPostByID(ctx context.Context, id utils.UUID) (*model.Post, error) {
	p, ok := m.Posts[id]
	if !ok {
		return nil, model.ErrPostNotFound
	}
	return p, nil
}

func (m *MockPostRepo) GetAllPosts(ctx context.Context, archived bool) ([]*model.Post, error) {
	var result []*model.Post
	for _, p := range m.Posts {
		if p.IsArchived == archived {
			result = append(result, p)
		}
	}
	return result, nil
}

func (m *MockPostRepo) ArchivePostTx(ctx context.Context, tx *sql.Tx, postID utils.UUID) error {
	if post, ok := m.Posts[postID]; ok {
		post.IsArchived = true
		m.ArchivedID = postID
		return nil
	}
	return model.ErrPostNotFound
}

func (m *MockPostRepo) UpdateUserNameForSession(ctx context.Context, sessionID utils.UUID, newName string) error {
	m.UpdatedName = true
	return nil
}

// ========== Mock CommentRepo ==========
type MockCommentRepo struct {
	CreatedComment *model.Comment
	LatestTime     *time.Time
	UpdatedName    bool
}

func (m *MockCommentRepo) CreateComment(ctx context.Context, comment *model.Comment) error {
	m.CreatedComment = comment
	return nil
}

func (m *MockCommentRepo) GetCommentByID(ctx context.Context, id utils.UUID) (*model.Comment, error) {
	return &model.Comment{CommentID: id, PostID: "post123", IsArchived: false}, nil
}

func (m *MockCommentRepo) GetCommentsByPostID(ctx context.Context, postID utils.UUID, includeArchived bool) ([]*model.Comment, error) {
	return []*model.Comment{{CommentID: "c1", PostID: postID, Content: "Sample"}}, nil
}

func (m *MockCommentRepo) GetLatestCommentTime(ctx context.Context, postID utils.UUID) (*time.Time, error) {
	return m.LatestTime, nil
}

func (m *MockCommentRepo) ArchiveCommentByPostIDTx(ctx context.Context, tx *sql.Tx, postID utils.UUID) error {
	return nil
}

func (m *MockCommentRepo) UpdateUserNameForSession(ctx context.Context, sessionID utils.UUID, newName string) error {
	m.UpdatedName = true
	return nil
}

// ========== Mock Uploader ==========
type MockUploader struct{}

func (m *MockUploader) UploadPostImage(postID, filename string, r io.Reader) (string, error) {
	return "https://mock.upload/post.png", nil
}

func (m *MockUploader) UploadCommentImage(postID, commentID, filename string, r io.Reader) (string, error) {
	return "https://mock.upload/comment.png", nil
}
