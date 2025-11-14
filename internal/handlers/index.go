package handlers

import (
	"deltapay/internal/models"
	"deltapay/internal/services"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

type PageData struct {
	Residents []models.Resident
}

type IndexHandler struct {
	template *template.Template
	db services.Database
}

func NewIndexHandler(db services.Database) (*IndexHandler, error) {
	leadingZeroFunc := template.FuncMap{"leadingZero": func(i int) string {
		return fmt.Sprintf("%02d", i)
	}}
	tmpl, err := template.New("index.html").Funcs(leadingZeroFunc).ParseFiles(filepath.Join("web", "templates", "index.html"))
	if err != nil {
	return nil, err
	}
	h := &IndexHandler{template: tmpl, db: db}
	return h, nil
}

func (h *IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	result, err := h.db.GetResidents()
	if err != nil {
		http.Error(w, "Failed to get residents", http.StatusInternalServerError)
		return
	}
	data := &PageData{Residents: result}
	err = h.template.Execute(w, data)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}
