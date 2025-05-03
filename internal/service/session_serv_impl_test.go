// File: internal/service/session_serv_impl_test.go
package service

import (
	"1337b04rd/internal/domain/model"
	"1337b04rd/pkg/utils"
	"context"
	"io"
	"log/slog"
	"testing"
	"time"
)

type sessionServiceWithMockAvatar struct {
	*SessionServiceImpl
	MockAvatarURL string
}

func (s *sessionServiceWithMockAvatar) fetchRandomAvatar(ctx context.Context) (string, error) {
	return s.MockAvatarURL, nil
}

func TestCreateSession(t *testing.T) {
	repo := &MockSessionRepo{Sessions: make(map[utils.UUID]*model.Session)}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	service := &sessionServiceWithMockAvatar{
		SessionServiceImpl: NewSessionServiceImpl(repo, nil, nil, logger),
		MockAvatarURL:      "https://rickandmortyapi.com/api/character/avatar/1.jpeg",
	}

	sess := &model.Session{}
	err := service.CreateSession(context.Background(), sess)
	if err != nil {
		t.Fatalf("CreateSession failed: %v", err)
	}
	if repo.Saved == nil {
		t.Error("expected session to be saved")
	}
	if sess.AvatarURL == "" {
		t.Error("expected avatar to be assigned")
	}
	if sess.SessionID == "" {
		t.Error("expected session ID to be generated")
	}
	if time.Until(sess.ExpiresAt) < 6*time.Hour {
		t.Error("expected session to last for ~7 days")
	}
}

func TestGetSessionByID(t *testing.T) {
	id := utils.UUID("sess123")
	session := &model.Session{SessionID: id, AvatarURL: "https://example.com/avatar.jpg"}
	repo := &MockSessionRepo{Sessions: map[utils.UUID]*model.Session{id: session}}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	svc := NewSessionServiceImpl(repo, nil, nil, logger)

	got, err := svc.GetSessionByID(context.Background(), id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.AvatarURL != session.AvatarURL {
		t.Errorf("avatar mismatch: got %v, want %v", got.AvatarURL, session.AvatarURL)
	}
}

func TestDeleteExpiredSessions(t *testing.T) {
	repo := &MockSessionRepo{Sessions: make(map[utils.UUID]*model.Session)}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	svc := NewSessionServiceImpl(repo, nil, nil, logger)

	err := svc.DeleteExpiredSessions(context.Background())
	if err != nil {
		t.Fatalf("DeleteExpiredSessions failed: %v", err)
	}
	if !repo.Expired {
		t.Error("expected expired sessions to be deleted")
	}
}

func TestOverrideUserName(t *testing.T) {
	postRepo := &MockPostRepo{}
	commentRepo := &MockCommentRepo{}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	svc := NewSessionServiceImpl(nil, postRepo, commentRepo, logger)

	err := svc.OverrideUserName(context.Background(), "session-abc", "Rick")
	if err != nil {
		t.Fatalf("OverrideUserName failed: %v", err)
	}
	if !postRepo.UpdatedName || !commentRepo.UpdatedName {
		t.Error("expected user names to be updated in both post and comment repos")
	}
}