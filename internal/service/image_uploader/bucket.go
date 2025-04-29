package imageuploader

import (
	"1337b04rd/internal/domain/model"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

// if my post is be deleted (archived), should I change anything with the folders?
// in theory, even if I delete a post, and its bucket is still alive, I should not be able
// to create another folder with the same name (uuid)
// how should I handle this?

// is it possible to upload an image to the "deleted" post? How do I check it?

// Need to store URL paths in DB!! (via service)

type LocalUploader struct {
	RootDir string
	Logger  *slog.Logger
}

func NewLocalUploader(RootDir string, logger *slog.Logger) *LocalUploader {
	return &LocalUploader{RootDir: RootDir, Logger: logger}
}

// Creates a folder for each post (bucketName == PostID)
func (u *LocalUploader) CreateBucket(bucketname string) error {

	// Check if rootdir ("data") is writable
	testFile := filepath.Join(u.RootDir, ".testfile")
	file, err := os.Create(testFile)
	if err != nil {
		// log this because it is a system-level failure
		u.Logger.Error("directory is not writable",
			slog.String("rootdir", u.RootDir),
			slog.Any("error", err))
		return fmt.Errorf("failed to verify rootdir writable: %w", err)
	}

	file.Close()
	os.Remove(testFile)

	// Check if this bucket already exists
	bucketPath := fmt.Sprintf("%s/%s", u.RootDir, bucketname)
	if _, err := os.Stat(bucketPath); !os.IsNotExist(err) {
		// Do not need to log it because it is expected error
		// u.Logger.Warn("bucket already exists", slog.String("bucketPath", bucketPath))
		return model.ErrBucketAlreadyExists
	} else if !os.IsNotExist(err) {
		// Log this because it is unexpected failure to stat file
		u.Logger.Error("failed to check bucket existence",
			slog.String("bucket", bucketPath),
			slog.Any("error", err),
		)
		return fmt.Errorf("failed to stat bucket: %w", err)
	}

	// Create bucket
	err = os.Mkdir(bucketPath, 0o755)
	if err != nil {
		u.Logger.Error("failed to create bucket",
			slog.String("bucketName", bucketname),
			slog.Any("error", err))
		return fmt.Errorf("failed to create bucket %s: %w", bucketname, err)
	}

	u.Logger.Info("bucket created", slog.String("bucketPath", bucketPath))
	return nil
}
