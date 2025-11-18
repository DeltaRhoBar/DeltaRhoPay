package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"
)

type LoginPageHandler struct {
	template *template.Template
}

func NewLoginPageHandler() (*LoginPageHandler, error) {
	tmpl, err := template.New("login.html").ParseFiles(filepath.Join("web", "templates", "login.html"))
	if err != nil {
	return nil, err
	}
	h := &LoginPageHandler{template: tmpl}
	return h, nil
}

func (h *LoginPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.template.Execute(w, nil)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}
