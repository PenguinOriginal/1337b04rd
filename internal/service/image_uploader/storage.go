package imageuploader

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

// This function is shared by UploadPostImage and UploadCommentImage
func SaveImageFile(fullPath string, r io.Reader, logger *slog.Logger) (string, error) {

	// Check if parent directory exists
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		logger.Error("failed to create directory", slog.String("path", filepath.Dir(fullPath)), slog.Any("error", err))
		return "", fmt.Errorf("could not create parent directory: %w", err)
	}

	// Check for filename collision
	if _, err := os.Stat(fullPath); err == nil {
		logger.Warn("file already exists", slog.String("imagePath", fullPath))
		return "", fmt.Errorf("file already exists")
	}

	// Create file
	dst, err := os.Create(fullPath)
	if err != nil {
		logger.Error("failed to create file",
			slog.String("imagePath", fullPath),
			slog.Any("error", err))
		return "", fmt.Errorf("could not create file: %w", err)
	}
	defer dst.Close()

	// Copy content
	if _, err := io.Copy(dst, r); err != nil {
		logger.Error("failed to copy image data",
			slog.String("imagePath", fullPath),
			slog.Any("error", err))
		return "", fmt.Errorf("copy failed: %w", err)
	}

	// Public URL format
	imageURL := fmt.Sprintf("/%s", fullPath)
	return imageURL, nil
}
