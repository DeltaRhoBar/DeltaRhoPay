package handlers

import (
	"deltapay/internal/services"
	"net/http"
	"encoding/json"
)

type RemoveBeverageHandler struct {
	db services.Database
}

type oldBeverage struct {
	Name string
}

func NewRemoveBeverageHandler(db services.Database) *RemoveBeverageHandler {
	h := &RemoveBeverageHandler{db: db}	
	return h
}

func (h *RemoveBeverageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var beverage oldBeverage
	if err := json.NewDecoder(r.Body).Decode(&beverage); err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	err := h.db.RemoveBeverage(beverage.Name)
	if err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
