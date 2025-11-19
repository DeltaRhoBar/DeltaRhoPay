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


type orderPageData struct {
	Orders []models.Order
}

type OrdersPageHandler struct {
	template *template.Template
	db services.Database
}

func NewOrdersPageHandler(db services.Database) (*OrdersPageHandler, error) {
	trailingZeroFunc := template.FuncMap{"trailingZero": func(i int) string {
		return fmt.Sprintf("%.2f", float64(i) / 100)
	},	
	}
	tmpl, err := template.New("orders.html").Funcs(trailingZeroFunc).ParseFiles(filepath.Join("web", "templates", "orders.html"))
	if err != nil {
	return nil, err
	}
	h := &OrdersPageHandler{template: tmpl, db: db}
	return h, nil
}

func (h *OrdersPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	orders, err := h.db.GetOrders(0)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to get orders", http.StatusInternalServerError)
		return
	}
	data := &orderPageData{Orders: orders}
	err = h.template.Execute(w, data)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		log.Println(err);
		return
	}
}
