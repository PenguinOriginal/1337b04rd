package imageuploader

// implementation of image uploader interface:
// implements port

// this function:
// validates allowed extensions
// calls into storage.go
// wraps errors
// returns URLS

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

func (u *LocalUploader) UploadImage(postID, filename string, r io.Reader) (string, error) {

	// validate file extension
	ext := strings.ToLower(filepath.Ext(filename))
	if !allowedExtensions[ext] {
		// Log this because it is unexpected user input
		u.Logger.Warn("invalid image extension",
			slog.String("filename", filename),
			slog.String("extension", ext))
		return "", logger.ErrorWrapper("image_uploader", "UploadImage", "checking extension", fmt.Errorf("unsupported image extension"))
	}

	imageURL, err := SaveImageFile(u.RootDir, postID, filename, r, u.Logger)
	if err != nil {
		return "", logger.ErrorWrapper("image_uploader", "UploadImage", "saving image failed", err)
	}

	u.Logger.Info("image uploaded successfully", slog.String("imageURL", imageURL))
	return imageURL, nil
}
