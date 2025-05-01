package imageuploader

// Implementation of image uploader interface (from the port)

// This file:
// Validates allowed extensions
// Calls into storage.go
// Wraps errors
// Returns image URLS

import (
	"1337b04rd/pkg/logger"
	"fmt"
	"io"
	"log/slog"
	"path/filepath"
	"strings"
)

var allowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".webp": true,
	".bmp":  true,
	".svg":  true,
}

// Shared function between UploadPostImage and UploadCommentImage
// Validating image extension
func validateExtension(filename string) error {
	// Avoid empty filenames and traversal
	if filename == "" || strings.Contains(filename, "..") || strings.HasPrefix(filename, ".") {
		return fmt.Errorf("invalid or unsafe filename: %s", filename)
	}

	ext := strings.ToLower(filepath.Ext(filename))
	if !allowedExtensions[ext] {
		return fmt.Errorf("unsupported image extension: %s", ext)
	}
	return nil
}

// Upload image for a post (path: /<postID>/<filename>)
func (u *LocalUploader) UploadPostImage(postID, filename string, r io.Reader) (string, error) {

	// Validate file extension
	if err := validateExtension(filename); err != nil {
		u.Logger.Warn("invalid image extension", slog.String("filename", filename))
		return "", logger.ErrorWrapper("image_uploader", "UploadPostImage", "extension check", err)
	}

	bucketPath := filepath.Join(u.RootDir, postID)
	fullPath := filepath.Join(bucketPath, filename)

	imageURL, err := SaveImageFile(fullPath, r, u.Logger)
	if err != nil {
		return "", logger.ErrorWrapper("image_uploader", "UploadPostImage", "saving post image failed", err)
	}

	u.Logger.Info("post image uploaded successfully", slog.String("imageURL", imageURL))
	return imageURL, nil
}

// Upload image for a comment (path: /<postID>/comments/<commentID>/<filename>)
func (u *LocalUploader) UploadCommentImage(postID, commentID, filename string, r io.Reader) (string, error) {

	// Validate file extension
	if err := validateExtension(filename); err != nil {
		u.Logger.Warn("invalid image extension", slog.String("filename", filename))
		return "", logger.ErrorWrapper("image_uploader", "UploadCommentImage", "extension check", err)
	}

	commentImagePath := filepath.Join(u.RootDir, postID, "comments", commentID)
	fullPath := filepath.Join(commentImagePath, filename)

	imageURL, err := SaveImageFile(fullPath, r, u.Logger)
	if err != nil {
		return "", logger.ErrorWrapper("image_uploader", "UploadCommentImage", "saving comment image failed", err)
	}

	u.Logger.Info("comment image uploaded successfully", slog.String("imageURL", imageURL))
	return imageURL, nil
}
