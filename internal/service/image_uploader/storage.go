package imageuploader

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

func SaveImageFile(rootDir, postID, filename string, r io.Reader, logger *slog.Logger) (string, error) {

	// Check if bucket exists
	bucketPath := filepath.Join(rootDir, postID)
	if _, err := os.Stat(bucketPath); os.IsNotExist(err) {
		logger.Error("bucket not found",
			slog.String("bucketPath", bucketPath))
		return "", fmt.Errorf("bucket not found: %w", err)
	}

	// Create full path for the image
	imagePath := filepath.Join(bucketPath, filename)

	// Check for filename collision
	if _, err := os.Stat(imagePath); err == nil {
		logger.Warn("file already exists",
			slog.String("imagePath", imagePath))
		return "", fmt.Errorf("file already exists")
	}

	// Create file
	dst, err := os.Create(imagePath)
	if err != nil {
		logger.Error("failed to create file",
			slog.String("imagePath", imagePath),
			slog.Any("error", err))
		return "", fmt.Errorf("could not create file: %w", err)
	}
	defer dst.Close()

	// Copy content
	if _, err := io.Copy(dst, r); err != nil {
		logger.Error("failed to copy image data",
			slog.String("imagePath", imagePath),
			slog.Any("error", err))
		return "", fmt.Errorf("copy failed: %w", err)
	}

	// Public URL format (can be replaced with S3-style later)
	imageURL := fmt.Sprintf("/%s/%s/%s", rootDir, postID, filename)
	return imageURL, nil
}
