package handlers

import (
	"deltapay/internal/models"
	"deltapay/internal/services"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)


type orderPageData struct {
	Orders []models.Order
}

type OrdersPageHandler struct {
	template *template.Template
	db services.Database
}

func NewOrdersPageHandler(db services.Database) (*OrdersPageHandler, error) {
	tmpl, err := template.New("orders.html").ParseFiles(filepath.Join("web", "templates", "orders.html"))
	if err != nil {
	return nil, err
	}
	h := &OrdersPageHandler{template: tmpl, db: db}
	return h, nil
}

func (h *OrdersPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	orders, err := h.db.GetOrders()
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
