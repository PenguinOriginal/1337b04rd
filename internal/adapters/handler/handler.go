// HTTP adapter
package handler

// all handlers
// Serve /static using http.FileServer

import (
	"1337b04rd/internal/domain/port"
	"log/slog"
)

type Handler struct {
	postService    port.PostService
	commentService port.CommentService
	sessionService port.SessionService
	logger         *slog.Logger
}

func NewHandler(post port.PostService, comment port.CommentService, session port.SessionService, logger *slog.Logger) *Handler {
	return &Handler{
		postService:    post,
		commentService: comment,
		sessionService: session,
		logger:         logger,
	}
}
