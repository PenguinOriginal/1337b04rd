// LATER: check dublicate error logging
package service

import (
	"1337b04rd/internal/domain/model"
	"1337b04rd/internal/domain/port"
	"1337b04rd/pkg/logger"
	"1337b04rd/pkg/utils"
	"context"
	"database/sql"
	"errors"
	"io"
	"log/slog"
	"time"
)

type PostServiceImpl struct {
	repo        port.PostRepo
	commentRepo port.CommentRepo
	db          *sql.DB
	uploader    port.ImageUploader
	logger      *slog.Logger
}

func NewPostServiceImpl(repo port.PostRepo, commentRepo port.CommentRepo, db *sql.DB, uploader port.ImageUploader, logger *slog.Logger) *PostServiceImpl {
	return &PostServiceImpl{
		repo:        repo,
		commentRepo: commentRepo,
		db:          db,
		uploader:    uploader,
		logger:      logger}
}

// imageData should come from handler
// LATER: extract sessionID from cookie in handler and inject into post.SessionID
func (s *PostServiceImpl) CreatePost(ctx context.Context, post *model.Post, imageData map[string]io.Reader) error {

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
		s.logger.Warn("invalid post input",
			slog.Any("error", err))
		return logger.ErrorWrapper("service", "CreatePost", "validation", err)
	}

	// Check if there are any images attached
	if len(imageData) > 0 {

		// Upload images to buckets
		var urls []string
		for filename, content := range imageData {
			url, err := s.uploader.UploadPostImage(string(post.PostID), filename, content)
			if err != nil {
				s.logger.Error("image upload failed", slog.String("filename", filename), slog.Any("error", err))
				return logger.ErrorWrapper("service", "CreatePost", "image uploading", err)
			}

			urls = append(urls, url)
		}

		// Save urls to db
		post.ImageURLs = urls
	}

	// Save to repo
	if err := s.repo.CreatePost(ctx, post); err != nil {
		s.logger.Error("failed to create post", slog.Any("err", err))
		return logger.ErrorWrapper("service", "CreatePost", "saving post to repo", err)
	}

	s.logger.Info("post created successfully", slog.String("postID", string(post.PostID)))
	return nil
}

// GetAllPosts retrieves all non-archived posts from the database.
// Used to display the post catalog.
func (s *PostServiceImpl) GetAllPosts(ctx context.Context) ([]*model.Post, error) {
	posts, err := s.repo.GetAllPosts(ctx)
	if err != nil {
		s.logger.Error("failed to fetch posts", slog.Any("error", err))
		return nil, logger.ErrorWrapper("service", "GetAllPosts", "fetching posts", err)
	}
	s.logger.Info("fetched posts successfully", slog.Int("count", len(posts)))
	return posts, nil
}

// GetPostByID retrieves a single post by its ID.
// Used to view the full thread along with comments.
func (s *PostServiceImpl) GetPostByID(ctx context.Context, id utils.UUID) (*model.Post, error) {
	post, err := s.repo.GetPostByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get post by id", slog.Any("error", err))
		return nil, logger.ErrorWrapper("service", "GetPostByID", "fetching post by id", err)
	}
	s.logger.Info("fetched post by id successfully", slog.String("post_id", string(id)))
	return post, nil
}

// ArchivePost marks a post (and its comments) as archived.
// Used for removing content from the board (e.g., moderation, TTL).
func (s *PostServiceImpl) ArchivePost(ctx context.Context, postID utils.UUID) error {

	// Get the post by ID to retrieve its latest comment
	post, err := s.repo.GetPostByID(ctx, postID)
	if err != nil {
		// Handle not found or database error
		return logger.ErrorWrapper("service", "ArchivePost", "getting post", err)
	}

	// Get the latest comment time
	latestCommentTime, err := s.commentRepo.GetLatestCommentTime(ctx, postID)
	if err != nil && !errors.Is(err, model.ErrCommentNotFound) {
		return logger.ErrorWrapper("service", "ArchivePost", "getting latest comment", err)
	}

	now := time.Now()
	shouldArchive := false

	if latestCommentTime == nil {
		// No comments — archive after 10 minutes
		if post.CreatedAt.Add(10 * time.Minute).Before(now) {
			shouldArchive = true
		}
	} else {
		// Has comments — archive after 15 minutes since latest
		if latestCommentTime.Add(15 * time.Minute).Before(now) {
			shouldArchive = true
		}
	}

	if !shouldArchive {
		s.logger.Info("post is not eligible for archival yet", slog.String("post_id", string(postID)))
		return nil
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		s.logger.Error("failed to begin transaction", slog.Any("error", err))
		return logger.ErrorWrapper("service", "ArchivePost", "starting tx", err)
	}
	// Archive the post
	if err := s.repo.ArchivePostTx(ctx, tx, postID); err != nil {
		tx.Rollback()
		s.logger.Error("failed to archive post", slog.Any("error", err))
		return logger.ErrorWrapper("service", "ArchivePost", "archiving post", err)
	}

	// Archive related comments
	if err := s.commentRepo.ArchiveCommentByPostIDTx(ctx, tx, postID); err != nil {
		tx.Rollback()
		s.logger.Error("failed to archive comments", slog.String("post_id", string(postID)), slog.Any("error", err))
		return logger.ErrorWrapper("service", "ArchivePost", "archiving comments", err)
	}

	// If both comments and post succeed, commit
	if err := tx.Commit(); err != nil {
		s.logger.Error("failed to commit transaction", slog.Any("error", err))
		return logger.ErrorWrapper("service", "ArchivePost", "committing tx", err)
	}

	s.logger.Info("post and comments are archived successfully", slog.String("post_id", string(postID)))
	return nil
}


// Check the comments implementation and repo. There are some related files
// What is ticker?