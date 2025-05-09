package handler

import (
	"html/template"
	"net/http"

	"1337b04rd/pkg/utils"
)

type ErrorData struct {
	Code    string
	Message string
}

func (h *Handler) ErrorPage(w http.ResponseWriter, r *http.Request) {
	const ep = "ErrorPage"

	// Extract optional error code and message
	// If the user visited /error?code=403, this sets code = "403"
	code := r.URL.Query().Get("code")
	if code == "" {
		code = "500"
	}

	msg := r.URL.Query().Get("msg")
	if msg == "" {
		msg = "An unexpected error has occurred."
	}

	// Load the error.html template
	tpl, err := template.ParseFiles(templates["error"])
	if err != nil {
		utils.LogError(h.logger, ep, "failed to parse error.html", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	data := ErrorData{
		Code:    code,
		Message: msg,
	}

	// Render the template
	if err := tpl.Execute(w, data); err != nil {
		utils.LogError(h.logger, ep, "failed to render error template", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	utils.LogInfo(h.logger, ep, "error page rendered successfully")
}
