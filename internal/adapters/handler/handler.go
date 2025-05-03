// HTTP adapter
package handler

// all handlers
// Serve /static using http.FileServer

import (
	"1337b04rd/internal/domain/port"
)

type Handler struct {
	postService    port.PostService
	commentService port.CommentService
	sessionService port.SessionService
}

func NewHandler(post port.PostService, comment port.CommentService, session port.SessionService) *Handler {
	return &Handler{
		postService:    post,
		commentService: comment,
		sessionService: session,
	}
}

// func GetSessionFromContext(ctx context.Context) *model.Session {
// 	val := ctx.Value(middleware.sessionKey)
// 	if session, ok := val.(*model.Session); ok {
// 		return session
// 	}
// 	return nil
// }
