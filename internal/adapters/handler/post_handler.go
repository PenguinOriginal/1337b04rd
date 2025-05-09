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
	// Only allow GET method
	if r.Method != http.MethodGet {
		utils.LogWarn(h.logger, "Catalog", "invalid method", "method", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session := CheckAndReturnSession(w, r, h.logger, "Catalog")
	if session == nil {
		return
	}

	posts, err := h.postService.GetAllPosts(r.Context(), false) // false = not archived
	if err != nil {
		utils.LogError(h.logger, "Catalog", "failed to get posts", err)
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}

	tpl, err := template.ParseFiles(templates["catalog"])
	if err != nil {
		utils.LogError(h.logger, "Catalog", "failed to load template", err)
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
		utils.LogError(h.logger, "Catalog", "failed to render template", err)
		http.Error(w, "render error", http.StatusInternalServerError)
	}
	utils.LogInfo(h.logger, "Catalog", "served catalog page")
}

// GET /archive
func (h *Handler) Archive(w http.ResponseWriter, r *http.Request) {
	// Only allow GET method
	if r.Method != http.MethodGet {
		utils.LogWarn(h.logger, "Archive", "invalid method", "method", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session := CheckAndReturnSession(w, r, h.logger, "Catalog")
	if session == nil {
		return
	}

	posts, err := h.postService.GetAllPosts(r.Context(), true)
	if err != nil {
		utils.LogError(h.logger, "Archive", "failed to get archived posts", err)
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}

	tpl, err := template.ParseFiles(templates["archive"])
	if err != nil {
		utils.LogError(h.logger, "Archive", "failed to load template", err)
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
		utils.LogError(h.logger, "Archive", "failed to render template", err)
		http.Error(w, "render error", http.StatusInternalServerError)
	}
	utils.LogInfo(h.logger, "Archive", "served archive page")
}

// GET /posts/{id}
func (h *Handler) Post(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.LogWarn(h.logger, "Post", "invalid method", "method", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session := CheckAndReturnSession(w, r, h.logger, "Catalog")
	if session == nil {
		return
	}

	// Extract post ID from path
	postID := r.URL.Path[len("/posts/"):]
	if postID == "" {
		utils.LogError(h.logger, "Post", "missing post ID", nil)
		http.Error(w, "missing post ID", http.StatusBadRequest)
		return
	}

	// Fetch the post
	post, err := h.postService.GetPostByID(r.Context(), utils.UUID(postID))
	if err != nil {
		utils.LogError(h.logger, "Post", "failed to get post", err)
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}

	// Fetch the comments
	comments, err := h.commentService.GetCommentsByPostID(r.Context(), post.PostID, post.IsArchived)
	if err != nil {
		utils.LogError(h.logger, "Post", "failed to get comments", err)
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}

	tplFile := templates["post"]
	if post.IsArchived {
		tplFile = templates["archive-post"]
	}

	tpl, err := template.ParseFiles(tplFile)
	if err != nil {
		utils.LogError(h.logger, "Post", "failed to load template", err)
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
		utils.LogError(h.logger, "Post", "failed to render post", err)
		http.Error(w, "render error", http.StatusInternalServerError)
	}
	utils.LogInfo(h.logger, "Post", "served post page")
}

// GET /create
func (h *Handler) CreatePostForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.LogWarn(h.logger, "CreatePostForm", "invalid method", "method", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session := CheckAndReturnSession(w, r, h.logger, "Catalog")
	if session == nil {
		return
	}

	tpl, err := template.ParseFiles(templates["create-post"])
	if err != nil {
		utils.LogError(h.logger, "CreatePostForm", "failed to load template", err)
		http.Error(w, "Template load error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Session *middleware.SessionData
	}{
		Session: &middleware.SessionData{AvatarURL: session.AvatarURL},
	}

	if err := tpl.Execute(w, data); err != nil {
		utils.LogError(h.logger, "CreatePostForm", "failed to render create post form", err)
		http.Error(w, "render error", http.StatusInternalServerError)
	}
	utils.LogInfo(h.logger, "CreatePostForm", "create post form submitted")
}

// POST /posts
func (h *Handler) SubmitPost(w http.ResponseWriter, r *http.Request) {
	// Only allow GET method
	if r.Method != http.MethodPost {
		utils.LogWarn(h.logger, "SubmitPost", "invalid method", "method", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session := middleware.GetSessionFromContext(r.Context())
	if session == nil {
		utils.LogError(h.logger, "SubmitPost", "session not found", nil)
		http.Error(w, "session not found", http.StatusUnauthorized)
		return
	}

	// Limit file upload size to ~10MB
	// Go default limit is 32MB
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.LogError(h.logger, "SubmitPost", "failed to parse multipart form", err)
		http.Error(w, "invalid form data", http.StatusBadRequest)
		return
	}

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
		if err := h.sessionService.OverrideUserName(r.Context(), session.SessionID, name); err != nil {
			utils.LogError(h.logger, "SubmitPost", "failed to override username", err)
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

	// Crete the post
	if err := h.postService.CreatePost(r.Context(), post, imageData); err != nil {
		utils.LogError(h.logger, "SubmitPost", "failed to create post", err)
		http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}
	utils.LogInfo(h.logger, "SubmitPost", "post created", "post_id", string(post.PostID))
	// Redirect to the new post page
	http.Redirect(w, r, "/posts/"+string(post.PostID), http.StatusSeeOther)
}
