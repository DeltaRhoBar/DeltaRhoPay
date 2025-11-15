package handlers

import (
	"deltapay/internal/models"
	"deltapay/internal/services"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

type IndexPageData struct {
	Residents []models.Resident
	Beverages []models.Beverage
}

type IndexHandler struct {
	template *template.Template
	db services.Database
}

func NewIndexHandler(db services.Database) (*IndexHandler, error) {
	leadingZeroFunc := template.FuncMap{"leadingZero": func(i int) string {
		return fmt.Sprintf("%02d", i)
	}}
	trailingZeroFunc := template.FuncMap{"trailingZero": func(i int) string {
		return fmt.Sprintf("%.2f", float64(i) / 100)
	},	
	}
	tmpl, err := template.New("index.html").Funcs(leadingZeroFunc).Funcs(trailingZeroFunc).ParseFiles(filepath.Join("web", "templates", "index.html"))
	if err != nil {
	return nil, err
	}
	h := &IndexHandler{template: tmpl, db: db}
	return h, nil
}

func (h *IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	residents, err := h.db.GetResidents()
	if err != nil {
		http.Error(w, "Failed to get residents", http.StatusInternalServerError)
		return
	}
	beverages, err := h.db.GetBeverages()
	if err != nil {
		http.Error(w, "Failed to get beverages", http.StatusInternalServerError)
		return
	}
	data := &IndexPageData{Residents: residents, Beverages: beverages}
	err = h.template.Execute(w, data)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}
