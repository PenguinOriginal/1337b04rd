// Add method checks
package handler

import (
	"1337b04rd/internal/adapters/middleware"
	"1337b04rd/internal/domain/model"
	"1337b04rd/pkg/utils"
	"html/template"
	"io"
	"net/http"
)

// GET /
func (h *Handler) Catalog(w http.ResponseWriter, r *http.Request) {

	// Check method first?

	session := middleware.GetSessionFromContext(r.Context())
	if session == nil {
		http.Error(w, "session not found", http.StatusUnauthorized)
		return
	}

	posts, err := h.postService.GetAllPosts(r.Context(), false) // false = not archived
	if err != nil {
		h.logger.Error("failed to get posts", "error", err)
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}

	tpl, err := template.ParseFiles("static/catalog.html")
	if err != nil {
		h.logger.Error("failed to load catalog.html", "error", err)
		http.Error(w, "Template load error", http.StatusInternalServerError)
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
	if err := tpl.Execute(w, data); err != nil {
		h.logger.Error("failed to render catalog", "error", err)
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
		h.logger.Error("failed to get archived posts", "error", err)
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}

	tpl, err := template.ParseFiles("static/archive.html")
	if err != nil {
		h.logger.Error("failed to load archive.html", "error", err)
		http.Error(w, "Template load error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Session *middleware.SessionData
		Posts   any
	}{
		Session: &middleware.SessionData{AvatarURL: session.AvatarURL},
		Posts:   posts,
	}

	if err := tpl.Execute(w, data); err != nil {
		h.logger.Error("failed to render archive", "error", err)
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
		h.logger.Error("failed to get post", "error", err)
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}

	// Fetch the comments
	comments, err := h.commentService.GetCommentsByPostID(r.Context(), post.PostID, post.IsArchived)
	if err != nil {
		h.logger.Error("failed to get comments", "error", err)
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}

	tplFile := "static/post.html"
	if post.IsArchived {
		tplFile = "static/archive-post.html"
	}

	tpl, err := template.ParseFiles(tplFile)
	if err != nil {
		h.logger.Error("failed to load post template", "file", tplFile, "error", err)
		http.Error(w, "Template load error", http.StatusInternalServerError)
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

	if err := tpl.Execute(w, data); err != nil {
		h.logger.Error("failed to render post", "error", err)
		http.Error(w, "render error", http.StatusInternalServerError)
	}
}

func (h *Handler) CreatePostForm(w http.ResponseWriter, r *http.Request) {
	session := middleware.GetSessionFromContext(r.Context())
	if session == nil {
		http.Error(w, "session not found", http.StatusUnauthorized)
		return
	}

	tpl, err := template.ParseFiles("static/create-post.html")
	if err != nil {
		h.logger.Error("failed to load create-post.html", "error", err)
		http.Error(w, "Template load error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Session *middleware.SessionData
	}{
		Session: &middleware.SessionData{AvatarURL: session.AvatarURL},
	}

	if err := tpl.Execute(w, data); err != nil {
		h.logger.Error("failed to render create-post", "error", err)
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

	// Create imageData for service methods
	var imageData map[string]io.Reader

	// Get uploaded file from html
	file, fileHeader, err := r.FormFile("file")

	if err == nil && file != nil {
		defer file.Close()
		imageData = map[string]io.Reader{fileHeader.Filename: file}
	}

	// If user entered a new name → override stored name
	if name != "" {
		err := h.sessionService.OverrideUserName(r.Context(), session.SessionID, name)
		if err != nil {
			h.logger.Error("failed to override username", "error", err)
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
		h.logger.Error("failed to create post", "error", err)
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}

	h.logger.Info("post created", "post_id", post.PostID)
	// Redirect to the new post page
	http.Redirect(w, r, "/posts/"+string(post.PostID), http.StatusSeeOther)
}
