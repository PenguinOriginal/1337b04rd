// Checked
package service

import (
	"1337b04rd/internal/domain/model"
	"1337b04rd/internal/domain/port"
	"1337b04rd/pkg/logger"
	"1337b04rd/pkg/utils"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math/rand"
	"net/http"
	"time"
)

type SessionServiceImpl struct {
	sessionRepo port.SessionRepo
	postRepo    port.PostRepo
	commentRepo port.CommentRepo
	logger      *slog.Logger
}

func NewSessionServiceImpl(sessionRepo port.SessionRepo, postRepo port.PostRepo, commentRepo port.CommentRepo, logger *slog.Logger) *SessionServiceImpl {
	return &SessionServiceImpl{
		sessionRepo: sessionRepo,
		postRepo:    postRepo,
		commentRepo: commentRepo,
		logger:      logger}
}

// Stores a new session
// Called when new user and does not have a session yet
func (s *SessionServiceImpl) CreateSession(ctx context.Context, session *model.Session) error {

	// Assign SessionID
	UUID, err := utils.GenerateUUID()
	if err != nil {
		s.logger.Error("failed to generate session UUID", slog.Any("error", err))
		return logger.ErrorWrapper("service", "CreateSession", "generating session UUID", err)
	}
	session.SessionID = UUID

	// Get avatar URL
	avatarURL, err := s.fetchRandomAvatar(ctx)
	if err != nil {
		s.logger.Error("failed to fetch random avatar", slog.Any("error", err))
		return logger.ErrorWrapper("service", "CreateSession", "fetching random avatar", err)
	}

	session.AvatarURL = avatarURL
	session.CreatedAt = time.Now()
	session.ExpiresAt = session.CreatedAt.Add(7 * 24 * time.Hour)

	// Save session info in db
	if err := s.sessionRepo.CreateSession(ctx, session); err != nil {
		return logger.ErrorWrapper("service", "CreateSession", "saving session to db", err)
	}
	s.logger.Info("session created", slog.String("session_id", string(session.SessionID)))
	return nil
}

// Retrieves a session by its UUID
// Used to check who is making a request and fetch avatar info
func (s *SessionServiceImpl) GetSessionByID(ctx context.Context, id utils.UUID) (*model.Session, error) {
	session, err := s.sessionRepo.GetSessionByID(ctx, id)
	if err != nil {
		return nil, logger.ErrorWrapper("service", "GetSessionByID", "retrieving session", err)
	}
	return session, nil
}

// DeleteExpiredSessions cleans up all expired sessions
func (s *SessionServiceImpl) DeleteExpiredSessions(ctx context.Context) error {
	if err := s.sessionRepo.DeleteExpiredSession(ctx); err != nil {
		return logger.ErrorWrapper("service", "DeleteExpiredSessions", "deleting expired sessions", err)
	}
	s.logger.Info("expired sessions deleted successfully")
	return nil
}

func (s *SessionServiceImpl) OverrideUserName(ctx context.Context, sessionID utils.UUID, newName string) error {
	// Update posts
	if err := s.postRepo.UpdateUserNameForSession(ctx, sessionID, newName); err != nil {
		return logger.ErrorWrapper("service", "OverrideUserName", "updating post usernames", err)
	}

	// Update comments
	if err := s.commentRepo.UpdateUserNameForSession(ctx, sessionID, newName); err != nil {
		return logger.ErrorWrapper("service", "OverrideUserName", "updating comment usernames", err)
	}
	s.logger.Info("user name updated successfully for session", slog.String("session_id", string(sessionID)), slog.String("new_name", newName))
	return nil
}

func (s *SessionServiceImpl) fetchRandomAvatar(ctx context.Context) (string, error) {
	rand.Seed(time.Now().UnixNano())
	randomID := rand.Intn(826) + 1

	url := fmt.Sprintf("https://rickandmortyapi.com/api/character/%d", randomID)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("building avatar request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("fetching avatar: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status from Rick and Morty API: %d", resp.StatusCode)
	}

	var data struct {
		Image string `json:"image"`
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading avatar response: %w", err)
	}

	if err := json.Unmarshal(body, &data); err != nil {
		return "", fmt.Errorf("decoding avatar JSON: %w", err)
	}

	return data.Image, nil
}
