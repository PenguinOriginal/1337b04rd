package service

import (
	"1337b04rd/internal/domain/model"
	"1337b04rd/internal/domain/port"
	"1337b04rd/pkg/logger"
	"1337b04rd/pkg/utils"
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"
)

type CommentServiceImpl struct {
	repo        port.PostRepo
	commentRepo port.CommentRepo
	logger      *slog.Logger
}

func NewCommentServiceImpl(repo port.PostRepo, commentRepo port.CommentRepo, logger *slog.Logger) *CommentServiceImpl {
	return &CommentServiceImpl{
		repo:        repo,
		commentRepo: commentRepo,
		logger:      logger}
}

// CreateComment adds a new comment or reply to a post
// LATER: extract SessionID and PostID from the handler
// IMPROVE: handle saving of the images from the comments
func (s *CommentServiceImpl) CreateComment(ctx context.Context, comment *model.Comment) error {

	// Assign CommentID and CreatedAt
	UUIDnum, err := utils.GenerateUUID()
	if err != nil {
		s.logger.Error("failed to assign UUID to PostID",
			slog.Any("error", err))
		return logger.ErrorWrapper("service", "CreatePost", "generating UUID", model.ErrUUIDGeneration)
	}
	comment.CommentID = UUIDnum
	comment.CreatedAt = time.Now()

	// Check if post exists
	post, err := s.repo.GetPostByID(ctx, comment.PostID)
	if err != nil {
		// Do I need to log this?
		s.logger.Error("post does not exist", slog.Any("error", err))
		return logger.ErrorWrapper("service", "CreateComment", "checking post existence", err)
	}

	// Do I need this part?
	if post.PostID != comment.PostID {
		s.logger.Error("post does not exist", slog.Any("error", err))
		return logger.ErrorWrapper("service", "CreateComment", "creating a comment on a given post", err)
	}

	if post.IsArchived {
		return errors.New("cannot comment on archived post")
	}

	// Check if the ParentCommentID exists in the db
	if len(comment.CommentID) > 0 {
		parent, err := s.commentRepo.GetCommentByID(ctx, comment.ParentCommentID)
		if err != nil {
			return logger.ErrorWrapper("service", "CreateComment", "checking parent comment", err)
		}
		if parent.IsArchived {
			return errors.New("cannot reply to archived comment")
		}
	}

	// Ensure content or image exists
	if strings.TrimSpace(comment.Content) == "" {
		return errors.New("comment cannot be empty")
	}

	return nil
}
