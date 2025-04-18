// Service implementation
package service

import (
	"1337b04rd/internal/domain/model"
	"1337b04rd/internal/domain/port"
	"1337b04rd/logger"
	"1337b04rd/pkg/utils"
	"context"
	"log/slog"
	"time"
)

// session creation, post expiry logic, validation

// What it should do?
// Validating user input before it touches the database.
// Checking that required fields are: Present, Not empty, In the right format
// Assigning default values (e.g., UUID, timestamps).

type PostServiceImpl struct {
	repo   port.PostRepo
	logger *slog.Logger
}

func NewPostServiceImpl(repo port.PostRepo, logger *slog.Logger) *PostServiceImpl {
	return &PostServiceImpl{repo: repo, logger: logger}
}

func (s *PostServiceImpl) CreatePost(ctx context.Context, post *model.Post) error {

	// Assign PostID and CreatedAt
	UUIDnum, err := utils.GenerateUUID()
	if err != nil {
		s.logger.Error("failed to assign UUID to PostID", slog.Any("error", err))
		return logger.ErrorWrapper("service", "CreatePost", "generating UUID", model.ErrUUIDGeneration)
	}
	post.PostID = UUIDnum
	post.CreatedAt = time.Now()

	// Check if title & session are not empty
	if err := post.ValidatePost(); err != nil {
		s.logger.Warn("invalid post input", slog.Any("error", err))
		return err
	}

	// Save to repo
	if err := s.repo.CreatePost(ctx, post); err != nil {
		s.logger.Error("failed to create post", slog.Any("err", err))
		return logger.ErrorWrapper("service", "CreatePost", "saving post to repo", err)
	}

	s.logger.Info("post created successfully", slog.String("postID", string(post.PostID)))
	return nil

	// create buckets for the images
	// save URLs of the images and attach to post.ImageURLs

	// need a separate function to save images to triple-s
	// should be before saving the post to repo!!!

}

// Work on this later:
// imageUploader is a function that returns err and urls, must be from triple-s

// // Optional: Save image files if any were attached
// if len(post.ImageURLs) > 0 {
// 	uploadedURLs := make([]string, 0, len(post.ImageURLs))
// 	for _, rawPath := range post.ImageURLs {
// 		url, err := s.imageUploader.Upload(ctx, rawPath) // Wrap this utility
// 		if err != nil {
// 			s.logger.Error("failed to upload image", slog.Any("error", err))
// 			return logger.ErrorWrapper("service", "CreatePost", "uploading images", model.ErrImageUpload)
// 		}
// 		uploadedURLs = append(uploadedURLs, url)
// 	}
// 	post.ImageURLs = uploadedURLs
// }


// Save image to triple-s
// get their url
// update the post creation with this urls?