package imageuploader

import (
	"bytes"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"testing"
)

func TestSaveImageFile_CreatesFile(t *testing.T) {
	tmpDir := t.TempDir()
	targetPath := filepath.Join(tmpDir, "nested", "image.png")

	content := []byte("hello")
	reader := bytes.NewReader(content)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	url, err := SaveImageFile(targetPath, reader, logger)
	if err != nil {
		t.Fatalf("SaveImageFile failed: %v", err)
	}

	if _, err := os.Stat(targetPath); err != nil {
		t.Errorf("expected file to exist: %v", err)
	}

	if url != "/"+targetPath {
		t.Errorf("expected URL /%s, got %s", targetPath, url)
	}
}

func TestSaveImageFile_FileExists(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "image.png")

	os.WriteFile(path, []byte("existing"), 0644)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	_, err := SaveImageFile(path, bytes.NewReader([]byte("new")), logger)
	if err == nil {
		t.Error("expected error for existing file, got nil")
	}
}
