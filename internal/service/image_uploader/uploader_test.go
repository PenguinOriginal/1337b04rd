package imageuploader

import (
	"bytes"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"testing"
)

func TestUploadPostImage_Valid(t *testing.T) {
	tmpDir := t.TempDir()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	uploader := NewLocalUploader(tmpDir, logger)

	content := []byte("fake image data")
	reader := bytes.NewReader(content)

	url, err := uploader.UploadPostImage("post1", "image.png", reader)
	if err != nil {
		t.Fatalf("UploadPostImage failed: %v", err)
	}

	expectedPath := filepath.Join(tmpDir, "post1", "image.png")
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("expected file at %s, but it doesn't exist", expectedPath)
	}

	if url != "/"+expectedPath {
		t.Errorf("unexpected image URL: %s", url)
	}
}

func TestUploadPostImage_InvalidExtension(t *testing.T) {
	tmpDir := t.TempDir()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	uploader := NewLocalUploader(tmpDir, logger)

	reader := bytes.NewReader([]byte("fake"))

	_, err := uploader.UploadPostImage("post1", "script.exe", reader)
	if err == nil {
		t.Error("expected error for invalid extension, got nil")
	}
}

func TestUploadCommentImage_Valid(t *testing.T) {
	tmpDir := t.TempDir()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	uploader := NewLocalUploader(tmpDir, logger)

	content := []byte("img")
	reader := bytes.NewReader(content)

	url, err := uploader.UploadCommentImage("post1", "cmt123", "cat.png", reader)
	if err != nil {
		t.Fatalf("UploadCommentImage failed: %v", err)
	}

	expectedPath := filepath.Join(tmpDir, "post1", "comments", "cmt123", "cat.png")
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("expected file at %s, but it doesn't exist", expectedPath)
	}

	if url != "/"+expectedPath {
		t.Errorf("unexpected image URL: %s", url)
	}
}
