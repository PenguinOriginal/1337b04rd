package handler

import (
	"1337b04rd/internal/adapters/middleware"
	"1337b04rd/internal/domain/model"
	"1337b04rd/pkg/utils"
	"html/template"
	"io"
	"log/slog"
	"net/http"
)

// Must panics if the static file is missing
// Parse html files
var catalogTpl = template.Must(template.ParseFiles("static/catalog.html"))
var archiveTpl = template.Must(template.ParseFiles("static/archive.html"))
var postTpl = template.Must(template.ParseFiles("static/post.html"))
var archivePostTpl = template.Must(template.ParseFiles("static/archive-post.html"))
var createPostTpl = template.Must(template.ParseFiles("static/create-post.html"))

// GET /
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
		Session *middleware.SessionData // avatar
		Posts   []*model.Post
	}{
		Session: &middleware.SessionData{AvatarURL: session.AvatarURL},
		Posts:   posts,
	}

	// Renders the catalog page with data struct
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

func (h *Handler) Post(w http.ResponseWriter, r *http.Request) {
	session := middleware.GetSessionFromContext(r.Context())
	if session == nil {
		http.Error(w, "session not found", http.StatusUnauthorized)
		return
	}

	// Extract post ID from path
	postID := r.URL.Path[len("/posts/"):]
	if postID == "" {
		http.Error(w, "missing post ID", http.StatusBadRequest)
		return
	}

	// Fetch the post
	post, err := h.postService.GetPostByID(r.Context(), utils.UUID(postID))
	if err != nil {
		slog.Error("failed to fetch post", slog.Any("error", err))
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}

	// Fetch the comments
	comments, err := h.commentService.GetCommentsByPostID(r.Context(), post.PostID, post.IsArchived)
	if err != nil {
		slog.Error("failed to fetch comments", slog.Any("error", err))
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}

	data := struct {
		Post     *model.Post
		Comments []*model.Comment
		Session  *middleware.SessionData
	}{
		Post:     post,
		Comments: comments,
		Session:  &middleware.SessionData{AvatarURL: session.AvatarURL},
	}

	var tpl *template.Template
	if post.IsArchived {
		tpl = archivePostTpl
	} else {
		tpl = postTpl
	}

	if err := tpl.Execute(w, data); err != nil {
		slog.Error("failed to render post page", slog.Any("error", err))
		http.Error(w, "render error", http.StatusInternalServerError)
	}
}

func (h *Handler) CreatePostForm(w http.ResponseWriter, r *http.Request) {
	session := middleware.GetSessionFromContext(r.Context())
	if session == nil {
		http.Error(w, "session not found", http.StatusUnauthorized)
		return
	}

	data := struct {
		Session *middleware.SessionData
	}{
		Session: &middleware.SessionData{AvatarURL: session.AvatarURL},
	}

	if err := createPostTpl.Execute(w, data); err != nil {
		slog.Error("failed to render create-post.html", slog.Any("error", err))
		http.Error(w, "render error", http.StatusInternalServerError)
	}
}

// POST /posts
func (h *Handler) SubmitPost(w http.ResponseWriter, r *http.Request) {
	session := middleware.GetSessionFromContext(r.Context())
	if session == nil {
		http.Error(w, "session not found", http.StatusUnauthorized)
		return
	}

	// Limit file upload size to ~10MB
	// Go default limit is 32MB
	r.ParseMultipartForm(10 << 20) // 10 MB

	// Get text input from html files
	title := r.FormValue("subject")
	content := r.FormValue("comment")
	name := r.FormValue("name")

	// Get uploaded file from html
	file, fileHeader, err := r.FormFile("file")
	// Create imageData for service methods
	var imageData map[string]io.Reader
	if err == nil && file != nil {
		defer file.Close()
		imageData = map[string]io.Reader{fileHeader.Filename: file}
	}

	// If user entered a new name → override stored name
	if name != "" {
		err := h.sessionService.OverrideUserName(r.Context(), session.SessionID, name)
		if err != nil {
			slog.Error("failed to override user name", slog.Any("error", err))
			http.Redirect(w, r, "/error", http.StatusSeeOther)
			return
		}
	}

	// Create the post model
	post := &model.Post{
		SessionID: session.SessionID,
		UserName:  name, // If blank, default will be preserved
		Title:     title,
		Content:   content,
	}

	if err := h.postService.CreatePost(r.Context(), post, imageData); err != nil {
		slog.Error("failed to create post", slog.Any("error", err))
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}

	// Redirect to the new post page
	http.Redirect(w, r, "/posts/"+string(post.PostID), http.StatusSeeOther)
}
