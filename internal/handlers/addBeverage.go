package handlers

import (
	"deltapay/internal/services"
	"net/http"
	"encoding/json"
)

type AddBeverageHandler struct {
	db services.Database
}

type newBeverage struct {
	Name string
	Price int
}

func NewAddBeverageHandler(db services.Database) *AddBeverageHandler {
	h := &AddBeverageHandler{db: db}	
	return h
}

func (h *AddBeverageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var beverage newBeverage
	if err := json.NewDecoder(r.Body).Decode(&beverage); err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	err := h.db.AddBeverage(beverage.Name, beverage.Price)
	if err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
