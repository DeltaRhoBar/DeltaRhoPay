package handlers

import (
	"deltapay/internal/models"
	"deltapay/internal/services"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)


type residentsPageData struct {
	Residents []models.Resident
}

type ResidentsPageHandler struct {
	template *template.Template
	db services.Database
}

func NewResidentsPageHandler(db services.Database) (*ResidentsPageHandler, error) {
	tmpl, err := template.New("residents.html").ParseFiles(filepath.Join("web", "templates", "residents.html"))
	if err != nil {
	return nil, err
	}
	h := &ResidentsPageHandler{template: tmpl, db: db}
	return h, nil
}

func (h *ResidentsPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	residents, err := h.db.GetAllResidents()
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to get residents", http.StatusInternalServerError)
		return
	}
	data := &residentsPageData{Residents: residents}
	err = h.template.Execute(w, data)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		log.Println(err);
		return
	}
}
