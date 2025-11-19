package handlers

import (
	"deltapay/internal/services"
	"encoding/json"
	"net/http"
	"strconv"
)

type pageNumber struct{
	Nr int
}

type GetOrdersHandler struct {
	db services.Database
}

func NewGetOrdersHandler(db services.Database) *GetOrdersHandler {
	h := &GetOrdersHandler{db: db}	
	return h
}

func (h *GetOrdersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	pageStr := r.URL.Query().Get("Nr")
	pageNr, err := strconv.Atoi(pageStr)
	if err != nil {
		pageNr = 1
	}
	orders, err := h.db.GetOrders(pageNr);
	if err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]any{
		"orders": orders,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
