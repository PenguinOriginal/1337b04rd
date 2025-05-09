package utils

import "log/slog"

func LogError(log *slog.Logger, function, message string, err error) {
	log.Error(message,
		slog.String("layer", "handler"),
		slog.String("function", function),
		slog.Any("error", err),
	)
}

func LogWarn(log *slog.Logger, function, message string, keyVals ...any) {
	fields := []any{
		slog.String("layer", "handler"),
		slog.String("function", function),
	}
	fields = append(fields, keyVals...)
	log.Warn(message, fields...)
}

// Success
func LogInfo(logger *slog.Logger, function, message string, kv ...any) {
	args := []any{
		slog.String("layer", "handler"),
		slog.String("function", function),
	}
	args = append(args, kv...)
	logger.Info(message, args...)
}
