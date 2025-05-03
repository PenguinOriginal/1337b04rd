package service

import (
	"1337b04rd/internal/domain/model"
	"1337b04rd/pkg/utils"
	"context"
	"io"
	"log/slog"
	"strings"
	"testing"
	"time"
)

func TestCreatePost(t *testing.T) {
	ctx := context.Background()
	post := &model.Post{
		Title:     "Test",
		Content:   "Hello world",
		SessionID: "session123",
	}

	mockRepo := &MockPostRepo{Posts: map[utils.UUID]*model.Post{}}
	svc := NewPostServiceImpl(mockRepo, nil, nil, &MockUploader{}, slog.New(slog.NewTextHandler(io.Discard, nil)))

	err := svc.CreatePost(ctx, post, map[string]io.Reader{
		"image.png": strings.NewReader("fake image data"),
	})

	if err != nil {
		t.Fatalf("CreatePost failed: %v", err)
	}

	if mockRepo.CreatedPost == nil {
		t.Errorf("expected CreatePost to store post")
	}
	if len(mockRepo.CreatedPost.ImageURLs) == 0 {
		t.Errorf("expected uploaded image URL, got none")
	}
}

func TestGetAllPosts(t *testing.T) {
	postID := utils.UUID("post1")
	mockRepo := &MockPostRepo{
		Posts: map[utils.UUID]*model.Post{
			postID: {PostID: postID, Title: "Sample", SessionID: "abc", IsArchived: false},
		},
	}
	svc := NewPostServiceImpl(mockRepo, nil, nil, nil, slog.New(slog.NewTextHandler(io.Discard, nil)))

	posts, err := svc.GetAllPosts(context.Background(), false)
	if err != nil {
		t.Fatalf("GetAllPosts failed: %v", err)
	}
	if len(posts) != 1 || posts[0].PostID != postID {
		t.Errorf("unexpected posts returned: %+v", posts)
	}
}

func TestGetPostByID(t *testing.T) {
	postID := utils.UUID("p-123")
	mockRepo := &MockPostRepo{
		Posts: map[utils.UUID]*model.Post{
			postID: {PostID: postID, Title: "Title", SessionID: "sess1"},
		},
	}
	svc := NewPostServiceImpl(mockRepo, nil, nil, nil, slog.New(slog.NewTextHandler(io.Discard, nil)))

	post, err := svc.GetPostByID(context.Background(), postID)
	if err != nil {
		t.Fatalf("GetPostByID failed: %v", err)
	}
	if post.PostID != postID {
		t.Errorf("unexpected post ID: got %v", post.PostID)
	}
}

func TestArchivePost_NoComments_OlderThan10Min(t *testing.T) {
	postID := utils.UUID("archive-post")
	createdAt := time.Now().Add(-11 * time.Minute)
	mockRepo := &MockPostRepo{
		Posts: map[utils.UUID]*model.Post{
			postID: {
				PostID:    postID,
				CreatedAt: createdAt,
			},
		},
	}
	mockComment := &MockCommentRepo{LatestTime: nil}
	svc := NewPostServiceImpl(mockRepo, mockComment, nil, nil, slog.New(slog.NewTextHandler(io.Discard, nil)))

	err := svc.ArchivePost(context.Background(), postID)
	if err != nil {
		t.Fatalf("ArchivePost failed: %v", err)
	}
	if !mockRepo.Posts[postID].IsArchived {
		t.Errorf("expected post to be archived")
	}
}
