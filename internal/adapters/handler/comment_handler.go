package handler

import (
	"1337b04rd/internal/adapters/middleware"
	"1337b04rd/internal/domain/model"
	"1337b04rd/pkg/utils"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

func (h *Handler) SubmitComment(w http.ResponseWriter, r *http.Request) {
	// Only allow POST method
	if r.Method != http.MethodPost {
		h.logger.Warn("invalid method", "method", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Session check
	session := middleware.GetSessionFromContext(r.Context())
	if session == nil {
		h.logger.Error("session not found in context")
		http.Error(w, "Session not found", http.StatusUnauthorized)
		return
	}

	// Parse form and limit upload size
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		h.logger.Error("failed to parse multipart form", "error", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Extract post ID from URL
	rawPath := strings.TrimPrefix(r.URL.Path, "/posts/")
	postID := strings.TrimSuffix(rawPath, "/comments")
	if postID == "" {
		http.Error(w, "missing post ID", http.StatusBadRequest)
		return
	}

	// Read form values from html
	content := r.FormValue("comment")
	replyTo := r.FormValue("reply_to")
	newName := r.FormValue("name")

	// Optional: update username
	if newName != "" {
		if err := h.sessionService.OverrideUserName(r.Context(), session.SessionID, newName); err != nil {
			h.logger.Error("failed to override username", "error", err)
		}
	}

	// Parse uploaded images (supports multiple)
	imageData := make(map[string]io.Reader)
	files := r.MultipartForm.File["file"]
	for _, fh := range files {
		file, err := fh.Open()
		if err != nil {
			h.logger.Warn("skipped broken uploaded file", slog.Any("error", err))
			continue
		}
		defer file.Close()
		imageData[fh.Filename] = file
	}

	// Construct comment
	comment := &model.Comment{
		SessionID:       session.SessionID,
		PostID:          utils.UUID(postID),
		UserName:        newName,
		ParentCommentID: utils.UUID(replyTo),
		Content:         content,
	}

	// Submit comment via service
	if err := h.commentService.CreateComment(r.Context(), comment, imageData); err != nil {
		h.logger.Error("failed to create comment", "post_id", postID, "error", err)
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}

	// Redirect back to post
	h.logger.Info("comment created successfully", "post_id", postID, "session_id", session.SessionID)
	http.Redirect(w, r, "/posts/"+postID, http.StatusSeeOther)
}
