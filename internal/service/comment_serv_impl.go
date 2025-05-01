package service

import (
	"1337b04rd/internal/domain/model"
	"1337b04rd/internal/domain/port"
	"1337b04rd/pkg/logger"
	"1337b04rd/pkg/utils"
	"context"
	"errors"
	"io"
	"log/slog"
	"strings"
	"time"
)

type CommentServiceImpl struct {
	repo        port.PostRepo
	commentRepo port.CommentRepo
	uploader    port.ImageUploader
	logger      *slog.Logger
}

func NewCommentServiceImpl(repo port.PostRepo, commentRepo port.CommentRepo, uploader port.ImageUploader, logger *slog.Logger) *CommentServiceImpl {
	return &CommentServiceImpl{
		repo:        repo,
		commentRepo: commentRepo,
		uploader:    uploader,
		logger:      logger}
}

// CreateComment adds a new comment or reply to a post
// LATER HANDLER: extract SessionID and PostID from the handler, prepare imageData map
func (s *CommentServiceImpl) CreateComment(ctx context.Context, comment *model.Comment, imageData map[string]io.Reader) error {

	// Assign CommentID and CreatedAt
	UUIDnum, err := utils.GenerateUUID()
	if err != nil {
		s.logger.Error("failed to assign UUID to CommentID", slog.Any("error", err))
		return logger.ErrorWrapper("service", "CreateComment", "generating UUID for commentID", model.ErrUUIDGeneration)
	}
	comment.CommentID = UUIDnum
	comment.CreatedAt = time.Now()

	// Check if post exists
	post, err := s.repo.GetPostByID(ctx, comment.PostID)
	if err != nil {
		return logger.ErrorWrapper("service", "CreateComment", "checking post existence", err)
	}

	// Check if it is archived post
	if post.IsArchived {
		return errors.New("cannot comment on archived post")
	}

	// Check if the ParentCommentID exists in the db
	if comment.ParentCommentID != "" {
		parent, err := s.commentRepo.GetCommentByID(ctx, comment.ParentCommentID)
		if err != nil {
			return logger.ErrorWrapper("service", "CreateComment", "checking parent comment", err)
		}
		// Check if it is not archived comment
		if parent.IsArchived {
			return errors.New("cannot reply to archived comment")
		}
		// Check if comment.ParentCommentID refers to the comment under the same PostID
		if parent.PostID != comment.PostID {
			return errors.New("cannot reply to comment from different post")
		}
	}

	// Ensure text or image exists
	if strings.TrimSpace(comment.Content) == "" && len(imageData) == 0 {
		return errors.New("comment cannot be empty")
	}

	// Check if there are any images attached
	if len(imageData) > 0 {

		// Upload comment images to buckets
		var urls []string
		for filename, content := range imageData {
			url, err := s.uploader.UploadCommentImage(string(post.PostID), string(comment.CommentID), filename, content)
			if err != nil {
				s.logger.Error("comment image upload failed", slog.String("filename", filename), slog.Any("error", err))
				return logger.ErrorWrapper("service", "CreateComment", "comment image uploading", err)
			}

			urls = append(urls, url)
		}

		// Save urls to db
		comment.ImageURLs = urls
	}

	// Save the comment to the repo
	if err := s.commentRepo.CreateComment(ctx, comment); err != nil {
		return logger.ErrorWrapper("service", "CreateComment", "saving comment to db", err)
	}

	return nil
}