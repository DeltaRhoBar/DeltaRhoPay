package handlers

import (
	"deltapay/internal/services"
	"net/http"
	"encoding/json"
)

type AddOrderHandler struct {
	db services.Database
}

type orderedBeverage struct {
	Name string
	Amount int
}

type order struct {
	Beverages []orderedBeverage
	R_floor int
	R_nr int
}

func NewAddOrderHandler(db services.Database) *AddOrderHandler {
	h := &AddOrderHandler{db: db}	
	return h
}

func (h *AddOrderHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var order order 
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	
	for _, beverage := range order.Beverages {
		err := h.db.AddOrder(beverage.Name,beverage.Amount, order.R_floor, order.R_nr)
		if err != nil {
			http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

