package handler

import (
	"1337b04rd/internal/adapters/middleware"
	"html/template"
	"log/slog"
	"net/http"
)

var catalogTpl = template.Must(template.ParseFiles("static/catalog.html"))
var archiveTpl = template.Must(template.ParseFiles("static/archive.html"))

func (h *Handler) Catalog(w http.ResponseWriter, r *http.Request) {
	session := middleware.GetSessionFromContext(r.Context())
	if session == nil {
		http.Error(w, "session not found", http.StatusUnauthorized)
		return
	}

	posts, err := h.postService.GetAllPosts(r.Context(), false) // false = not archived
	if err != nil {
		slog.Error("failed to get posts", slog.Any("error", err))
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}

	data := struct {
		Session *middleware.SessionData
		Posts   any
	}{
		Session: &middleware.SessionData{AvatarURL: session.AvatarURL},
		Posts:   posts,
	}

	if err := catalogTpl.Execute(w, data); err != nil {
		slog.Error("failed to render catalog", slog.Any("error", err))
		http.Error(w, "render error", http.StatusInternalServerError)
	}
}

func (h *Handler) Archive(w http.ResponseWriter, r *http.Request) {
	session := middleware.GetSessionFromContext(r.Context())
	if session == nil {
		http.Error(w, "session not found", http.StatusUnauthorized)
		return
	}

	posts, err := h.postService.GetAllPosts(r.Context(), true) // true = archived
	if err != nil {
		slog.Error("failed to get archived posts", slog.Any("error", err))
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}

	data := struct {
		Session *middleware.SessionData
		Posts   any
	}{
		Session: &middleware.SessionData{AvatarURL: session.AvatarURL},
		Posts:   posts,
	}

	if err := archiveTpl.Execute(w, data); err != nil {
		slog.Error("failed to render archive", slog.Any("error", err))
		http.Error(w, "render error", http.StatusInternalServerError)
	}
}
