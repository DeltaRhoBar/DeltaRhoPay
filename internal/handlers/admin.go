package handlers

import (
	"deltapay/internal/services"
	"deltapay/internal/models"
	"html/template"
	"net/http"
	"path/filepath"
	"fmt"	
)

type AdminPageData struct {
	Beverages []models.Beverage
}

type AdminHandler struct {
	template *template.Template
	db services.Database
}

func NewAdminHandler(db services.Database) (*AdminHandler, error) {
	trailingZeroFunc := template.FuncMap{"trailingZero": func(i int) string {
		return fmt.Sprintf("%.2f", float64(i) / 100)
	},	
	}
	tmpl, err := template.New("admin.html").Funcs(trailingZeroFunc).ParseFiles(filepath.Join("web", "templates", "admin.html"))
	if err != nil {
		return nil, err
	}
	h := &AdminHandler{template: tmpl, db: db}
	return h, nil
}

func (h *AdminHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	result, err := h.db.GetBeverages()
	if err != nil {
		http.Error(w, "Failed to get beverages", http.StatusInternalServerError)
		return
	}
	data := &AdminPageData{Beverages: result}
	err = h.template.Execute(w, data)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}
