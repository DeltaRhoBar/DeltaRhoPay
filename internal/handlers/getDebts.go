package handlers

import (
	"deltapay/internal/services"
	"encoding/json"
	"net/http"
)

type GetDebtsHandler struct {
	db services.Database
}

func NewGetDebtsHandler(db services.Database) *GetDebtsHandler {
	h := &GetDebtsHandler{db: db}	
	return h
}

func (h *GetDebtsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	orders, err := h.db.GetOrders(0);
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
