package handler

import (
	"1337b04rd/internal/adapters/middleware"
	"1337b04rd/internal/domain/model"
	"1337b04rd/pkg/utils"
	"io"
	"log/slog"
	"net/http"
)

func (h *Handler) SubmitComment(w http.ResponseWriter, r *http.Request) {
	session := middleware.GetSessionFromContext(r.Context())
	if session == nil {
		http.Error(w, "session not found", http.StatusUnauthorized)
		return
	}

	// Parse form and limit upload size
	r.ParseMultipartForm(10 << 20) // 10MB

	// Get post ID from URL path
	postID := r.URL.Path[len("/posts/"):]
	if idx := len(postID) - len("/comments"); idx > 0 {
		postID = postID[:idx]
	}
	if postID == "" {
		http.Error(w, "missing post ID", http.StatusBadRequest)
		return
	}

	// Read form values from html
	content := r.FormValue("comment")
	replyTo := r.FormValue("reply_to") // optional

	// Parse uploaded images (supports multiple)
	imageData := make(map[string]io.Reader)
	files := r.MultipartForm.File["file"]
	for _, fh := range files {
		file, err := fh.Open()
		if err != nil {
			slog.Warn("skipped broken uploaded file", slog.Any("error", err))
			continue
		}
		defer file.Close()
		imageData[fh.Filename] = file
	}

	// Override name
	name := r.FormValue("name")
	if name != "" {
		if err := h.sessionService.OverrideUserName(r.Context(), session.SessionID, name); err != nil {
			slog.Error("failed to override username", slog.Any("error", err))
		}
	}

	// Construct comment
	comment := &model.Comment{
		SessionID:       session.SessionID,
		PostID:          utils.UUID(postID),
		UserName:        name,
		ParentCommentID: utils.UUID(replyTo),
		Content:         content,
	}

	// Submit comment via service
	if err := h.commentService.CreateComment(r.Context(), comment, imageData); err != nil {
		slog.Error("failed to create comment", slog.Any("error", err))
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}

	// Redirect back to post
	http.Redirect(w, r, "/posts/"+postID, http.StatusSeeOther)
}
