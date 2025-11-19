package handlers

import (
	"deltapay/internal/models"
	"deltapay/internal/services"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"fmt"
)


type debtPageData struct {
	Debts []models.Debt
}

type DebtPageHandler struct {
	template *template.Template
	db services.Database
}

func NewDebtPageHandler(db services.Database) (*DebtPageHandler, error) {
	trailingZeroFunc := template.FuncMap{"trailingZero": func(i int) string {
		return fmt.Sprintf("%.2f", float64(i) / 100)
	},	
	}
	tmpl, err := template.New("debt.html").Funcs(trailingZeroFunc).ParseFiles(filepath.Join("web", "templates", "debt.html"))
	if err != nil {
	return nil, err
	}
	h := &DebtPageHandler{template: tmpl, db: db}
	return h, nil
}

func (h *DebtPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	debts, err := h.db.GetDebts()
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to get debts", http.StatusInternalServerError)
		return
	}
	data := &debtPageData{Debts: debts}
	err = h.template.Execute(w, data)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		log.Println(err);
		return
	}
}
