package handler

import (
	"1337b04rd/internal/adapters/middleware"
	"1337b04rd/internal/domain/model"
	"1337b04rd/pkg/utils"
	"io"
	"net/http"
	"strings"
)

func (h *Handler) SubmitComment(w http.ResponseWriter, r *http.Request) {
	const fn = "SubmitComment"

	// Only allow POST method
	if r.Method != http.MethodPost {
		utils.LogWarn(h.logger, fn, "invalid method", "method", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Session check
	session := middleware.GetSessionFromContext(r.Context())
	if session == nil {
		utils.LogError(h.logger, fn, "session not found", nil)
		http.Error(w, "Session not found", http.StatusUnauthorized)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.LogError(h.logger, fn, "failed to parse multipart form", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Extract post ID from URL
	rawPath := strings.TrimPrefix(r.URL.Path, "/posts/")
	postID := strings.TrimSuffix(rawPath, "/comments")
	if postID == "" {
		utils.LogError(h.logger, fn, "missing post ID in URL", nil)
		http.Error(w, "Missing post ID", http.StatusBadRequest)
		return
	}

	// Read form values from html
	content := r.FormValue("comment")
	replyTo := r.FormValue("reply_to")
	newName := r.FormValue("name")

	// Optional: update username
	if newName != "" {
		if err := h.sessionService.OverrideUserName(r.Context(), session.SessionID, newName); err != nil {
			utils.LogError(h.logger, fn, "failed to override username", err)
		}
	}

	// Parse uploaded images (supports multiple)
	imageData := make(map[string]io.Reader)
	files := r.MultipartForm.File["file"]
	for _, fh := range files {
		file, err := fh.Open()
		if err != nil {
			utils.LogWarn(h.logger, fn, "skipped broken uploaded file", "file", fh.Filename, "error", err.Error())
			continue
		}
		defer file.Close()
		imageData[fh.Filename] = file
	}

	// Construct comment model
	comment := &model.Comment{
		SessionID:       session.SessionID,
		PostID:          utils.UUID(postID),
		UserName:        newName,
		ParentCommentID: utils.UUID(replyTo),
		Content:         content,
	}

	// Submit comment via service
	if err := h.commentService.CreateComment(r.Context(), comment, imageData); err != nil {
		utils.LogError(h.logger, fn, "failed to create comment", err)
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}

	utils.LogInfo(h.logger, fn, "comment created successfully", "post_id", postID, "session_id", string(session.SessionID))
	http.Redirect(w, r, "/posts/"+postID, http.StatusSeeOther)
}
