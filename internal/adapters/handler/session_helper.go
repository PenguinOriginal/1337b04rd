package handler

import (
	"1337b04rd/internal/adapters/middleware"
	"1337b04rd/pkg/utils"
	"log/slog"
	"net/http"
)

// CheckAndReturnSession safely extracts the session and logs error if missing.
func CheckAndReturnSession(w http.ResponseWriter, r *http.Request, logger *slog.Logger, functionName string) *middleware.SessionData {
	session := middleware.GetSessionFromContext(r.Context())
	if session == nil {
		utils.LogError(logger, functionName, "session not found", nil)
		http.Error(w, "session not found", http.StatusUnauthorized)
		return nil
	}
	return &middleware.SessionData{AvatarURL: session.AvatarURL}
}
