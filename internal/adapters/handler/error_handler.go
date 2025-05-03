// Fix later
package handler

import (
	"html/template"
	"log/slog"
	"net/http"
)

var errorTpl = template.Must(template.ParseFiles("static/error.html"))

type ErrorData struct {
	Code    string
	Message string
}

func (h *Handler) ErrorPage(w http.ResponseWriter, r *http.Request) {
	// Optional: allow ?code=XXX&msg=Some+message
	code := r.URL.Query().Get("code")
	if code == "" {
		code = "500"
	}

	msg := r.URL.Query().Get("msg")
	if msg == "" {
		msg = "An unexpected error has occurred."
	}

	data := ErrorData{
		Code:    code,
		Message: msg,
	}

	if err := errorTpl.Execute(w, data); err != nil {
		slog.Error("failed to render error page", slog.Any("error", err))
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}
