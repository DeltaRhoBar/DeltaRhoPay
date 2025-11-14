
package handlers

import (
    "net/http"
	"html/template"
	"path/filepath"
)

type AdminHandler struct {
	template *template.Template
}

func NewAdminHandler() (*AdminHandler, error) {
	t, err := template.ParseFiles(filepath.Join("web", "templates", "admin.html"))
	if err != nil {
	return nil, err
	}
	h := &AdminHandler{template: t}
	return h, nil
}

func (h *AdminHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.template.Execute(w, nil)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}
