package logger

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
)

var MyLogger *slog.Logger

// GetLoggerObject creates and returns a structured logger that writes to the given file path.
func GetLoggerObject(filePath string) *slog.Logger {
	// Ensure parent directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}

	// File flags: create if not exists, append, read/write
	option := os.O_CREATE | os.O_RDWR | os.O_APPEND
	file, err := os.OpenFile(filePath, option, 0666)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}

	// Create the structured logger
	logger := slog.New(slog.NewJSONHandler(file, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))

	return logger
}

func ErrorWrapper(layer, functionName, context string, err error) error {
	return fmt.Errorf("%s %w\n", fmt.Sprintf("[Layer:%s,Function: %s,Context: %s]--->", layer, functionName, context), err)
}
