package service

import (
	"1337b04rd/internal/domain/model"
	"1337b04rd/pkg/utils"
	"context"
	"io"
	"log/slog"
	"strings"
	"testing"
)

func TestCreateComment_Basic(t *testing.T) {
	ctx := context.Background()
	postID := utils.UUID("post123")

	comment := &model.Comment{
		PostID:    postID,
		Content:   "Test comment",
		SessionID: "sess123",
	}

	mockPost := &MockPostRepo{
		Posts: map[utils.UUID]*model.Post{
			postID: {PostID: postID, IsArchived: false},
		},
	}
	mockComment := &MockCommentRepo{}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	svc := NewCommentServiceImpl(mockPost, mockComment, &MockUploader{}, logger)

	err := svc.CreateComment(ctx, comment, map[string]io.Reader{
		"file.png": strings.NewReader("dummy"),
	})
	if err != nil {
		t.Fatalf("CreateComment failed: %v", err)
	}
	if mockComment.CreatedComment == nil {
		t.Error("expected comment to be created")
	}
	if len(mockComment.CreatedComment.ImageURLs) == 0 {
		t.Error("expected image to be uploaded")
	}
	if mockComment.CreatedComment.CommentID == "" {
		t.Error("expected comment ID to be generated")
	}
}

func TestCreateComment_MissingPost(t *testing.T) {
	ctx := context.Background()
	comment := &model.Comment{PostID: "non-existent", Content: "Hello", SessionID: "sess123"}

	mockPost := &MockPostRepo{Posts: map[utils.UUID]*model.Post{}}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	svc := NewCommentServiceImpl(mockPost, &MockCommentRepo{}, &MockUploader{}, logger)

	err := svc.CreateComment(ctx, comment, nil)
	if err == nil {
		t.Fatal("expected error for missing post, got nil")
	}
}

func TestGetCommentsByPostID(t *testing.T) {
	ctx := context.Background()
	postID := utils.UUID("post123")

	mockRepo := &MockCommentRepo{}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	svc := NewCommentServiceImpl(nil, mockRepo, nil, logger)

	comments, err := svc.GetCommentsByPostID(ctx, postID, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(comments) != 1 {
		t.Errorf("expected 1 comment, got %d", len(comments))
	}
	if comments[0].PostID != postID {
		t.Errorf("unexpected comment postID: %v", comments[0].PostID)
	}
}
