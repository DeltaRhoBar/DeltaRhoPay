package handlers

import (
	"deltapay/internal/services"
	"net/http"
	"encoding/json"
)

type AddDebtHandler struct {
	db services.Database
}

type order struct {
	Amount int
	R_floor int
	R_nr int
}

func NewAddDebtHandler(db services.Database) *AddDebtHandler {
	h := &AddDebtHandler{db: db}	
	return h
}

func (h *AddDebtHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var order order 
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	err := h.db.AddDebt(order.Amount, order.R_floor, order.R_nr)
	if err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

